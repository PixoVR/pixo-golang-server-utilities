package client

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewInClusterClient() (kubernetes.Interface, error) {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to build K8s config using in-cluster config")
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
