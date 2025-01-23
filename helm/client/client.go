package helm

import (
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

type Client interface {
	Namespace() string
	Install(chart Chart, values map[string]interface{}) error
	Upgrade(chart Chart, values map[string]interface{}) error
	Exists(chart Chart) (bool, error)
	Uninstall(chart Chart) error
	DownloadChart(chartURL string) (string, error)
	LoadChart(chart Chart) (*chart.Chart, error)
}

type Chart struct {
	Name        string
	RepoURL     string
	Version     string
	ReleaseName string
	Namespace   string
}

type clientImpl struct {
	config       ClientConfig
	actionConfig *action.Configuration
}

type ClientConfig struct {
	Namespace       string
	Driver          string
	ChartsDirectory string
}

func NewClient(config ClientConfig) (Client, error) {
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

	return &clientImpl{
		actionConfig: actionConfig,
		config:       config,
	}, nil
}

func (c clientImpl) Namespace() string {
	return c.config.Namespace
}
