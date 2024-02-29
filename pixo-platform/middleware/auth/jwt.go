package auth

import (
	"errors"

	jwt2 "github.com/go-jose/go-jose/v3/jwt"
)

type RawToken struct {
	Authorized    bool    `json:"authorized"`
	UserID        int     `json:"userId"`
	OrgType       string  `json:"orgType"`
	OrgID         int     `json:"orgId"`
	FirstName     string  `json:"given_name"`
	LastName      string  `json:"family_name"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Audience      string  `json:"aud"`
	Expiration    int64   `json:"exp"`
	IAT           float64 `json:"iat"`
	Issuer        string  `json:"iss"`
	Sub           string  `json:"sub"`
	JTI           string  `json:"jti"`
	EmailVerified bool    `json:"email_verified"`
	Hd            string  `json:"hd"`
	Data          string  `json:"data"`
}

func ParseJWT(tokenString string) (RawToken, error) {
	var claims map[string]interface{}

	var token *jwt2.JSONWebToken
	token, err := jwt2.ParseSigned(tokenString)
	if err != nil {
		return RawToken{}, errors.New("no token found")
	}

	_ = token.UnsafeClaimsWithoutVerification(&claims)

	var rawToken RawToken

	if userID, ok := extractClaim(claims, "userId").(float64); ok {
		rawToken.UserID = int(userID)
	}

	if authorized, ok := extractClaim(claims, "authorized").(bool); ok {
		rawToken.Authorized = authorized
	}

	if audience, ok := extractClaim(claims, "aud").(string); ok {
		rawToken.Audience = audience
	}

	if expiration, ok := extractClaim(claims, "exp").(float64); ok {
		rawToken.Expiration = int64(expiration)
	}

	if iat, ok := extractClaim(claims, "iat").(float64); ok {
		rawToken.IAT = iat
	}

	if issuer, ok := extractClaim(claims, "iss").(string); ok {
		rawToken.Issuer = issuer
	}

	if sub, ok := extractClaim(claims, "sub").(string); ok {
		rawToken.Sub = sub
	}

	if jti, ok := extractClaim(claims, "jti").(string); ok {
		rawToken.JTI = jti
	}

	if firstName, ok := extractClaim(claims, "given_name").(string); ok {
		rawToken.FirstName = firstName
	}

	if lastName, ok := extractClaim(claims, "family_name").(string); ok {
		rawToken.LastName = lastName
	}

	if email, ok := extractClaim(claims, "email").(string); ok {
		rawToken.Email = email
	}

	if orgID, ok := extractClaim(claims, "orgId").(float64); ok {
		rawToken.OrgID = int(orgID)
	}

	if orgType, ok := extractClaim(claims, "orgType").(string); ok {
		rawToken.OrgType = orgType
	}

	if role, ok := extractClaim(claims, "role").(string); ok {
		rawToken.Role = role
	}

	if emailVerified, ok := extractClaim(claims, "email_verified").(bool); ok {
		rawToken.EmailVerified = emailVerified
	}

	if hd, ok := extractClaim(claims, "hd").(string); ok {
		rawToken.Hd = hd
	}

	if data, ok := extractClaim(claims, "data").(string); ok {
		rawToken.Data = data
	}

	return rawToken, nil
}

func extractClaim(claims map[string]interface{}, key string) interface{} {
	if value, ok := claims[key]; ok {
		return value
	}

	return nil
}
