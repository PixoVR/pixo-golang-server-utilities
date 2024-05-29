package base

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func GetConfigUsingKubeconfig() (*rest.Config, error) {
	configPath, exists := os.LookupEnv("KUBECONFIG")
	if !exists {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/workspace"
		}
		configPath = filepath.Join(home, ".kube", "config")
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to build K8s config using kubeconfig")
		return nil, err
	}

	return kubeconfig, nil
}

func GetConfigUsingInCluster() (*rest.Config, error) {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to create K8s kubeconfig")
		return nil, err
	}

	return kubeconfig, nil
}
