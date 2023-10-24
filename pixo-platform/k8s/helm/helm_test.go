package helm_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/helm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm", func() {

	var (
		helmClient *helm.Client
	)

	BeforeEach(func() {
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
		chartURL := "https://github.com/PixoVR/helm-charts/releases/download/multiplayer-build-trigger-0.0.2/multiplayer-build-trigger-0.0.2.tgz"
		filepath, err := helmClient.DownloadChart(chartURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(filepath).To(ContainSubstring("multiplayer-build-trigger-0.0.2.tgz"))
	})

})
