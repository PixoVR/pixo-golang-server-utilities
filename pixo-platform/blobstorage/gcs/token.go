package gcs

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)

func getAccessToken() (string, error) {
	var token *oauth2.Token
	ctx := context.Background()

	scopes := []string{
		CloudPlatformScope,
	}

	credentials, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		log.Error().Err(err).Msg("unable to find default credentials")
		return "", err
	}

	token, err = credentials.TokenSource.Token()
	if err != nil {
		log.Error().Err(err).Msg("unable to get token credentials")
		return "", err
	}

	return token.AccessToken, nil
}
