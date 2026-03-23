package workflows

import (
	"bytes"
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
	"time"
)

type Log struct {
	Step  string `json:"step"`
	Lines string `json:"lines"`
}

type LogsStreamer struct {
	workflowName         string
	argoClient           *Client
	k8sClient            kubernetes.Interface
	storageClient        blobstorage.StorageClient
	logsCache            *redis.Client
	bucketName           string
	namespace            string
	streams              map[string]chan Log
	combinedStream       chan Log
	numNodes             int
	numDone              int
	closed               bool
	combinedStreamClosed bool
	cancelFunc           context.CancelFunc
	mtx                  sync.Mutex
	markCompleteOnce     sync.Once
	combineWg            sync.WaitGroup
	done                 chan struct{}
}

type StreamerConfig struct {
	K8sClient     kubernetes.Interface
	ArgoClient    *Client
	StorageClient blobstorage.StorageClient
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
		done:          make(chan struct{}),
	}, nil
}

func (s *LogsStreamer) Start(ctx context.Context) (chan Log, error) {
	workflow, err := s.argoClient.GetWorkflow(ctx, s.namespace, s.workflowName)
	if err != nil {
		return nil, errors.New("workflow not found")
	}

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	s.cancelFunc = cancelFunc

	s.addStreams(workflow)
	s.combineStreams()
	go s.startStreaming(cancelCtx, workflow)

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
		s.combineWg.Add(1)
		go s.combineStream(stream)
	}

	go func() {
		s.combineWg.Wait()
		s.markCompleteOnce.Do(func() {
			s.markComplete()
		})
	}()
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
	defer s.combineWg.Done()
	for newLog := range stream {
		if newLog.Step != "" && newLog.Lines != "" {
			log.Debug().Msgf("New log: %s %s", newLog.Step, newLog.Lines)
			select {
			case s.combinedStream <- newLog:
			case <-s.done:
				return
			}
		}
	}

	log.Debug().Msgf("Closed stream. Num total: %d Num closed: %d", s.NumNodes(), s.NumDone())
}

func (s *LogsStreamer) waitForTail(c context.Context, template v1alpha1.Template, workflow *v1alpha1.Workflow) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Done():
			log.Debug().Msgf("Context cancelled, stopping waitForTail for template %s", template.Name)
			s.markStreamDone(template.Name)
			return
		case <-s.done:
			log.Debug().Msgf("LogsStreamer closed, stopping waitForTail for template %s", template.Name)
			s.markStreamDone(template.Name)
			return
		case <-ticker.C:
			s.mtx.Lock()
			closed := s.closed
			s.mtx.Unlock()

			if closed {
				log.Debug().Msgf("LogsStreamer closed, stopping waitForTail for template %s", template.Name)
				s.markStreamDone(template.Name)
				return
			}

			if _, err := s.tail(c, template.Name, workflow); err == nil {
				return
			}
		}
	}
}

func (s *LogsStreamer) tail(ctx context.Context, templateName string, workflow *v1alpha1.Workflow) (chan Log, error) {
	if workflow == nil {
		return nil, errors.New("workflow may not be nil")
	}

	if templateName == "" {
		return nil, errors.New("templateName may not be empty")
	}

	node, err := s.argoClient.GetNode(ctx, workflow, templateName)
	if err != nil {
		return nil, errors.Join(err, errors.New("unable to find node from template name"))
	}
	if node == nil {
		return nil, errors.New("node is nil after initial GetNode call")
	}

	containerName := "main"
	podName := FormatPodName(node)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Debug().Msgf("Timeout waiting for node %s to be ready", templateName)
			return s.getStream(templateName), errors.New("timeout waiting for node")
		case <-ctx.Done():
			log.Debug().Msgf("Context cancelled for node %s", templateName)
			return s.getStream(templateName), ctx.Err()
		case <-ticker.C:
			s.mtx.Lock()
			closed := s.closed
			s.mtx.Unlock()
			
			if closed {
				log.Debug().Msgf("LogsStreamer closed, stopping tail for node %s", templateName)
				return s.getStream(templateName), errors.New("streamer closed")
			}

			if node == nil {
				log.Debug().Msgf("Node %s is nil before GetNode call, retrying", templateName)
				continue
			}
			node, err = s.argoClient.GetNode(ctx, workflow, node.TemplateName)
			if err != nil {
				log.Debug().Err(err).Msgf("Error getting node %s, retrying", templateName)
				continue
			}
			if node == nil {
				log.Debug().Msgf("Node %s is nil, retrying", templateName)
				continue
			}
			if node.Pending() {
				continue
			}
		}

		readCloser, err := s.k8sClient.CoreV1().
			Pods(s.namespace).
			GetLogs(podName, &corev1.PodLogOptions{
				Container: containerName,
				Follow:    true,
			}).Stream(ctx)
		if err != nil {
			log.Debug().Err(err).Msgf("unable to get logs for pod %s", podName)
			continue
		}

		if node == nil {
			log.Debug().Msgf("Node %s is nil after select loop, cannot start log reading", templateName)
			readCloser.Close()
			return s.getStream(templateName), errors.New("node is nil after select loop")
		}
		
		nodeTemplateName := node.TemplateName
		go s.readLogsForNode(ctx, nodeTemplateName, readCloser)
		
		return s.getStream(nodeTemplateName), nil
	}
}

func (s *LogsStreamer) streamArchive(ctx context.Context, nodeName string) {
	archiveReadCloser, err := s.GetArchivedLogsForTemplate(ctx, nodeName)
	if err != nil {
		s.markStreamDone(nodeName)
		return
	}

	go s.readLogsForNode(ctx, nodeName, archiveReadCloser)
}

func (s *LogsStreamer) readLogsForNode(ctx context.Context, nodeName string, readCloser io.ReadCloser) {
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
		s.mtx.Lock()
		closed := s.closed
		s.mtx.Unlock()
		
		if closed {
			log.Debug().Msgf("LogsStreamer closed, stopping log reading for node %s", nodeName)
			s.markStreamDone(nodeName)
			break
		}

		select {
		case <-ctx.Done():
			log.Debug().Msgf("Context cancelled, stopping log reading for node %s", nodeName)
			s.markStreamDone(nodeName)
			return
		default:
		}

		buf := new(bytes.Buffer)
		if written, err := io.Copy(buf, readCloser); err != nil {
			log.Debug().Err(err).Msgf("unable to copy logs for %s", nodeName)
			s.markStreamDone(nodeName)
			break

		} else if written == 0 {

			if !s.nodeIsDone(ctx, nodeName) {
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
				select {
				case stream <- Log{
					Step:  nodeName,
					Lines: buf.String(),
				}:
					log.Debug().Msgf("streamed log for %s", nodeName)
				case <-s.done:
					log.Debug().Msgf("LogsStreamer done, stopping send for %s", nodeName)
					s.markStreamDone(nodeName)
					return
				}
			}

	}
}

func (s *LogsStreamer) GetArchivedLogsForTemplate(ctx context.Context, templateName string) (io.ReadCloser, error) {

	workflow, err := s.argoClient.GetWorkflow(ctx, s.namespace, s.workflowName)
	if err != nil || workflow == nil {
		return nil, errors.Join(err, errors.New("unable to get workflow"))
	}

	node, err := s.argoClient.GetNode(ctx, workflow, templateName)
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

	readCloser, err := s.storageClient.ReadFile(ctx, archive)
	if err != nil {
		return nil, errors.Join(err, errors.New("unable to read archived logs"))
	}

	return readCloser, nil
}

func (s *LogsStreamer) Close() error {
	s.mtx.Lock()
	if s.closed {
		s.mtx.Unlock()
		return nil
	}
	s.closed = true

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	close(s.done)

	// Force-close all individual stream channels to unblock combineStream
	// goroutines stuck on `range stream`. This is necessary for streams
	// where no producer goroutine was started (e.g., workflow in Pending
	// phase). The nil-guard prevents double-close if markStreamDone already
	// closed a stream. The done channel was closed first, so any active
	// producer will see it and exit before attempting to send.
	for name, stream := range s.streams {
		if stream != nil {
			close(stream)
			s.streams[name] = nil
		}
	}
	s.numDone = s.numNodes

	s.mtx.Unlock()

	// combinedStream is closed by the combineWg goroutine after all
	// combineStream goroutines finish their range loops.

	return nil
}

func hasLogs(template v1alpha1.Template) bool {
	return template.GetNodeType() == v1alpha1.NodeTypePod
}
