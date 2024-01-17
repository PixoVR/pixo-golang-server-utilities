package engine

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"strconv"

	env "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	platformAuth "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CustomContext struct {
}

type Config struct {
	AddCustomContextMiddleware func(*gin.Engine)
	Port                       int
	BasePath                   string

	InternalRoutes bool
	ExternalRoutes bool
}

type CustomEngine struct {
	port     int
	basePath string

	engine *gin.Engine

	PublicRouteGroup   *gin.RouterGroup
	InternalRouteGroup *gin.RouterGroup
	ExternalRouteGroup *gin.RouterGroup
}

func NewEngine(config Config) *CustomEngine {

	e := &CustomEngine{}

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

	e.engine = gin.New()

	if config.AddCustomContextMiddleware != nil {
		config.AddCustomContextMiddleware(e.engine)
	}

	e.PublicRouteGroup = e.engine.Group(e.basePath)
	e.PublicRouteGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	e.PublicRouteGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.InternalRoutes {
		e.InternalRouteGroup = e.engine.Group(config.BasePath)
		e.InternalRouteGroup.Use(platformAuth.SecretKeyAuthMiddleware())
	}

	if config.ExternalRoutes {
		e.ExternalRouteGroup = e.engine.Group(config.BasePath)
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

func (e *CustomEngine) Engine() *gin.Engine {
	return e.engine
}

func (e *CustomEngine) Use(middleware gin.HandlerFunc) {
	e.engine.Use(middleware)
}

func (e *CustomEngine) Run() {
	if err := e.engine.Run(fmt.Sprint(e.PortString())); err != nil {
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
