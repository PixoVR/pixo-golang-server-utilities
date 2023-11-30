package argo

import (
	"bytes"
	"context"
	"errors"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/redis/go-redis/v9"
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
	workflowName   string
	argoClient     *Client
	k8sClient      *base.Client
	storageClient  client.StorageClient
	logsCache      *redis.Client
	bucketName     string
	namespace      string
	streams        map[string]chan Log
	combinedStream chan Log
	numNodes       int
	numDone        int
	mtx            sync.Mutex
}

type StreamerConfig struct {
	K8sClient     *base.Client
	ArgoClient    *Client
	StorageClient client.StorageClient
	LogsCache     *redis.Client
	Namespace     string
	WorkflowName  string
}

func (sc StreamerConfig) isValid() error {
	if sc.Namespace == "" {
		return errors.New("namespace may not be empty")
	}

	if sc.WorkflowName == "" {
		return errors.New("workflowName may not be empty")
	}

	if sc.StorageClient == nil {
		return errors.New("storage client may not be nil")
	}

	if sc.K8sClient == nil {
		return errors.New("k8s client may not be nil")
	}

	if sc.ArgoClient == nil {
		return errors.New("argo client may not be nil")
	}

	return nil
}

func NewLogsStreamer(config StreamerConfig) (*LogsStreamer, error) {
	if err := config.isValid(); err != nil {
		return nil, err
	}

	return &LogsStreamer{
		k8sClient:     config.K8sClient,
		argoClient:    config.ArgoClient,
		logsCache:     config.LogsCache,
		storageClient: config.StorageClient,
		workflowName:  config.WorkflowName,
		namespace:     config.Namespace,
		streams:       make(map[string]chan Log),
	}, nil
}

func (s *LogsStreamer) Start(ctx context.Context) (chan Log, error) {
	workflow, err := s.argoClient.GetWorkflow(s.namespace, s.workflowName)
	if err != nil {
		return nil, err
	}

	if workflow == nil {
		return nil, errors.New("workflow not found")
	}

	s.addStreams(workflow)
	s.combineStreams()
	go s.startStreaming(ctx, workflow)

	return s.combinedStream, nil
}

func (s *LogsStreamer) addStreams(workflow *v1alpha1.Workflow) {
	for _, template := range workflow.Spec.Templates {
		if !hasLogs(template) {
			continue
		}

		s.addStream(template.Name)
	}
}

func (s *LogsStreamer) combineStreams() {
	s.makeCombinedStream()

	for _, stream := range s.streams {
		go s.combineStream(stream)
	}
}

func (s *LogsStreamer) startStreaming(ctx context.Context, workflow *v1alpha1.Workflow) {
	for _, template := range workflow.Spec.Templates {
		if workflow.Status.Phase == v1alpha1.WorkflowSucceeded || workflow.Status.Phase == v1alpha1.WorkflowFailed {
			go s.streamArchive(ctx, template.Name)

		} else if workflow.Status.Phase == v1alpha1.WorkflowRunning {
			go s.waitForTail(ctx, template, workflow)
		}
	}
}

func (s *LogsStreamer) combineStream(stream chan Log) {
	for newLog := range stream {
		if newLog.Step != "" && newLog.Lines != "" {
			log.Debug().Msgf("New log: %s %s", newLog.Step, newLog.Lines)
			s.combinedStream <- newLog
		}
	}

	log.Debug().Msgf("Closed stream. Num total: %d Num closed: %d", s.NumNodes(), s.NumDone())

	if s.IsDone() {
		s.markComplete()
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

	for {
		time.Sleep(1 * time.Second)

		node, err = s.argoClient.GetNode(workflow, node.TemplateName)
		if err != nil || node == nil || node.Pending() {
			continue
		}

		readCloser, err := s.k8sClient.GetPodLogs(c, s.namespace, podName, containerName)
		if err != nil {
			log.Debug().Err(err).Msgf("unable to get logs for pod %s", podName)
			continue
		}

		go s.readLogsForNode(node.TemplateName, readCloser)
		break
	}

	return s.getStream(node.TemplateName), nil
}

func (s *LogsStreamer) streamArchive(ctx context.Context, nodeName string) {
	archiveReadCloser, err := s.GetArchivedLogsForTemplate(ctx, nodeName)
	if err != nil {
		return
	}

	go s.readLogsForNode(nodeName, archiveReadCloser)
}

func (s *LogsStreamer) readLogsForNode(nodeName string, readCloser io.ReadCloser) {
	if nodeName == "" || readCloser == nil {
		return
	}
	defer readCloser.Close()

	stream := s.getStream(nodeName)
	if stream == nil {
		log.Debug().Msgf("stream for %s is nil", nodeName)
		return
	}

	log.Debug().Msgf("started streaming logs for %s", nodeName)
	for {
		buf := new(bytes.Buffer)
		if written, err := io.Copy(buf, readCloser); err != nil {
			log.Debug().Err(err).Msgf("unable to copy logs for %s", nodeName)
			s.markStreamDone(nodeName)
			break

		} else if written == 0 {

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

		if buf.String() != "" {
			stream <- Log{
				Step:  nodeName,
				Lines: buf.String(),
			}

			log.Debug().Msgf("streamed log for %s", nodeName)
		}

	}
}

func (s *LogsStreamer) GetArchivedLogsForTemplate(c context.Context, templateName string) (io.ReadCloser, error) {

	workflow, err := s.argoClient.GetWorkflow(s.namespace, s.workflowName)
	if err != nil || workflow == nil {
		return nil, errors.Join(err, errors.New("unable to get workflow"))
	}

	node, err := s.argoClient.GetNode(workflow, templateName)
	if err != nil || node == nil {
		return nil, errors.Join(err, errors.New("unable to get node"))
	}

	archive := Archive{
		BucketName:   s.bucketName,
		WorkflowName: s.workflowName,
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

func hasLogs(template v1alpha1.Template) bool {
	return template.GetNodeType() == v1alpha1.NodeTypePod
}
