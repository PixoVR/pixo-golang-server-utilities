package agones_test

import (
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
