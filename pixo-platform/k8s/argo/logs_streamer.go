package argo

import (
	"bytes"
	"context"
	"errors"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/gcs"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
	"time"
)

type Log struct {
	Step  string `json:"step"`
	Lines string `json:"lines"`
}

type LogsStreamer struct {
	WorkflowName  string
	argoClient    Client
	k8sClient     base.Client
	storageClient client.StorageClient
	bucketName    string
	namespace     string
	streams       map[string]chan Log
	numNodes      int
	numDone       int
	mtx           sync.Mutex
}

func NewLogsStreamer(k8sClient base.Client, argoClient Client, namespace, workflowName, bucketName string) (*LogsStreamer, error) {
	if namespace == "" {
		return nil, errors.New("namespace may not be empty")
	}

	if workflowName == "" {
		return nil, errors.New("workflowName may not be empty")
	}

	if bucketName == "" {
		return nil, errors.New("bucketName may not be empty")
	}

	storageClient, err := gcs.NewClient(gcs.Config{BucketName: bucketName})
	if err != nil {
		return nil, err
	}

	return &LogsStreamer{
		k8sClient:     k8sClient,
		argoClient:    argoClient,
		storageClient: storageClient,
		WorkflowName:  workflowName,
		namespace:     namespace,
		bucketName:    bucketName,
		streams:       make(map[string]chan Log),
	}, nil
}

func (s *LogsStreamer) Start(c context.Context) (chan Log, error) {
	workflow, err := s.argoClient.GetWorkflow(s.namespace, s.WorkflowName)
	if err != nil {
		return nil, err
	}

	if workflow == nil {
		return nil, errors.New("workflow not found")
	}

	for _, template := range workflow.Spec.Templates {

		if template.GetNodeType() != v1alpha1.NodeTypePod || s.getStream(template.Name) != nil {
			continue
		}

		s.addStream(template.Name)

		if workflow.Status.Phase == v1alpha1.WorkflowSucceeded || workflow.Status.Phase == v1alpha1.WorkflowFailed {
			go s.streamArchive(template.Name)

		} else if workflow.Status.Phase == v1alpha1.WorkflowRunning {
			go s.waitForTail(c, template, workflow)
		}
	}

	combinedStream := make(chan Log, 1000)

	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, stream := range s.streams {
		go s.combineStream(combinedStream, stream)
	}

	return combinedStream, nil
}

func (s *LogsStreamer) combineStream(combinedStream, stream chan Log) {
	for !IsClosed(stream) {
		newLog := <-stream
		if newLog.Step != "" && newLog.Lines != "" {
			log.Debug().Msgf("New log: %s %s", newLog.Step, newLog.Lines)
			combinedStream <- newLog
		}
	}

	log.Debug().Msgf("Closed stream. Num total: %d Num closed: %d", s.NumNodes(), s.NumDone())

	if !IsClosed(combinedStream) && s.IsDone() {
		log.Debug().Msg("All streams are closed. Closing combined stream")
		close(combinedStream)
	}
}

func (s *LogsStreamer) waitForTail(c context.Context, template v1alpha1.Template, workflow *v1alpha1.Workflow) {
	for {
		time.Sleep(1 * time.Second)

		if _, err := s.tail(c, template.Name, workflow); err == nil {
			break
		}
	}
}

func (s *LogsStreamer) tail(c context.Context, templateName string, workflow *v1alpha1.Workflow) (chan Log, error) {
	if workflow == nil {
		return nil, errors.New("workflow may not be nil")
	}

	if templateName == "" {
		return nil, errors.New("templateName may not be empty")
	}

	node, err := s.argoClient.GetNode(workflow, templateName)
	if err != nil {
		return nil, errors.Join(err, errors.New("unable to find node from template name"))
	}

	containerName := "main"
	podName := FormatPodName(node)

	var ioStream io.ReadCloser

	for {
		time.Sleep(1 * time.Second)

		node, err = s.argoClient.GetNode(workflow, node.TemplateName)
		if err != nil || node == nil || node.Pending() {
			continue
		}

		ioStream, err = s.k8sClient.GetPodLogs(c, s.namespace, podName, containerName)
		if err != nil {
			log.Debug().Err(err).Msgf("unable to get logs for pod %s", podName)
			continue
		}

		break
	}

	go s.readLogsForNode(node.TemplateName, ioStream)

	log.Debug().Msgf("started tailing logs for node %s", node.TemplateName)

	return s.getStream(node.TemplateName), nil
}

func (s *LogsStreamer) streamArchive(nodeName string) {
	readCloser, err := s.GetArchivedLogsForTemplate(context.Background(), nodeName)
	if err != nil {
		return
	}
	go s.readLogsForNode(nodeName, readCloser)
}

func (s *LogsStreamer) readLogsForNode(nodeName string, ioStream io.ReadCloser) {
	if nodeName == "" || ioStream == nil {
		return
	}
	defer ioStream.Close()

	stream := s.getStream(nodeName)
	if stream == nil {
		log.Debug().Msgf("stream for %s is nil", nodeName)
		return
	}

	log.Debug().Msgf("started streaming logs for %s", nodeName)
	for {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, ioStream); err != nil {
			log.Debug().Err(err).Msgf("unable to copy logs for %s", nodeName)
			s.markStreamDone(nodeName)
			break

		} else if buf.Len() == 0 {

			if !s.nodeIsDone(nodeName) {
				log.Debug().Msgf("no logs for running node %s", nodeName)
				time.Sleep(1 * time.Second)
				continue

			} else {
				log.Debug().Msgf("no logs for completed node %s", nodeName)
				s.markStreamDone(nodeName)
				break
			}

		}

		stream <- Log{
			Step:  nodeName,
			Lines: buf.String(),
		}

		log.Debug().Msgf("streamed log for %s", nodeName)
	}
}

func (s *LogsStreamer) GetArchivedLogsForTemplate(c context.Context, templateName string) (io.ReadCloser, error) {

	workflow, err := s.argoClient.GetWorkflow(s.namespace, s.WorkflowName)
	if err != nil || workflow == nil {
		return nil, errors.Join(err, errors.New("unable to get workflow"))
	}

	node, err := s.argoClient.GetNode(workflow, templateName)
	if err != nil || node == nil {
		return nil, errors.Join(err, errors.New("unable to get node"))
	}

	archive := Archive{
		BucketName:   s.bucketName,
		WorkflowName: s.WorkflowName,
		PodName:      FormatPodName(node),
	}

	if node.Phase != v1alpha1.NodeSucceeded && node.Phase != v1alpha1.NodeFailed {
		log.Debug().Msgf("unable to get archives, node %s is not done", node.TemplateName)
		return nil, errors.New("node is not done")
	}

	readCloser, err := s.storageClient.ReadFile(c, archive)
	if err != nil {
		return nil, errors.Join(err, errors.New("unable to read archived logs"))
	}

	return readCloser, nil
}
