package base

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	kubernetes.Interface
}

func NewInClusterK8sClient() (Client, error) {
	kubeconfig, err := GetConfigUsingInCluster()
	if err != nil {
		return Client{}, err
	}

	clientset, err := getClientsetFromConfig(kubeconfig)
	if err != nil {
		return Client{}, err
	}

	return Client{clientset}, nil
}

func NewLocalClient() (Client, error) {
	kubeconfig, err := GetConfigUsingKubeconfig()
	if err != nil {
		return Client{}, err
	}

	clientset, err := getClientsetFromConfig(kubeconfig)
	if err != nil {
		return Client{}, err
	}

	return Client{clientset}, nil
}

func getClientsetFromConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create K8s clientset")
		return nil, err
	}

	return clientset, nil
}
