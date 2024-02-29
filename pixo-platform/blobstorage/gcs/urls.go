package gcs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/storage"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var expireDuration = 120 * time.Minute

func (c Client) getBucketName(object blobstorage.UploadableObject) string {
	bucketName := object.GetBucketName()
	if bucketName != "" {
		return bucketName
	}

	if c.bucketName == "" {
		c.bucketName = os.Getenv("GOOGLE_STORAGE_BUCKET")
	}

	return c.bucketName
}

func (c Client) GetPublicURL(object blobstorage.UploadableObject) string {
	bucketName := c.getBucketName(object)

	fileLocation := object.GetFileLocation()
	if fileLocation == "" {
		return ""
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileLocation)
}

func (c Client) GetSignedURL(ctx context.Context, object blobstorage.UploadableObject, options ...blobstorage.Option) (string, error) {
	data, err := os.ReadFile(os.Getenv("GOOGLE_JSON_KEY"))
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
	if len(options) > 0 {
		option := options[0]

		if option.ContentDisposition != "" {
			opts.QueryParameters = url.Values{
				"response-content-disposition": {option.ContentDisposition},
			}
		}

		if option.Expires != nil {
			opts.Expires = *option.Expires
		}

		if option.Method != "" {
			opts.Method = option.Method
		}
	}

	return storageClient.Bucket(c.getBucketName(object)).SignedURL(object.GetFileLocation(), opts)
}
