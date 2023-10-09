package base

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var (
	K8sClient = GetK8sClient()
	isLocal   = config.GetEnvOrReturn("IS_LOCAL", "false") == "true"
	Namespace = config.GetEnvOrReturn("NAMESPACE", "dev-multiplayer")
)

type Client struct {
	Clientset kubernetes.Interface
}

func GetK8sClient() *kubernetes.Clientset {
	if config.GetLifecycle() == "" {
		return nil
	}

	kubeconfig := GetKubeConfig()

	return getClientsetFromConfig(kubeconfig)
}

func GetKubeConfig() *rest.Config {
	var kubeconfig *rest.Config
	if isLocal {
		kubeconfig = GetConfigUsingKubeconfig()
	} else {
		kubeconfig = getConfigUsingInCluster()
	}

	return kubeconfig
}

func GetConfigUsingKubeconfig() *rest.Config {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	configPath := filepath.Join(home, ".kube", "config")

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create K8s kubeconfig")
	}

	return kubeconfig
}

func getConfigUsingInCluster() *rest.Config {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create K8s kubeconfig")
	}

	return kubeconfig
}

func getClientsetFromConfig(config *rest.Config) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create K8s clientset")
	}

	return clientset
}
