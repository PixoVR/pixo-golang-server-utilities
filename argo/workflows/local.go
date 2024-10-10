package workflows

import (
	"github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func NewLocalClient() (*Client, error) {
	kubeconfig, err := GetConfigUsingKubeconfig()
	if err != nil {
		return nil, err
	}

	clientset, err := versioned.NewForConfig(kubeconfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create argo client")
		return nil, err
	}
	return &Client{Clientset: clientset}, nil
}

func NewLocalBaseClient() (kubernetes.Interface, error) {
	kubeconfig, err := GetConfigUsingKubeconfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeconfig)
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

	return clientcmd.BuildConfigFromFlags("", configPath)
}
