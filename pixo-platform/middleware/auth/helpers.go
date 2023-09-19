package auth

import (
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

func GetParsedJWT(req *http.Request) (*RawToken, error) {
	tokenString := ExtractToken(req)

	if tokenString == "" {
		log.Debug().Msg("token not found")
		return nil, errors.New("token not found")
	}

	rawToken, err := ParseJWT(tokenString)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing JWT")
		return nil, err
	}

	return &rawToken, nil
}
