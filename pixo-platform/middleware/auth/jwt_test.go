package auth_test

import (
	"encoding/base64"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseJWT", func() {
	It("extracts legacy compact claims without validating the signature or expiry", func() {
		claims := jwt.MapClaims{
			"authorized": true, "userId": 42, "orgId": 7, "orgType": "platform",
			"given_name": "Ada", "family_name": "Lovelace", "email": "ada@example.com",
			"role": "admin", "aud": "pixo", "exp": 1, "iat": 0.5, "iss": "issuer",
			"sub": "subject", "jti": "id", "email_verified": true, "hd": "example.com", "data": "value",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte("unverified-test-key"))
		Expect(err).NotTo(HaveOccurred())

		raw, err := auth.ParseJWT(" \n" + signed + "\t")
		Expect(err).NotTo(HaveOccurred())
		Expect(raw).To(Equal(auth.RawToken{
			Authorized: true, UserID: 42, OrgID: 7, OrgType: "platform",
			FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com",
			Role: "admin", Audience: "pixo", Expiration: 1, IAT: 0.5, Issuer: "issuer",
			Sub: "subject", JTI: "id", EmailVerified: true, Hd: "example.com", Data: "value",
		}))
	})

	It("extracts claims from flattened and general JSON JWS", func() {
		key := []byte("0123456789abcdef0123456789abcdef")
		signingKey := jose.SigningKey{Algorithm: jose.HS256, Key: key}
		flattened, err := jose.NewSigner(signingKey, nil)
		Expect(err).NotTo(HaveOccurred())
		general, err := jose.NewMultiSigner([]jose.SigningKey{signingKey, signingKey}, nil)
		Expect(err).NotTo(HaveOccurred())

		for _, signer := range []jose.Signer{flattened, general} {
			signed, err := signer.Sign([]byte(`{"userId":42,"aud":["not a legacy string"],"orgId":"not a legacy number"}`))
			Expect(err).NotTo(HaveOccurred())
			raw, err := auth.ParseJWT(signed.FullSerialize())
			Expect(err).NotTo(HaveOccurred())
			Expect(raw).To(Equal(auth.RawToken{UserID: 42}))
		}
	})

	DescribeTable("returns empty claims for unusable payloads in valid envelopes", func(payload string) {
		signer, err := jose.NewSigner(jose.SigningKey{
			Algorithm: jose.HS256, Key: []byte("0123456789abcdef0123456789abcdef"),
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		signed, err := signer.Sign([]byte(payload))
		Expect(err).NotTo(HaveOccurred())
		compact, err := signed.CompactSerialize()
		Expect(err).NotTo(HaveOccurred())
		for _, serialized := range []string{compact, signed.FullSerialize()} {
			raw, err := auth.ParseJWT(serialized)
			Expect(err).NotTo(HaveOccurred())
			Expect(raw).To(BeZero())
		}
	},
		Entry("an empty payload", ""),
		Entry("malformed JSON", "{"),
		Entry("a non-object payload", "[]"),
	)

	It("accepts an unregistered compact algorithm when extracting claims", func() {
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"custom"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(`{"userId":42}`))
		raw, err := auth.ParseJWT(header + "." + payload + ".")
		Expect(err).NotTo(HaveOccurred())
		Expect(raw).To(Equal(auth.RawToken{UserID: 42}))
	})

	DescribeTable("rejects malformed envelopes", func(serialized string) {
		raw, err := auth.ParseJWT(serialized)
		Expect(err).To(MatchError("no token found"))
		Expect(raw).To(BeZero())
	},
		Entry("missing segments", "not-a-token"),
		Entry("an invalid compact signature encoding", "eyJhbGciOiJIUzI1NiJ9.e30.!"),
		Entry("JSON JWS without a payload", `{"signature":"AA"}`),
	)
})
