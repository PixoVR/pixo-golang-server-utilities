package auth_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
)

var _ = Describe("Auth", func() {

	var (
		originalLifecycle string
	)

	BeforeEach(func() {
		originalLifecycle = os.Getenv("SECRET_KEY")
		os.Setenv("SECRET_KEY", "test")
	})

	AfterEach(func() {
		os.Setenv("SECRET_KEY", originalLifecycle)
	})

	It("can determine if a secret key is valid", func() {
		isSecretKey := auth.IsValidSecretKey("test")
		Expect(isSecretKey).To(BeTrue())
	})

	It("can determine if a secret key is invalid", func() {
		isNotSecretKey := auth.IsValidSecretKey("")
		Expect(isNotSecretKey).To(BeFalse())
	})

	Context("middleware test", func() {

		var (
			endpoint = "/test"
			host     = "127.0.0.1:8000"
			ip       = "127.0.0.1"

			engine *gin.Engine
			w      *httptest.ResponseRecorder
			req    *http.Request
		)

		BeforeEach(func() {
			engine = gin.Default()
			engine.Use(auth.HostMiddleware())
			w = httptest.NewRecorder()

			req, _ = http.NewRequest(http.MethodGet, endpoint, nil)
			req.RemoteAddr = host
		})

		It("will return the context with the gin context", func() {
			engine.GET(endpoint, func(c *gin.Context) {
				Expect(config.GetGinContext(c)).NotTo(BeNil())
			})

			engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("will return the context with the ip address", func() {
			engine.GET(endpoint, func(c *gin.Context) {
				Expect(config.GetIPAddress(c)).To(Equal(ip))
			})

			engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})

	})

})
