package argo

import "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

func (s *LogsStreamer) nodeIsDone(nodeName string) bool {
	if nodeName == "" {
		return false
	}

	workflow, err := s.argoClient.GetWorkflow(s.namespace, s.WorkflowName)
	if err != nil {
		return false
	}

	newNode, err := s.argoClient.GetNode(workflow, nodeName)
	if err != nil || newNode == nil {
		return false
	}

	return newNode.Phase == v1alpha1.NodeSucceeded || newNode.Phase == v1alpha1.NodeFailed
}
