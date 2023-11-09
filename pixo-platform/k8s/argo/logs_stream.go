package argo

import (
	"bytes"
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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
	WorkflowName string
	Streams      map[string]chan Log
}

func NewLogsStreamer(k8sClient *base.Client, workflowName string) *LogsStreamer {
	return &LogsStreamer{
		k8sClient:    k8sClient,
		WorkflowName: workflowName,
		Streams:      make(map[string]chan Log),
	}
}

func (s *LogsStreamer) Tail(c context.Context, namespace string, node *v1alpha1.NodeStatus) (chan Log, error) {
	if node == nil {
		return nil, errors.New("node may not be nil")
	}

	if s.Streams[node.TemplateName] != nil {
		return s.Streams[node.TemplateName], nil
	}

	s.Streams[node.TemplateName] = make(chan Log, 100)

	containerName := "main"

	var ioStream io.ReadCloser
	var err error
	for {
		time.Sleep(1 * time.Second)
		ioStream, err = s.k8sClient.GetPodLogs(c, namespace, node.Name, containerName)
		if err != nil {
			if strings.Contains(err.Error(), "waiting to start") {
				continue
			} else {
				return nil, err
			}
		}

		break
	}

	go s.GetLogsForNode(node, ioStream)

	return s.Streams[node.TemplateName], nil
}

func (s *LogsStreamer) GetLogsForNode(node *v1alpha1.NodeStatus, ioStream io.ReadCloser) {
	if ioStream == nil || s.Streams[node.TemplateName] == nil {
		return
	}
	defer ioStream.Close()

	for {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, ioStream); err != nil {
			break
		}

		s.Streams[node.TemplateName] <- Log{
			Step:  node.Name,
			Lines: buf.String(),
		}
	}
}
