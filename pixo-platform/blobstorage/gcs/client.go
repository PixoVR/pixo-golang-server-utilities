package gcs

import (
	"os"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

var _ blobstorage.StorageClient = (*Client)(nil)

type Config struct {
	BucketName string
	Path       string
}

type Client struct {
	bucketName string
	path       string
}

func NewClient(config Config) (Client, error) {
	if config.BucketName == "" {
		config.BucketName = os.Getenv("GOOGLE_STORAGE_BUCKET")
	}

	return Client{
		bucketName: config.BucketName,
		path:       config.Path,
	}, nil
}
