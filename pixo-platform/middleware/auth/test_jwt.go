package auth

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func GetTestJWT(userID int) string {
	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.GetEnvOrReturn("SECRET_KEY", "fake-key")))
	if err != nil {
		log.Panic().Err(err).Msg("error signing token")
	}

	return signedToken
}
