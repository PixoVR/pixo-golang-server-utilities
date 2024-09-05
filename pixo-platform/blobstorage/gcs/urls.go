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

var DefaultExpireDuration = 120 * time.Minute

func (c Client) getBucketName(object blobstorage.UploadableObject) string {
	bucketName := object.GetBucketName()
	if bucketName != "" {
		return bucketName
	}

	if c.config.BucketName == "" {
		c.config.BucketName = os.Getenv("GOOGLE_STORAGE_BUCKET")
	}

	return c.config.BucketName
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
	if signedURL := c.cacheGet(ctx, object); signedURL != "" {
		return signedURL, nil
	}

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
		Expires:        time.Now().Add(DefaultExpireDuration),
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
	}
	if len(options) > 0 {
		opt := options[0]

		if opt.ContentDisposition != "" {
			opts.QueryParameters = url.Values{
				"response-content-disposition": {opt.ContentDisposition},
			}
		}

		if opt.Lifetime.String() != "0s" {
			opts.Expires = time.Now().Add(opt.Lifetime)
		}

		if opt.Method != "" {
			opts.Method = opt.Method
		}
	}

	signedURL, err := storageClient.Bucket(c.getBucketName(object)).SignedURL(object.GetFileLocation(), opts)
	if err != nil {
		return "", err
	}

	c.cacheSet(ctx, signedURL, object, options...)
	return signedURL, nil
}

func (c Client) cacheSet(ctx context.Context, signedURL string, object blobstorage.UploadableObject, options ...blobstorage.Option) {
	if c.config.Cache != nil {
		expiration := DefaultExpireDuration
		if len(options) > 0 && options[0].Lifetime != 0 {
			expiration = options[0].Lifetime - 1*time.Second
		}
		c.config.Cache.Set(ctx, c.CacheKey(object), signedURL, expiration)
	}
}

func (c Client) cacheGet(ctx context.Context, object blobstorage.UploadableObject) string {
	if c.config.Cache != nil {
		return c.config.Cache.Get(ctx, c.CacheKey(object)).Val()
	}

	return ""
}

func (c Client) CacheKey(object blobstorage.UploadableObject) string {
	return fmt.Sprintf("signed-url:%s/%s", c.getBucketName(object), object.GetFileLocation())
}
