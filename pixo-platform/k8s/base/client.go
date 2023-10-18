package base

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	Clientset kubernetes.Interface
}

func NewInClusterK8sClient() (*kubernetes.Clientset, error) {
	kubeconfig, err := GetConfigUsingInCluster()
	if err != nil {
		return nil, err
	}

	return getClientsetFromConfig(kubeconfig)
}

func NewLocalK8sClient() (*kubernetes.Clientset, error) {
	kubeconfig, err := GetConfigUsingKubeconfig()
	if err != nil {
		return nil, err
	}

	return getClientsetFromConfig(kubeconfig)
}

func getClientsetFromConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create K8s clientset")
		return nil, err
	}

	return clientset, nil
}
