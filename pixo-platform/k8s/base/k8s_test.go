package base_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("K8s", func() {

	It("can connect to your local cluster", func() {
		client := base.GetK8sClient()
		Expect(client).To(Not(BeNil()))
	})

})
