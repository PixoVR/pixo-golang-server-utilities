package cdn_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCDN(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CDN Suite")
}
