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
		internalEngine := e.Engine()
		Expect(e).NotTo(BeNil())
		req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
		Expect(err).NotTo(HaveOccurred())

		internalEngine.ServeHTTP(w, req)

		Expect(err).To(BeNil())
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(e.Engine()).NotTo(BeNil())
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
		engine := engine.NewEngine(config)
		Expect(engine).NotTo(BeNil())
		Expect(engine.Port()).To(Equal(config.Port))
		Expect(engine.BasePath()).To(Equal(config.BasePath))
		Expect(engine.PublicRouteGroup).NotTo(BeNil())
		Expect(engine.InternalRouteGroup).NotTo(BeNil())
		Expect(engine.ExternalRouteGroup).NotTo(BeNil())

		e := engine.Engine()
		Expect(e).NotTo(BeNil())

		req, err := http.NewRequest(http.MethodGet, "/api/v2/health", nil)
		Expect(err).NotTo(HaveOccurred())

		e.ServeHTTP(w, req)

		Expect(err).To(BeNil())
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("uses middleware to add a host to the context", func() {
		ip := "127.0.0.1"
		e.Engine().GET("/test", func(c *gin.Context) {
			Expect(config.GetIPAddress(c)).To(Equal(ip))
		})
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "127.0.0.1:8000"

		e.Engine().ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK))
	})

})
