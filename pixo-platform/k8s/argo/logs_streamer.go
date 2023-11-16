package argo

import (
	"bytes"
	"context"
	"errors"
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
	k8sClient    base.Client
	argoClient   Client
	WorkflowName string
	namespace    string
	streams      map[string]chan Log
	numStreams   int
	numClosed    int
	mtx          sync.Mutex
}

func NewLogsStreamer(k8sClient base.Client, argoClient Client, namespace, workflowName string) (*LogsStreamer, error) {
	if namespace == "" {
		return nil, errors.New("namespace may not be empty")
	}

	if workflowName == "" {
		return nil, errors.New("workflowName may not be empty")
	}

	return &LogsStreamer{
		k8sClient:    k8sClient,
		argoClient:   argoClient,
		WorkflowName: workflowName,
		namespace:    namespace,
		streams:      make(map[string]chan Log),
	}, nil
}

func (s *LogsStreamer) GetLogsStream() chan Log {
	combinedStream := make(chan Log)

	go func() {
		streams, err := s.TailAll(context.Background())
		if err != nil {
			return
		}

		for _, stream := range streams {
			go s.CombineStream(combinedStream, stream)
		}
	}()

	return combinedStream
}

func (s *LogsStreamer) CombineStream(combinedStream, stream chan Log) {
	for !IsClosed(stream) {
		newLog := <-stream
		log.Debug().Msgf("New log: %s %s", newLog.Step, newLog.Lines)
		combinedStream <- newLog
	}

	log.Debug().Msgf("Closed stream. Num total: %d Num closed: %d", s.NumTotalStreams(), s.NumClosed())

	if !IsClosed(combinedStream) && s.IsDone() {
		log.Debug().Msg("All streams are closed. Closing combined stream")
		close(combinedStream)
	}
}

func (s *LogsStreamer) TailAll(c context.Context) (map[string]chan Log, error) {
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

		stream := s.getStream(template.Name)

		if stream == nil {
			s.addStream(template.Name)
		}

		go s.waitForTail(c, template, workflow)
	}

	return s.streams, nil
}

func (s *LogsStreamer) waitForTail(c context.Context, template v1alpha1.Template, workflow *v1alpha1.Workflow) {
	for {
		time.Sleep(1 * time.Second)

		if _, err := s.Tail(c, template.Name, workflow); err == nil {
			break
		}
	}
}

func (s *LogsStreamer) Tail(c context.Context, templateName string, workflow *v1alpha1.Workflow) (chan Log, error) {
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

	go s.StreamLogsForNode(node, ioStream)

	log.Debug().Msgf("started tailing logs for node %s", node.TemplateName)

	return s.getStream(node.TemplateName), nil
}

func (s *LogsStreamer) StreamLogsForNode(node *v1alpha1.NodeStatus, ioStream io.ReadCloser) {
	if ioStream == nil || node == nil {
		return
	}
	defer ioStream.Close()

	stream := s.getStream(node.TemplateName)
	if stream == nil {
		log.Debug().Msgf("stream for %s is nil", node.TemplateName)
		return
	}

	log.Debug().Msgf("started streaming logs for %s", node.TemplateName)
	for {
		buf := new(bytes.Buffer)

		if _, err := io.Copy(buf, ioStream); err != nil {
			log.Debug().Err(err).Msgf("unable to copy logs for %s", node.TemplateName)
			s.closeStream(node.TemplateName)
			break

		} else if buf.Len() == 0 {
			if !s.nodeIsDone(node) {
				log.Debug().Msgf("no logs for %s", node.TemplateName)
				time.Sleep(1 * time.Second)
				continue
			} else {
				s.closeStream(node.TemplateName)
				break
			}
		}

		stream <- Log{
			Step:  node.DisplayName,
			Lines: buf.String(),
		}

		log.Debug().Msgf("streamed log for %s", node.TemplateName)
	}
}

func (s *LogsStreamer) ReadFromStream(name string) *Log {
	stream := s.getStream(name)
	if IsClosed(stream) {
		return nil
	}

	newLog := <-stream
	return &newLog
}

func (s *LogsStreamer) IsDone() bool {
	isDone := s.NumTotalStreams() == s.NumClosed()

	log.Debug().Msgf("is done %t", isDone)
	return isDone
}

func (s *LogsStreamer) NumTotalStreams() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	log.Debug().Msgf("num streams %d", s.numStreams)
	return s.numStreams
}

func (s *LogsStreamer) NumClosed() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	log.Debug().Msgf("num closed %d", s.numClosed)
	return s.numClosed
}

func (s *LogsStreamer) addStream(name string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.streams[name] != nil && !IsClosed(s.streams[name]) {
		return false
	}

	s.streams[name] = make(chan Log)
	s.numStreams++

	log.Debug().Msgf("added stream %s", name)
	return true
}

func (s *LogsStreamer) getStream(name string) chan Log {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	log.Debug().Msgf("getting stream %s", name)
	return s.streams[name]
}

func (s *LogsStreamer) closeStream(name string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.streams, name)
	s.numClosed++

	log.Debug().Msgf("closed stream %s", name)

}

func (s *LogsStreamer) nodeIsDone(node *v1alpha1.NodeStatus) bool {
	if node == nil {
		return false
	}

	workflow, err := s.argoClient.GetWorkflow(s.namespace, s.WorkflowName)
	if err != nil {
		return false
	}

	newNode, err := s.argoClient.GetNode(workflow, node.TemplateName)
	if err != nil || newNode == nil {
		return false
	}

	return newNode.Phase == v1alpha1.NodeSucceeded || newNode.Phase == v1alpha1.NodeFailed
}
