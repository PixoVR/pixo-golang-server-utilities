package helm

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path"
)

func (c clientImpl) LoadChart(chart Chart) (*chart.Chart, error) {
	getterProviders := getter.Providers{
		{Schemes: []string{"http", "https"}, New: getter.NewHTTPGetter},
	}

	chartURL, err := repo.FindChartInRepoURL(chart.RepoURL, chart.Name, chart.Version, "", "", "", getterProviders)
	if err != nil {
		return nil, err
	}

	filepath, err := c.DownloadChart(chartURL)
	if err != nil {
		return nil, err
	}

	helmChart, err := loader.Load(filepath)
	if err != nil {
		return nil, err
	}

	return helmChart, nil
}

func (c clientImpl) DownloadChart(chartURL string) (string, error) {
	dest, ok := os.LookupEnv("TMP_DIR")
	if !ok {
		dest = "/tmp"
	}

	client, err := getter.NewHTTPGetter()
	if err != nil {
		return "", err
	}

	fileBuffer, err := client.Get(chartURL)
	if err != nil {
		return "", err
	}

	if err = os.MkdirAll(dest, 0755); err != nil {
		return "", err
	}

	filepath := path.Join(dest, path.Base(chartURL))

	if err = os.WriteFile(filepath, fileBuffer.Bytes(), 0644); err != nil {
		return "", err
	}

	return filepath, nil
}
