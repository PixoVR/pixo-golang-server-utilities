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

func NewAgonesClient(namespace string) *Client {
	kubeconfig := base.GetKubeConfig()
	return &Client{
		Clientset: getAgonesClientsetFromConfig(kubeconfig),
		Namespace: namespace,
	}
}

func getAgonesClientsetFromConfig(config *rest.Config) *versioned.Clientset {
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create agones client")
	}

	return clientset
}
