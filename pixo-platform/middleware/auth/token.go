package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ValidateUserFunc func(UserID int) error

type ValidateAPIKey func(context.Context, string) (*platform.User, error)

func TokenAuthMiddleware(validateUser ValidateUserFunc, validateAPIKey ValidateAPIKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ExtractToken(c.Request) != "" {
			user, err := validateByToken(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}

			if err := validateUser(user.ID); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}
			c.Set(config.ContextRequestAuthentication.String(), user)
		} else {
			user, err := validateAPIKey(c.Request.Context(), ExtractToken(c.Request, APIKeyHeader))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}

			c.Set(config.ContextRequestAuthentication.String(), user)
		}

		c.Next()
	}
}

func validateByToken(c *gin.Context) (*platform.User, error) {
	if err := TokenValid(c.Request); err != nil {
		return nil, err
	}
	return GetParsedJWT(c.Request)
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetSecretKey()), nil
	})

	return err
}

func ExtractToken(r *http.Request, headerKeyInput ...string) string {
	headerKey := SecretKeyHeader
	if len(headerKeyInput) != 0 {
		headerKey = headerKeyInput[0]
	}

	accessToken := r.Header.Get(headerKey)
	if accessToken != "" {
		return accessToken
	}

	authToken := r.Header.Get(AuthorizationHeader)
	if len(strings.Split(authToken, " ")) == 2 {
		return strings.Split(authToken, " ")[1]
	}

	return ""
}
