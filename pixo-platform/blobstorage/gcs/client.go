package gcs

import (
	"github.com/redis/go-redis/v9"
	"os"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

var _ blobstorage.StorageClient = (*Client)(nil)

type Config struct {
	BucketName string
	Path       string
	Cache      *redis.Client
}

type Client struct {
	config Config
}

func NewClient(config Config) (Client, error) {
	if config.BucketName == "" {
		config.BucketName = os.Getenv("GOOGLE_STORAGE_BUCKET")
	}

	return Client{
		config: config,
	}, nil
}
