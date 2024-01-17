package engine_test

import (
	"os"
	"testing"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const fakeKey = "fake-key"

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	config.LoadEnvVars()
	RunSpecs(t, "API Suite")
}

var _ = BeforeSuite(func() {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", fakeKey)
})
