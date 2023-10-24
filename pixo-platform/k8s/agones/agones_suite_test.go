package agones_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAgones(t *testing.T) {
	RegisterFailHandler(Fail)
	_ = os.Setenv("IS_LOCAL", "true")
	RunSpecs(t, "Agones Client Suite")
}

var (
	namespace    = "dev-multiplayer"
	agonesClient *agones.Client
)

var _ = BeforeSuite(func() {
	baseClient, err := base.NewLocalClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(baseClient).To(Not(BeNil()))

	agonesClient, err = agones.NewLocalAgonesClient(*baseClient)
	Expect(err).NotTo(HaveOccurred())
	Expect(agonesClient).To(Not(BeNil()))
})
