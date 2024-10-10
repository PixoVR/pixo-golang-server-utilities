package client

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func NewLocalClient() (kubernetes.Interface, error) {
	kubeconfig, err := GetConfigUsingKubeconfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to build K8s clientset")
		return nil, err
	}

	return clientset, nil
}

func GetConfigUsingKubeconfig() (*rest.Config, error) {
	configPath, exists := os.LookupEnv("KUBECONFIG")
	if !exists {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/workspace"
		}
		configPath = filepath.Join(home, ".kube", "config")
	}

	log.Debug().Str("configPath", configPath).Msg("Using KUBECONFIG for K8s config")

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to build K8s config using kubeconfig")
		return nil, err
	}

	return kubeconfig, nil
}
