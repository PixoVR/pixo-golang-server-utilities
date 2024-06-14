package engine

import (
	"fmt"
	env "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uber/jaeger-client-go"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
	"os"
	"strconv"
	"strings"

	platformAuth "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type CustomContext struct {
}

type Config struct {
	AddCustomContextMiddleware func(*gin.Engine)
	Port                       int
	BasePath                   string

	Lifecycle         string
	Tracing           bool
	CollectorEndpoint string
	InternalRoutes    bool
	ExternalRoutes    bool
}

type CustomEngine struct {
	*gin.Engine
	port     int
	basePath string

	PublicRouteGroup   *gin.RouterGroup
	InternalRouteGroup *gin.RouterGroup
	ExternalRouteGroup *gin.RouterGroup
}

func NewEngine(config Config) *CustomEngine {

	e := &CustomEngine{
		Engine: gin.New(),
	}

	if config.Port != 0 {
		e.port = config.Port
	} else {
		e.port = e.findPort()
	}

	if config.BasePath != "" {
		e.basePath = config.BasePath
	} else {
		e.basePath = DefaultBasePath
	}

	lifecycle := env.GetLifecycle()
	if lifecycle == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	e.Use(gin.Recovery())
	e.Use(platformAuth.HostMiddleware())

	if config.Tracing {
		if os.Getenv("DD_ENV") != "" {
			e.Use(gintrace.Middleware(strings.ReplaceAll(e.basePath, "/", "")))
			tracer.Start(tracer.WithRuntimeMetrics())

		} else {
			if config.CollectorEndpoint == "" {
				if strings.ToLower(config.Lifecycle) == "local" {
					config.CollectorEndpoint = "http://localhost:16686"
				}
				config.CollectorEndpoint = "http://jaeger.linkerd-jaeger.svc:16686"
			}

			cfg := &jaegerConfig.Configuration{
				ServiceName: fmt.Sprintf("%s-service", e.basePath),
				Sampler: &jaegerConfig.SamplerConfig{
					Type:  "const",
					Param: 1,
				},
				Reporter: &jaegerConfig.ReporterConfig{
					LogSpans:          true,
					CollectorEndpoint: config.CollectorEndpoint,
				},
				// Token configuration
				//Tags: []opentracing.Tag{ // Set the tag, where information such as token can be stored.
				//	{Key: "token", Value: token},
				//},
			}

			t, _, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to initialize tracer")
			}
			//defer closer.Close() // nolint: errcheck

			e.Use(ginhttp.Middleware(t))
		}
	}

	if config.AddCustomContextMiddleware != nil {
		config.AddCustomContextMiddleware(e.Engine)
	}

	e.PublicRouteGroup = e.Group(e.basePath)
	e.PublicRouteGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	e.PublicRouteGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.InternalRoutes {
		e.InternalRouteGroup = e.Group(config.BasePath)
		e.InternalRouteGroup.Use(platformAuth.SecretKeyAuthMiddleware())
	}

	if config.ExternalRoutes {
		e.ExternalRouteGroup = e.Group(config.BasePath)
		e.ExternalRouteGroup.Use(platformAuth.JWTOrSecretKeyAuthMiddleware(func(c *gin.Context) error {
			return nil
		}))
	}

	return e
}

func (e *CustomEngine) Port() int {
	return e.port
}

func (e *CustomEngine) PortString() string {
	return fmt.Sprintf(":%d", e.port)
}

func (e *CustomEngine) BasePath() string {
	return e.basePath
}

func (e *CustomEngine) Start() {
	if err := e.Run(fmt.Sprint(e.PortString())); err != nil {
		log.Fatal().Err(err).Msg("Unable to start server")
	}
}

func (e *CustomEngine) findPort(portInput ...int) int {

	if len(portInput) > 0 {
		return portInput[0]
	}

	portString, ok := os.LookupEnv("PORT")
	if !ok {
		return DefaultPort
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return DefaultPort
	}

	return port
}
