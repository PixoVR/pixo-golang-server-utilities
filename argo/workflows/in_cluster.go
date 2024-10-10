package workflows

import (
	"github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

func NewInClusterArgoClient() (*Client, error) {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to build K8s config using in-cluster config")
		return nil, err
	}

	clientset, err := versioned.NewForConfig(kubeconfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create argo client")
		return nil, err
	}

	return &Client{Clientset: clientset}, nil
}
