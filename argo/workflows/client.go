package workflows

import (
	"github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
)

type Client struct {
	*versioned.Clientset
}
