package base_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("K8s", func() {

	var (
		baseClient *base.Client
	)

	BeforeEach(func() {
		baseClient, _ = base.NewLocalClient()
		Expect(baseClient).NotTo(BeNil())
	})

	It("can get a pod by name", func() {
		pods, err := baseClient.GetPods("dev-multiplayer")
		Expect(err).NotTo(HaveOccurred())
		Expect(pods).NotTo(BeNil())

		for _, p := range pods.Items {
			Expect(p).NotTo(BeNil())
			Expect(p.Name).NotTo(BeNil())
			pod, err := baseClient.GetPod(p.Namespace, p.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod).NotTo(BeNil())
			Expect(pod.Name).To(Equal(p.Name))
		}
	})

})
