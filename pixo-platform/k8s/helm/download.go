package helm

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path"
)

func (c Client) loadChart(repoURL, chartName, chartVersion string) (*chart.Chart, error) {
	getterProviders := getter.Providers{
		{Schemes: []string{"http", "https"}, New: getter.NewHTTPGetter},
	}

	chartURL, err := repo.FindChartInRepoURL(repoURL, chartName, chartVersion, "", "", "", getterProviders)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to find chart in repo")
		return nil, err
	}

	filepath, err := c.DownloadChart(chartURL)
	if err != nil {
		return nil, err
	}

	helmChart, err := loader.Load(filepath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to load chart")
		return nil, err
	}

	return helmChart, nil
}

func (c Client) DownloadChart(chartURL string) (string, error) {
	dest := config.GetEnvOrReturn("TMP_DIR", "/tmp")

	client, err := getter.NewHTTPGetter()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create http getter")
		return "", err
	}

	fileBuffer, err := client.Get(chartURL)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to download chart")
		return "", err
	}

	if err = os.MkdirAll(dest, 0755); err != nil {
		log.Error().Err(err).Msgf("Failed to directory")
	}

	filepath := path.Join(dest, path.Base(chartURL))

	if err = os.WriteFile(filepath, fileBuffer.Bytes(), 0644); err != nil {
		log.Error().Err(err).Msgf("Failed to write file")
		return "", err
	}

	return filepath, nil
}
