package helm_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/helm/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("Helm", Ordered, func() {

	var (
		chart          helm.Chart
		helmClient     *helm.Client
		NginxRepoURL   = "https://kubernetes.github.io/ingress-nginx"
		sampleChartURL = "https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.11.2/ingress-nginx-4.11.2.tgz"
	)

	BeforeAll(func() {
		chart = helm.Chart{
			RepoURL:     NginxRepoURL,
			Name:        "ingress-nginx",
			Namespace:   namespace,
			Version:     "4.11.2",
			ReleaseName: fmt.Sprintf("nginx-test-%d", rand.Intn(1000000)),
		}

		_ = os.Setenv("NAMESPACE", namespace)
		_ = os.Setenv("HELM_DRIVER", "memory")

		var err error
		helmClient, err = helm.NewClient(helm.ClientConfig{
			ChartsDirectory: "/tmp",
			Driver:          "",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(helmClient).NotTo(BeNil())
	})

	It("can get the namespace", func() {
		Expect(helmClient.Namespace()).To(Equal(namespace))
	})

	Context("downloading", func() {

		It("returns an error if the chart URL is invalid", func() {
			_, err := helmClient.DownloadChart("invalid-url")
			Expect(err).To(HaveOccurred())
		})

		It("can download a chart", func() {
			filepath, err := helmClient.DownloadChart(sampleChartURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(filepath).To(ContainSubstring("ingress-nginx-4.11.2.tgz"))
		})

		It("can load a chart", func() {
			installedChart, err := helmClient.LoadChart(chart)
			Expect(err).NotTo(HaveOccurred())
			Expect(installedChart).NotTo(BeNil())
			Expect(installedChart.Name()).To(Equal(chart.Name))
		})

	})

	Context("installing", func() {

		BeforeAll(func() {
			values := map[string]interface{}{
				"controller": map[string]interface{}{
					"service": map[string]interface{}{
						"type": "NodePort",
					},
				},
			}
			Expect(helmClient.Install(chart, values)).To(Succeed())
		})

		AfterAll(func() {
			Expect(helmClient.Uninstall(chart)).To(Succeed())
		})

		It("can return an error upgrading a non-existent chart", func() {
			values := map[string]interface{}{}
			err := helmClient.Upgrade(helm.Chart{
				RepoURL:     chart.RepoURL,
				Name:        chart.Name,
				ReleaseName: "nonexistent-chart",
			}, values)
			Expect(err).To(HaveOccurred())
		})

		It("can upgrade a chart", func() {
			values := map[string]interface{}{}
			Expect(helmClient.Upgrade(chart, values)).To(Succeed())
		})

		It("can tell if a chart is installed", func() {
			isInstalled, err := helmClient.Exists(chart)
			Expect(err).NotTo(HaveOccurred())
			Expect(isInstalled).To(BeTrue())

			isInstalled, err = helmClient.Exists(helm.Chart{
				RepoURL:     chart.RepoURL,
				Name:        chart.Name,
				ReleaseName: "nonexistent-chart",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(isInstalled).To(BeFalse())
		})

		It("returns an error if trying to uninstall a non-existent chart", func() {
			err := helmClient.Uninstall(helm.Chart{
				RepoURL:     chart.RepoURL,
				Name:        chart.Name,
				ReleaseName: "nonexistent-chart",
			})
			Expect(err).To(HaveOccurred())
		})

	})

})
