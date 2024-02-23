package auth

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/gin-gonic/gin"
)

func HostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.ContextRequestHost.String(), c.ClientIP())
		c.Next()
	}
}
