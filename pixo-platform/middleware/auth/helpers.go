package auth

import (
	"errors"
	"net/http"

	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/rs/zerolog/log"
)

func GetParsedJWT(req *http.Request) (*platform.User, error) {
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

	user := platform.User{
		ID:        rawToken.UserID,
		Email:     rawToken.Email,
		FirstName: rawToken.FirstName,
		LastName:  rawToken.LastName,
		Role:      rawToken.Role,
		OrgID:     rawToken.OrgID,
	}

	return &user, nil
}
