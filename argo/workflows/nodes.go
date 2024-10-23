package workflows

import (
	"context"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

func (s *LogsStreamer) nodeIsDone(ctx context.Context, nodeName string) bool {
	if nodeName == "" {
		return false
	}

	workflow, err := s.argoClient.GetWorkflow(ctx, s.namespace, s.workflowName)
	if err != nil {
		return false
	}

	newNode, err := s.argoClient.GetNode(ctx, workflow, nodeName)
	if err != nil || newNode == nil {
		return false
	}

	return newNode.Phase == v1alpha1.NodeSucceeded || newNode.Phase == v1alpha1.NodeFailed
}
