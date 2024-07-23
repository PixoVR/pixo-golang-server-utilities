package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func GetTestJWT(userID int) string {
	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(auth.GetSecretKey()))
	Expect(err).NotTo(HaveOccurred())

	return signedToken
}

var _ = Describe("Auth", func() {
	var (
		originalLifecycle string
		endpoint          = "/test"
		engine            *gin.Engine
		w                 *httptest.ResponseRecorder
		req               *http.Request
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

	Describe("Authentication middleware", func() {
		var (
			host = "127.0.0.1:8000"
			ip   = "127.0.0.1"
		)

		validateUserFunc := func(userID int) error {
			return nil
		}

		BeforeEach(func() {
			engine = gin.Default()
			engine.Use(gin.Recovery())
			engine.Use(auth.HostMiddleware())
			w = httptest.NewRecorder()

			req, _ = http.NewRequest(http.MethodGet, endpoint, nil)
			req.RemoteAddr = host
		})

		Context("JWT validation", func() {
			validateAPIKey := func(ctx context.Context, apiKey string) (*auth.User, error) {
				return &auth.User{ID: 1}, nil
			}
			BeforeEach(func() {
				jwtTokenString := GetTestJWT(1)
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenString))
			})

			It("will return the context with the ip address", func() {
				engine.GET(endpoint, func(c *gin.Context) {
					Expect(config.GetIPAddress(c)).To(Equal(ip))
				})

				engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("will return the context with the user id", func() {
				engine.Use(auth.TokenAuthMiddleware(validateUserFunc, validateAPIKey))
				engine.GET(endpoint, func(c *gin.Context) {
					Expect(auth.GetCurrentUserID(c)).To(BeNumerically(">", 0))
				})

				engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("returns unauthorized if the validator function fails", func() {
				failedValidateUserFunc := func(userID int) error {
					return fmt.Errorf("failed to validate user")
				}
				engine.Use(auth.TokenAuthMiddleware(failedValidateUserFunc, validateAPIKey))
				engine.GET(endpoint, func(c *gin.Context) {
					Expect(auth.GetCurrentUserID(c)).To(Equal(0))
				})

				engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("API key validation", func() {
			validateAPIKey := func(ctx context.Context, apiKey string) (*auth.User, error) {
				if apiKey != "test" {
					return nil, fmt.Errorf("invalid API key")
				}
				return &auth.User{ID: 1}, nil
			}

			BeforeEach(func() {
				engine.Use(auth.TokenAuthMiddleware(validateUserFunc, validateAPIKey))
			})

			It("will return the context with the user id for valid API-KEY", func() {
				req.Header.Add(auth.APIKeyHeader, "test")
				engine.GET(endpoint, func(c *gin.Context) {
					user := auth.GetUser(c)
					Expect(user).NotTo(BeNil())
					Expect(user.ID).To(Equal(1))
				})

				engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("will return unauthorized for invalid API-KEY", func() {
				req.Header.Add(auth.APIKeyHeader, "invalid-key")
				engine.GET(endpoint, func(c *gin.Context) {
					user := auth.GetUser(c)
					Expect(user).To(BeNil())
				})

				engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})
