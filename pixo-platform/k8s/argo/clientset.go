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

func NewArgoClient(namespace string) *Client {
	kubeconfig := base.GetKubeConfig()
	return &Client{
		Clientset: getArgoClientsetFromConfig(kubeconfig),
		Namespace: namespace,
	}
}

func getArgoClientsetFromConfig(config *rest.Config) *versioned.Clientset {
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create argo client")
	}

	return clientset
}
