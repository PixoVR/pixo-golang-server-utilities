package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const Issuer = "PIXO VR"

type UserClaims struct {
	UserId             int    `json:"userId"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	Platform           string `json:"platform"`
	OrgID              int    `json:"orgId"`
	OrgType            string `json:"orgType"`
	DeviceSerialNumber string `json:"deviceSerialNumber,omitempty"`
	FingerPrintHash    string `json:"data"`
	jwt.RegisteredClaims
}

func (m UserClaims) Validate() error {
	if m.UserId == 0 {
		return errors.New("invalid user id")
	}
	if m.Role == "" {
		return errors.New("invalid user role")
	}
	if m.OrgType == "" {
		return errors.New("invalid user org type")
	}
	if m.OrgID == 0 {
		return errors.New("invalid user org id")
	}
	return nil
}

func (c UserClaims) GenerateAccessToken() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	key := os.Getenv("SECRET_KEY")
	if strings.TrimSpace(key) == "" {
		return "", errors.New("valid key is required")
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return accessToken.SignedString([]byte(key))
}

func (c UserClaims) GenerateExternalAPIAccessToken() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	key := os.Getenv("EXTERNAL_SECRET_KEY")
	if strings.TrimSpace(key) == "" {
		return "", errors.New("valid key is required")
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return accessToken.SignedString([]byte(key))
}

func ParseAccessTokenWithExpiration(accessToken string) (*UserClaims, error) {
	key := os.Getenv("SECRET_KEY")
	claims, err := ParseAccessToken(key, accessToken)
	if err != nil {
		return nil, err
	}

	validateClaims := jwt.NewValidator(jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err = validateClaims.Validate(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func ParseExternalAPIAccessToken(accessToken string) (*UserClaims, error) {
	key := os.Getenv("EXTERNAL_SECRET_KEY")

	return ParseAccessToken(key, accessToken)
}

func ParseAccessToken(key, accessToken string) (*UserClaims, error) {
	if strings.TrimSpace(key) == "" {
		return nil, errors.New("valid key is required")
	}

	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok || !parsedAccessToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
