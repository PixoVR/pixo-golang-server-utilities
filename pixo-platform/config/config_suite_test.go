package config_test

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var (
	originalLifecycle string
)

var _ = BeforeSuite(func() {
	gin.SetMode(gin.TestMode)
	originalLifecycle = os.Getenv("LIFECYCLE")
	os.Setenv("LIFECYCLE", "test")
})

var _ = AfterSuite(func() {
	os.Setenv("LIFECYCLE", originalLifecycle)
})
