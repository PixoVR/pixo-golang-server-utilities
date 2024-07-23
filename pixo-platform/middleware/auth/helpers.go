package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type User struct {
	ID        int    `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func GetParsedJWT(req *http.Request) (*User, error) {
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

	user := User{
		ID:        rawToken.UserID,
		Email:     rawToken.Email,
		FirstName: rawToken.FirstName,
		LastName:  rawToken.LastName,
	}

	return &user, nil
}
