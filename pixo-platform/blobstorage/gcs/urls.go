package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
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
		c.bucketName = os.Getenv("GOOGLE_STORAGE_BUCKET")
	}

	return c.bucketName
}

func (c Client) GetPublicURL(object client.UploadableObject) (string, error) {
	bucketName := c.getBucketName(object)
	return "https://storage.googleapis.com/" + bucketName + "/" + object.GetFileLocation(), nil
}

func (c Client) GetSignedURL(ctx context.Context, object client.UploadableObject) (string, error) {
	jsonKeyPath := os.Getenv("GOOGLE_JSON_KEY")
	data, err := os.ReadFile(jsonKeyPath)
	if err != nil {
		return "", err
	}

	conf, err := google.JWTConfigFromJSON(data, storage.ScopeReadOnly)
	if err != nil {
		return "", err
	}

	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(conf.TokenSource(ctx)))
	if err != nil {
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		Expires:        time.Now().Add(expireDuration),
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
	}

	url, err := storageClient.Bucket(c.getBucketName(object)).SignedURL(object.GetFileLocation(), opts)
	if err != nil {
		return "", err
	}

	return url, nil
}
