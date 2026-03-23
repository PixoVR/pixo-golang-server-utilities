package auth_test

import (
	"os"
	"time"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func validUserClaims() auth.UserClaims {
	return auth.UserClaims{
		UserId:  1,
		Email:   "test@example.com",
		Role:    "admin",
		OrgID:   1,
		OrgType: "platform",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    auth.Issuer,
		},
	}
}

var _ = Describe("UserClaims", func() {

	Describe("Validate", func() {
		It("returns nil for valid claims", func() {
			claims := validUserClaims()
			Expect(claims.Validate()).NotTo(HaveOccurred())
		})

		It("returns error if user id is 0", func() {
			claims := validUserClaims()
			claims.UserId = 0
			Expect(claims.Validate()).To(MatchError("invalid user id"))
		})

		It("returns error if role is empty", func() {
			claims := validUserClaims()
			claims.Role = ""
			Expect(claims.Validate()).To(MatchError("invalid user role"))
		})

		It("returns error if org type is empty", func() {
			claims := validUserClaims()
			claims.OrgType = ""
			Expect(claims.Validate()).To(MatchError("invalid user org type"))
		})

		It("returns error if org id is 0", func() {
			claims := validUserClaims()
			claims.OrgID = 0
			Expect(claims.Validate()).To(MatchError("invalid user org id"))
		})
	})

	Describe("GenerateAccessToken", func() {
		BeforeEach(func() {
			Expect(os.Setenv("SECRET_KEY", "test-secret")).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.Unsetenv("SECRET_KEY")).NotTo(HaveOccurred())
		})

		It("generates a token for valid claims", func() {
			claims := validUserClaims()
			token, err := claims.GenerateAccessToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())
		})

		It("returns error for invalid claims", func() {
			claims := validUserClaims()
			claims.UserId = 0
			_, err := claims.GenerateAccessToken()
			Expect(err).To(HaveOccurred())
		})

		It("returns error if SECRET_KEY is empty", func() {
			Expect(os.Setenv("SECRET_KEY", "")).NotTo(HaveOccurred())
			claims := validUserClaims()
			_, err := claims.GenerateAccessToken()
			Expect(err).To(MatchError("valid key is required"))
		})

		It("returns error if SECRET_KEY is unset", func() {
			Expect(os.Unsetenv("SECRET_KEY")).NotTo(HaveOccurred())
			claims := validUserClaims()
			_, err := claims.GenerateAccessToken()
			Expect(err).To(MatchError("valid key is required"))
		})
	})

	Describe("GenerateExternalAPIAccessToken", func() {
		BeforeEach(func() {
			Expect(os.Setenv("EXTERNAL_SECRET_KEY", "external-test-secret")).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.Unsetenv("EXTERNAL_SECRET_KEY")).NotTo(HaveOccurred())
		})

		It("generates a token for valid claims", func() {
			claims := validUserClaims()
			token, err := claims.GenerateExternalAPIAccessToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())
		})

		It("returns error for invalid claims", func() {
			claims := validUserClaims()
			claims.Role = ""
			_, err := claims.GenerateExternalAPIAccessToken()
			Expect(err).To(HaveOccurred())
		})

		It("returns error if EXTERNAL_SECRET_KEY is empty", func() {
			Expect(os.Setenv("EXTERNAL_SECRET_KEY", "")).NotTo(HaveOccurred())
			claims := validUserClaims()
			_, err := claims.GenerateExternalAPIAccessToken()
			Expect(err).To(MatchError("valid key is required"))
		})

		It("returns error if EXTERNAL_SECRET_KEY is unset", func() {
			Expect(os.Unsetenv("EXTERNAL_SECRET_KEY")).NotTo(HaveOccurred())
			claims := validUserClaims()
			_, err := claims.GenerateExternalAPIAccessToken()
			Expect(err).To(MatchError("valid key is required"))
		})
	})

	Describe("ParseAccessTokenWithExpiration", func() {
		BeforeEach(func() {
			Expect(os.Setenv("SECRET_KEY", "test-secret")).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.Unsetenv("SECRET_KEY")).NotTo(HaveOccurred())
		})

		It("parses a valid token and returns claims", func() {
			claims := validUserClaims()
			token, err := claims.GenerateAccessToken()
			Expect(err).NotTo(HaveOccurred())

			parsed, err := auth.ParseAccessTokenWithExpiration(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.UserId).To(Equal(claims.UserId))
			Expect(parsed.Email).To(Equal(claims.Email))
			Expect(parsed.Role).To(Equal(claims.Role))
			Expect(parsed.OrgID).To(Equal(claims.OrgID))
			Expect(parsed.OrgType).To(Equal(claims.OrgType))
		})

		It("returns error for empty token", func() {
			_, err := auth.ParseAccessTokenWithExpiration("")
			Expect(err).To(HaveOccurred())
		})

		It("returns error for invalid token", func() {
			_, err := auth.ParseAccessTokenWithExpiration("invalid-token")
			Expect(err).To(HaveOccurred())
		})

		It("returns error if SECRET_KEY is empty", func() {
			Expect(os.Setenv("SECRET_KEY", "")).NotTo(HaveOccurred())
			_, err := auth.ParseAccessTokenWithExpiration("some-token")
			Expect(err).To(MatchError("valid key is required"))
		})

		It("returns error for expired token", func() {
			claims := validUserClaims()
			claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
			token, err := claims.GenerateAccessToken()
			Expect(err).NotTo(HaveOccurred())

			_, err = auth.ParseAccessTokenWithExpiration(token)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ParseExternalAPIAccessToken", func() {
		BeforeEach(func() {
			Expect(os.Setenv("EXTERNAL_SECRET_KEY", "external-test-secret")).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.Unsetenv("EXTERNAL_SECRET_KEY")).NotTo(HaveOccurred())
		})

		It("parses a valid external API token", func() {
			claims := validUserClaims()
			token, err := claims.GenerateExternalAPIAccessToken()
			Expect(err).NotTo(HaveOccurred())

			parsed, err := auth.ParseExternalAPIAccessToken(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.UserId).To(Equal(claims.UserId))
			Expect(parsed.Role).To(Equal(claims.Role))
			Expect(parsed.OrgID).To(Equal(claims.OrgID))
			Expect(parsed.OrgType).To(Equal(claims.OrgType))
		})

		It("returns error for empty token", func() {
			_, err := auth.ParseExternalAPIAccessToken("")
			Expect(err).To(HaveOccurred())
		})

		It("returns error if EXTERNAL_SECRET_KEY is empty", func() {
			Expect(os.Setenv("EXTERNAL_SECRET_KEY", "")).NotTo(HaveOccurred())
			_, err := auth.ParseExternalAPIAccessToken("some-token")
			Expect(err).To(MatchError("valid key is required"))
		})
	})

	Describe("ParseAccessToken", func() {
		It("returns error if key is empty", func() {
			_, err := auth.ParseAccessToken("", "some-token")
			Expect(err).To(MatchError("valid key is required"))
		})

		It("returns error if key is whitespace only", func() {
			_, err := auth.ParseAccessToken("   ", "some-token")
			Expect(err).To(MatchError("valid key is required"))
		})

		It("parses a valid token with the correct key", func() {
			claims := validUserClaims()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signed, err := token.SignedString([]byte("my-key"))
			Expect(err).NotTo(HaveOccurred())

			parsed, err := auth.ParseAccessToken("my-key", signed)
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.UserId).To(Equal(1))
		})

		It("returns error for wrong key", func() {
			claims := validUserClaims()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signed, err := token.SignedString([]byte("correct-key"))
			Expect(err).NotTo(HaveOccurred())

			_, err = auth.ParseAccessToken("wrong-key", signed)
			Expect(err).To(HaveOccurred())
		})
	})
})
