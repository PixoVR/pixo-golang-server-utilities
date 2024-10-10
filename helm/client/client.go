package helm

import (
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

type Chart struct {
	Name        string
	RepoURL     string
	Version     string
	ReleaseName string
	Namespace   string
}

type Client struct {
	config       ClientConfig
	actionConfig *action.Configuration
}

type ClientConfig struct {
	Namespace       string
	Driver          string
	ChartsDirectory string
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Namespace == "" {
		config.Namespace = os.Getenv("NAMESPACE")
	}

	if config.Driver == "" {
		config.Driver = os.Getenv("HELM_DRIVER")
	}

	actionConfig := new(action.Configuration)
	options := &genericclioptions.ConfigFlags{Namespace: &config.Namespace}

	if err := actionConfig.Init(options, config.Namespace, config.Driver, log.Printf); err != nil {
		log.Error().Err(err).Msgf("Failed to initialize helm provider")
		return nil, err
	}

	return &Client{
		actionConfig: actionConfig,
		config:       config,
	}, nil
}

func (c *Client) Namespace() string {
	return c.config.Namespace
}
