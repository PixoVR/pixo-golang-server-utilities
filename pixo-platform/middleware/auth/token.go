package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		user, err := GetParsedJWT(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Set(UserKey, user)

		c.Next()
	}
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

func ExtractToken(r *http.Request) string {
	accessToken := r.Header.Get(SecretKeyHeader)
	if accessToken != "" {
		return accessToken
	}

	authToken := r.Header.Get(AuthorizationHeader)
	if len(strings.Split(authToken, " ")) == 2 {
		return strings.Split(authToken, " ")[1]
	}

	return ""
}
