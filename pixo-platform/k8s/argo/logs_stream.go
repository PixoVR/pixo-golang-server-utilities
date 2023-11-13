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
	"time"
)

type Log struct {
	Step  string `json:"step"`
	Lines string `json:"line"`
}

type LogsStreamer struct {
	k8sClient    *base.Client
	argoClient   *Client
	WorkflowName string
	streams      map[string]chan Log
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
		if template.GetNodeType() != v1alpha1.NodeTypePod {
			continue
		}

		failureThreshold := 30
		for {
			time.Sleep(1 * time.Second)
			if _, err = s.Tail(c, namespace, template.Name, workflow); err == nil || failureThreshold == 0 {
				break
			}

			failureThreshold--
		}
	}

	return s.streams, nil
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

	if s.streams[node.TemplateName] != nil {
		return s.streams[node.TemplateName], nil
	}

	s.streams[node.TemplateName] = make(chan Log)

	containerName := "main"
	podName := FormatPodName(node)

	var ioStream io.ReadCloser

	failureThreshold := 30
	for {
		time.Sleep(1 * time.Second)

		node, err = s.argoClient.GetNode(workflow, node.TemplateName)
		if err != nil {
			if failureThreshold == 0 {
				return nil, err
			}

			failureThreshold--
		}

		ioStream, err = s.k8sClient.GetPodLogs(c, namespace, podName, containerName)
		if err != nil {
			if strings.Contains(err.Error(), "waiting to start") {
				continue
			} else if failureThreshold == 0 {
				return nil, err
			}

			failureThreshold--
		}

		break
	}

	go s.StreamLogsForNode(node, ioStream)

	return s.streams[node.TemplateName], nil
}

func (s *LogsStreamer) StreamLogsForNode(node *v1alpha1.NodeStatus, ioStream io.ReadCloser) {
	if ioStream == nil || node == nil || s.streams[node.TemplateName] == nil {
		return
	}
	defer ioStream.Close()

	for {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, ioStream); err != nil {
			break
		}

		s.streams[node.TemplateName] <- Log{
			Step:  node.Name,
			Lines: buf.String(),
		}
	}
}

func (s *LogsStreamer) ReadFromStream(name string) *Log {
	if s.streams[name] == nil {
		return nil
	}

	newLog := <-s.streams[name]
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
