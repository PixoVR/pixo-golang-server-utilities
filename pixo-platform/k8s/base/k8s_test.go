package base_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("K8s", func() {

	var (
		baseClient base.Client
	)

	BeforeEach(func() {
		var err error
		baseClient, err = base.NewLocalClient()
		Expect(baseClient).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())
	})

	It("can get a pod by name", func() {
		pods, err := baseClient.GetPods(context.Background(), namespace)
		Expect(err).NotTo(HaveOccurred())
		Expect(pods).NotTo(BeNil())

		for _, p := range pods.Items {
			Expect(p).NotTo(BeNil())
			Expect(p.Name).NotTo(BeNil())
			pod, err := baseClient.GetPod(context.Background(), p.Namespace, p.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod).NotTo(BeNil())
			Expect(pod.Name).To(Equal(p.Name))
		}
	})

})
