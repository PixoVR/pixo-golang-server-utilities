package agones

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

type Client struct {
	*versioned.Clientset
	Namespace string
}

func NewAgonesClient(namespace string) (*Client, error) {
	kubeconfig, err := base.GetConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := getAgonesClientsetFromConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset: clientset,
		Namespace: namespace,
	}, nil
}

func getAgonesClientsetFromConfig(config *rest.Config) (*versioned.Clientset, error) {
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create agones client")
		return nil, err
	}

	return clientset, nil
}
