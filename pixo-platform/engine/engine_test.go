package engine_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/engine"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("CustomEngine", func() {

	var (
		e *engine.CustomEngine
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		e = engine.NewEngine(engine.Config{})
		Expect(e).NotTo(BeNil())
		Expect(e.Port()).To(Equal(engine.DefaultPort))
		Expect(e.PortString()).To(Equal(fmt.Sprintf(":%d", engine.DefaultPort)))
		Expect(e.BasePath()).To(Equal(engine.DefaultBasePath))
	})

	It("can create an engine with defaults", func() {
		Expect(e).NotTo(BeNil())
		req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
		Expect(err).NotTo(HaveOccurred())

		e.ServeHTTP(w, req)

		Expect(err).To(BeNil())
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("can create an custom engine with a basic route and tracing", func() {
		config := engine.Config{
			Port:              8002,
			BasePath:          "/api/v2",
			Tracing:           true,
			CollectorEndpoint: "localhost:55678",
			InternalRoutes:    true,
			ExternalRoutes:    true,
		}
		customEngine := engine.NewEngine(config)
		Expect(customEngine).NotTo(BeNil())
		Expect(customEngine.Port()).To(Equal(config.Port))
		Expect(customEngine.BasePath()).To(Equal(config.BasePath))
		Expect(customEngine.PublicRouteGroup).NotTo(BeNil())
		Expect(customEngine.InternalRouteGroup).NotTo(BeNil())
		Expect(customEngine.ExternalRouteGroup).NotTo(BeNil())

		req, err := http.NewRequest(http.MethodGet, "/api/v2/health", nil)
		Expect(err).NotTo(HaveOccurred())

		customEngine.ServeHTTP(w, req)

		Expect(err).To(BeNil())
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("uses middleware to add a host to the context", func() {
		ip := "127.0.0.1"
		e.GET("/test", func(c *gin.Context) {
			Expect(config.GetIPAddress(c)).To(Equal(ip))
		})
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "127.0.0.1:8000"

		e.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK))
	})

})
