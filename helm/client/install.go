package helm

import (
	"helm.sh/helm/v3/pkg/action"
)

func (c Client) Install(chart Chart, values map[string]interface{}) error {
	helmChart, err := c.LoadChart(chart)
	if err != nil {
		return err
	}

	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = chart.ReleaseName
	client.Namespace = chart.Namespace

	if _, err = client.Run(helmChart, values); err != nil {
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
		return err
	}

	return nil
}

func (c Client) Exists(chart Chart) (bool, error) {
	client := action.NewList(c.actionConfig)

	releases, err := client.Run()
	if err != nil {
		return false, err
	}

	for _, release := range releases {
		if release.Name == chart.ReleaseName {
			return true, nil
		}
	}

	return false, nil
}

func (c Client) Uninstall(chart Chart) error {
	client := action.NewUninstall(c.actionConfig)

	if _, err := client.Run(chart.ReleaseName); err != nil {
		return err
	}

	return nil
}
