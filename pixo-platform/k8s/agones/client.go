package agones

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/rs/zerolog/log"
)

type Client struct {
	*versioned.Clientset
	BaseClient base.Client
}

func NewInClusterAgonesClient(baseClient base.Client) (Client, error) {
	kubeconfig, err := base.GetConfigUsingInCluster()
	if err != nil {
		return Client{}, err
	}

	clientset, err := versioned.NewForConfig(kubeconfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create agones client")
		return Client{}, err
	}

	return Client{
		Clientset:  clientset,
		BaseClient: baseClient,
	}, nil
}

func NewLocalAgonesClient(baseClient base.Client) (Client, error) {
	kubeconfig, err := base.GetConfigUsingKubeconfig()
	if err != nil {
		return Client{}, err
	}

	clientset, err := versioned.NewForConfig(kubeconfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create agones client")
		return Client{}, err
	}

	return Client{
		Clientset:  clientset,
		BaseClient: baseClient,
	}, nil
}
