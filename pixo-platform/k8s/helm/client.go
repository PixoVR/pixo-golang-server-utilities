package helm

import (
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

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

func (c Client) Install(chart Chart, values map[string]interface{}) error {
	helmChart, err := c.LoadChart(chart)
	if err != nil {
		return err
	}

	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = chart.ReleaseName
	client.Namespace = chart.Namespace

	if _, err = client.Run(helmChart, values); err != nil {
		log.Error().Err(err).Msgf("Failed to install chart")
		return err
	}

	return nil
}

func (c Client) Upgrade(chart Chart, values map[string]interface{}) error {
	helmChart, err := c.LoadChart(chart)
	if err != nil {
		return err
	}

	client := action.NewUpgrade(c.actionConfig)
	client.Namespace = chart.Namespace

	if _, err = client.Run(chart.ReleaseName, helmChart, values); err != nil {
		log.Error().Err(err).Msgf("Failed to upgrade chart")
		return err
	}

	return nil
}

func (c Client) Uninstall(chart Chart) error {
	client := action.NewUninstall(c.actionConfig)

	if _, err := client.Run(chart.ReleaseName); err != nil {
		log.Error().Err(err).Msgf("Failed to uninstall chart")
		return err
	}

	return nil
}
