package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTOrSecretKeyAuthMiddleware(getCurrentUser func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !IsValidSecretKey(ExtractToken(c.Request)) {

			if TokenValid(c.Request) != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}

			if err := getCurrentUser(c); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}
		}

		c.Next()
	}
}

func SecretKeyAuthMiddleware(headerKeyInput ...string) gin.HandlerFunc {
	var headerKey string
	if len(headerKeyInput) == 0 {
		headerKey = SecretKeyHeader
	} else {
		headerKey = headerKeyInput[0]
	}

	return func(c *gin.Context) {

		if !IsValidSecretKey(ExtractToken(c.Request, headerKey)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid secret key",
			})
			return
		}

		c.Next()
	}
}
