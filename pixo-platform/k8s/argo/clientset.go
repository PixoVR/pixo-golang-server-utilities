package argo

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

type Client struct {
	*versioned.Clientset
	Namespace string
}

func NewArgoClient(namespace string) (*Client, error) {
	kubeconfig, err := base.GetConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create argo client")
		return nil, err
	}

	clientset, err := getArgoClientsetFromConfig(kubeconfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create argo client")
		return nil, err
	}

	return &Client{
		Clientset: clientset,
		Namespace: namespace,
	}, nil
}

func getArgoClientsetFromConfig(config *rest.Config) (*versioned.Clientset, error) {
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create argo client")
		return nil, err
	}

	return clientset, nil
}
