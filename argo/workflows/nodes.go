package workflows

import (
	"context"
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

	return newNode.Fulfilled()
}

func (s *LogsStreamer) workflowIsDone(ctx context.Context) bool {
	workflow, err := s.argoClient.GetWorkflow(ctx, s.namespace, s.workflowName)
	if err != nil || workflow == nil {
		return false
	}

	return workflow.Status.Phase.Completed()
}
