package base_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestK8s(t *testing.T) {
	RegisterFailHandler(Fail)
	_ = os.Setenv("IS_LOCAL", "true")
	RunSpecs(t, "K8s Client Suite")
}

var (
	namespace = config.GetEnvOrReturn("NAMESPACE", "test")
)
