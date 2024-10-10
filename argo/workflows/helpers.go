package workflows

import (
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rs/zerolog/log"
	"strings"
)

func FormatPodName(node *v1alpha1.NodeStatus) string {
	if node == nil {
		log.Debug().Msg("Node is nil")
		return ""
	}

	nodeID := strings.Split(node.ID, "-")
	podName := fmt.Sprintf("%s-%s-%s", node.BoundaryID, node.TemplateName, nodeID[len(nodeID)-1])

	return podName
}
