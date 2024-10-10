package helm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

var (
	namespace string
)

func TestHelm(t *testing.T) {
	RegisterFailHandler(Fail)

	var ok bool
	namespace, ok = os.LookupEnv("NAMESPACE")
	if !ok {
		namespace = "test"
	}

	RunSpecs(t, "Helm Suite")
}
