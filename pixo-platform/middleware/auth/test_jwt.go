package auth

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

func GetTestJWT(user platform.User) string {
	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.GetEnvOrReturn("SECRET_KEY", "fake-key")))
	if err != nil {
		log.Panic().Err(err).Msg("error signing token")
	}

	return signedToken
}
