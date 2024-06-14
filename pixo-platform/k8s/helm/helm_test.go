package helm_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/helm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm", Ordered, func() {

	var (
		helmClient     helm.Client
		chart          helm.Chart
		sampleChartURL = "https://github.com/PixoVR/helm-charts/releases/download/multiplayer-build-trigger-0.0.2/multiplayer-build-trigger-0.0.2.tgz"
		namespace      = config.GetEnvOrReturn("NAMESPACE", "test")
	)

	BeforeEach(func() {
		chart = helm.Chart{
			RepoURL:     helm.PixoRepoURL,
			Name:        "multiplayer-module",
			Namespace:   namespace,
			Version:     "0.0.23",
			ReleaseName: "helm-test",
		}

		var err error
		helmClient, err = helm.NewClient(helm.ClientConfig{
			ChartsDirectory: "/tmp",
			Namespace:       namespace,
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
			"create_infra":   false,
		}
		Expect(helmClient.Install(chart, values)).To(Succeed())
	})

	It("can upgrade a chart", func() {
		values := map[string]interface{}{}
		Expect(helmClient.Upgrade(chart, values)).To(Succeed())
	})

	It("can tell if a chart is installed", func() {
		installed, err := helmClient.Exists(chart)
		Expect(err).NotTo(HaveOccurred())
		Expect(installed).To(BeTrue())

		installed, err = helmClient.Exists(helm.Chart{
			RepoURL:     chart.RepoURL,
			Name:        chart.Name,
			ReleaseName: "nonexistent-chart",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(installed).To(BeFalse())
	})

	It("can uninstall a chart", func() {
		Expect(helmClient.Uninstall(chart)).To(Succeed())

		installed, err := helmClient.Exists(chart)
		Expect(err).NotTo(HaveOccurred())
		Expect(installed).To(BeFalse())
	})

})
