package agones

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

type Client struct {
	*versioned.Clientset
	client     rest.Interface
	BaseClient base.Client
}

func NewInClusterAgonesClient(baseClient base.Client) (*Client, error) {
	kubeconfig, err := base.GetConfigUsingInCluster()
	if err != nil {
		return nil, err
	}

	clientset, err := getAgonesClientsetFromConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset: clientset,
	}, nil
}

func NewLocalAgonesClient(baseClient base.Client) (*Client, error) {
	kubeconfig, err := base.GetConfigUsingKubeconfig()
	if err != nil {
		return nil, err
	}

	clientset, err := getAgonesClientsetFromConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset:  clientset,
		BaseClient: baseClient,
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
