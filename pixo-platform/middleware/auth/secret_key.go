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

func SecretKeyAuthMiddleware(getCurrentUser func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !IsValidSecretKey(ExtractToken(c.Request)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}
