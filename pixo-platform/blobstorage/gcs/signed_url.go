package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"time"
)

var (
	expireDuration = 120 * time.Minute
)

func (c Client) getBucketName(object client.UploadableObject) string {

	bucketName := object.GetBucketName()
	if bucketName != "" {
		return bucketName
	}

	if c.bucketName == "" {
		return c.bucketName
	}

	return os.Getenv("GOOGLE_STORAGE_BUCKET")
}

func (c Client) getFullPath(object client.UploadableObject) string {
	destination := fmt.Sprintf("%s/%s", object.GetUploadDestination(), object.GetFilename())
	return destination
}

func (c Client) GetSignedURL(ctx context.Context, object client.UploadableObject) (string, error) {
	jsonKeyPath := os.Getenv("GOOGLE_JSON_KEY")
	data, err := os.ReadFile(jsonKeyPath)
	if err != nil {
		log.Error().Err(err).Msgf("unable to read JSON key file at filepath: %s", jsonKeyPath)
		return "", err
	}

	conf, err := google.JWTConfigFromJSON(data, storage.ScopeReadOnly)
	if err != nil {
		log.Error().Err(err).Msg("unable to create JWT config from JSON key data")
		return "", err
	}

	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(conf.TokenSource(ctx)))
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage client with JSON key")
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		Expires:        time.Now().Add(expireDuration),
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
	}

	url, err := storageClient.Bucket(object.GetBucketName()).SignedURL(c.getFullPath(object), opts)
	if err != nil {
		log.Error().Err(err).Msg("unable to create signed URL")
		return "", err
	}

	log.Debug().Msgf("Created signed URL for %s", c.getFullPath(object))
	return url, nil
}
