//go:build !release

package helm

import "helm.sh/helm/v3/pkg/chart"

var _ Client = &MockClient{}

type MockClient struct {
	NamespaceReturns string

	CalledInstallChartWith [][]interface{}
	InstallError           error

	CalledUpgradeChartWith [][]interface{}
	UpgradeError           error

	CalledExistsWith []Chart
	ExistsReturns    bool
	ExistsError      error

	CalledDownloadChartWith []string
	DownloadReturns         string
	DownloadError           error

	CalledLoadChartWith []Chart
	LoadChartReturns    *chart.Chart
	LoadError           error

	CalledUninstallChartWith []Chart
	UninstallError           error
}

func (c *MockClient) Namespace() string {
	return c.NamespaceReturns
}

func (c *MockClient) Install(chart Chart, values map[string]interface{}) error {
	c.CalledInstallChartWith = append(c.CalledInstallChartWith, []interface{}{chart, values})
	return c.InstallError
}

func (c *MockClient) Upgrade(chart Chart, values map[string]interface{}) error {
	c.CalledUpgradeChartWith = append(c.CalledUpgradeChartWith, []interface{}{chart, values})
	return c.UpgradeError
}

func (c *MockClient) Exists(chart Chart) (bool, error) {
	c.CalledExistsWith = append(c.CalledExistsWith, chart)
	if c.ExistsError != nil {
		return false, c.ExistsError
	}
	return c.ExistsReturns, nil
}

func (c *MockClient) DownloadChart(chartURL string) (string, error) {
	c.CalledDownloadChartWith = append(c.CalledDownloadChartWith, chartURL)
	if c.DownloadError != nil {
		return "", c.DownloadError
	}
	return c.DownloadReturns, nil
}

func (c *MockClient) LoadChart(chart Chart) (*chart.Chart, error) {
	c.CalledLoadChartWith = append(c.CalledLoadChartWith, chart)
	if c.LoadError != nil {
		return nil, c.LoadError
	}
	return c.LoadChartReturns, nil
}

func (c *MockClient) Uninstall(chart Chart) error {
	c.CalledUninstallChartWith = append(c.CalledUninstallChartWith, chart)
	return c.UninstallError
}
