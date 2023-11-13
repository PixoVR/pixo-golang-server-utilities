package argo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
	"sync"
	"time"
)

type Log struct {
	Step  string `json:"step"`
	Lines string `json:"lines"`
}

type LogsStreamer struct {
	k8sClient    *base.Client
	argoClient   *Client
	WorkflowName string
	streams      map[string]chan Log
	mtx          sync.Mutex
}

func NewLogsStreamer(k8sClient *base.Client, argoClient *Client, workflowName string) *LogsStreamer {
	return &LogsStreamer{
		k8sClient:    k8sClient,
		argoClient:   argoClient,
		WorkflowName: workflowName,
		streams:      make(map[string]chan Log),
	}
}

func (s *LogsStreamer) TailAll(c context.Context, namespace string, workflowName string) (map[string]chan Log, error) {
	if workflowName == "" {
		return nil, errors.New("workflowName may not be empty")
	}

	workflow, err := s.argoClient.GetWorkflow(namespace, workflowName)
	if err != nil {
		return nil, err
	}

	if workflow == nil {
		return nil, errors.New("workflow not found")
	}

	for _, template := range workflow.Spec.Templates {
		if template.GetNodeType() != v1alpha1.NodeTypePod || s.streams[template.Name] != nil {
			continue
		}

		s.mtx.Lock()
		stream := s.streams[template.Name]
		if stream == nil {
			s.streams[template.Name] = make(chan Log)
		}
		s.mtx.Unlock()

		go s.waitForTail(c, namespace, template, workflow)
	}

	return s.streams, nil
}

func (s *LogsStreamer) waitForTail(c context.Context, namespace string, template v1alpha1.Template, workflow *v1alpha1.Workflow) {
	for {
		time.Sleep(1 * time.Second)
		if _, err := s.Tail(c, namespace, template.Name, workflow); err == nil {
			break
		}
	}
}

func (s *LogsStreamer) Tail(c context.Context, namespace, templateName string, workflow *v1alpha1.Workflow) (chan Log, error) {
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
		if err != nil || node.Pending() {
			continue
		}

		ioStream, err = s.k8sClient.GetPodLogs(c, namespace, podName, containerName)
		if err != nil {
			log.Debug().Err(err).Msgf("unable to get logs for pod %s", podName)
			continue
		}

		break
	}

	go s.StreamLogsForNode(node, ioStream)

	log.Debug().Msgf("started tailing logs for node %s", node.TemplateName)

	s.mtx.Lock()
	stream := s.streams[node.TemplateName]
	s.mtx.Unlock()

	return stream, nil
}

func (s *LogsStreamer) StreamLogsForNode(node *v1alpha1.NodeStatus, ioStream io.ReadCloser) {
	if ioStream == nil || node == nil {
		return
	}

	s.mtx.Lock()
	stream := s.streams[node.TemplateName]
	s.mtx.Unlock()

	if stream == nil {
		return
	}

	defer ioStream.Close()

	log.Debug().Msgf("started streaming logs for %s", node.TemplateName)
	for {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, ioStream); err != nil {
			break
		}

		stream <- Log{
			Step:  node.Name,
			Lines: buf.String(),
		}
		log.Debug().Msgf("streamed log for %s", node.TemplateName)
	}
}

func (s *LogsStreamer) ReadFromStream(name string) *Log {
	s.mtx.Lock()
	stream := s.streams[name]
	s.mtx.Unlock()

	if stream == nil {
		return nil
	}

	newLog := <-stream
	return &newLog
}

func FormatPodName(node *v1alpha1.NodeStatus) string {
	if node == nil {
		return ""
	}

	nodeID := strings.Split(node.ID, "-")
	podName := fmt.Sprintf("%s-%s-%s", node.BoundaryID, node.TemplateName, nodeID[len(nodeID)-1])

	log.Debug().Msgf("Formatted pod name: %s", podName)

	return podName
}
