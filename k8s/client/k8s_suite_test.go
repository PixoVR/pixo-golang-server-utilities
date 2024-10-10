package client_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	namespace string
)

func TestK8s(t *testing.T) {
	RegisterFailHandler(Fail)

	_ = os.Setenv("IN_CLUSTER", "false")

	var ok bool
	namespace, ok = os.LookupEnv("NAMESPACE")
	if !ok {
		namespace = "default"
	}

	RunSpecs(t, "K8s Client Suite")
}
