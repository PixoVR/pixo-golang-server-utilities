package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"unicode"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
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
	claims, err := parseLegacyJWTClaims(tokenString)
	if err != nil {
		return RawToken{}, errors.New("no token found")
	}

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

func parseLegacyJWTClaims(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, tokenString)
	parser := jwt.NewParser()
	if strings.HasPrefix(tokenString, "{") {
		return parseLegacyJSONJWS(tokenString)
	}

	return parseLegacyCompactJWT(tokenString, parser)
}

func parseLegacyCompactJWT(tokenString string, parser *jwt.Parser) (jwt.MapClaims, error) {
	segments := strings.SplitN(tokenString, ".", 4)
	if len(segments) != 3 {
		return nil, jwt.ErrTokenMalformed
	}
	for i := range segments {
		segments[i] = strings.TrimRight(segments[i], "=")
	}
	emptyObject := base64.RawURLEncoding.EncodeToString([]byte("{}"))
	if segments[0] == "" {
		segments[0] = emptyObject
	}
	tokenString = strings.Join(segments, ".")
	claims := jwt.MapClaims{}
	token, parts, err := parser.ParseUnverified(tokenString, claims)
	if errors.Is(err, jwt.ErrTokenMalformed) && len(parts) == 3 {
		payload, decodeErr := parser.DecodeSegment(parts[1])
		if decodeErr == nil && json.Unmarshal(payload, &claims) != nil {
			claims = jwt.MapClaims{}
			parts[1] = emptyObject
			token, parts, err = parser.ParseUnverified(strings.Join(parts, "."), claims)
		}
	}
	if errors.Is(err, jwt.ErrTokenUnverifiable) {
		if alg := token.Header["alg"]; alg != nil {
			if _, ok := alg.(string); !ok {
				return nil, err
			}
		}
		_, err = parser.DecodeSegment(parts[2])
	}
	if err != nil {
		return nil, err
	}
	for _, field := range []string{"jwk", "kid", "nonce", "x5c"} {
		if token.Header[field] != nil {
			alg, _ := token.Header["alg"].(string)
			if _, err := jose.ParseSignedCompact(tokenString, []jose.SignatureAlgorithm{jose.SignatureAlgorithm(alg)}); err != nil {
				return nil, err
			}
			break
		}
	}
	return claims, nil
}

func parseLegacyJSONJWS(tokenString string) (jwt.MapClaims, error) {
	algorithms := []jose.SignatureAlgorithm{
		"", "none", jose.EdDSA,
		jose.HS256, jose.HS384, jose.HS512,
		jose.RS256, jose.RS384, jose.RS512,
		jose.ES256, jose.ES384, jose.ES512,
		jose.PS256, jose.PS384, jose.PS512,
	}
	token, err := jose.ParseSignedJSON(tokenString, algorithms)
	if err != nil {
		return nil, err
	}
	var claims jwt.MapClaims
	_ = json.Unmarshal(token.UnsafePayloadWithoutVerification(), &claims)
	return claims, nil
}

func extractClaim(claims map[string]interface{}, key string) interface{} {
	if value, ok := claims[key]; ok {
		return value
	}

	return nil
}
