package argo

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rs/zerolog/log"
	"strings"
)

func IsClosed(ch <-chan Log) bool {
	if ch == nil {
		return true
	}

	select {
	case <-ch:
		return true
	default:
		return false
	}
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
