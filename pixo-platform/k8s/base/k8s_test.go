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

})
