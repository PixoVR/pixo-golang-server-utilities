package helm_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/helm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm", Ordered, func() {

	var (
		helmClient     helm.Client
		chart          helm.Chart
		sampleChartURL = "https://github.com/PixoVR/helm-charts/releases/download/multiplayer-build-trigger-0.0.2/multiplayer-build-trigger-0.0.2.tgz"
	)

	BeforeEach(func() {
		chart = helm.Chart{
			RepoURL:     helm.PixoRepoURL,
			Name:        "multiplayer-module",
			Namespace:   "dev-multiplayer",
			Version:     "0.0.23",
			ReleaseName: "helm-test",
		}

		var err error
		helmClient, err = helm.NewClient(helm.ClientConfig{
			ChartsDirectory: "/tmp",
			Namespace:       "dev-multiplayer",
			Driver:          "",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(helmClient).NotTo(BeNil())
	})

	It("can download a chart", func() {
		filepath, err := helmClient.DownloadChart(sampleChartURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(filepath).To(ContainSubstring("multiplayer-build-trigger-0.0.2.tgz"))
	})

	It("can load a chart", func() {
		chart, err := helmClient.LoadChart(chart)
		Expect(err).NotTo(HaveOccurred())
		Expect(chart).NotTo(BeNil())
		Expect(chart.Name()).To(Equal("multiplayer-module"))
	})

	It("can install a chart", func() {
		values := map[string]interface{}{
			"app_project_id": "pixo-dev",
		}
		err := helmClient.Install(chart, values)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can upgrade a chart", func() {
		values := map[string]interface{}{}
		err := helmClient.Upgrade(chart, values)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can tell if a chart is installed", func() {
		installed, err := helmClient.Exists(chart)
		Expect(err).NotTo(HaveOccurred())
		Expect(installed).To(BeTrue())

		nonexistentChart := helm.Chart{
			RepoURL:     chart.RepoURL,
			Name:        chart.Name,
			ReleaseName: "nonexistent-chart",
		}
		installed, err = helmClient.Exists(nonexistentChart)
		Expect(err).NotTo(HaveOccurred())
		Expect(installed).To(BeFalse())
	})

	It("can uninstall a chart", func() {
		err := helmClient.Uninstall(chart)
		Expect(err).NotTo(HaveOccurred())
	})

})
