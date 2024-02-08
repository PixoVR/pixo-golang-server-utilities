package gcs

import (
	"errors"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

var _ client.StorageClient = (*Client)(nil)

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
		return Client{}, errors.New("no bucket name given")
	}

	return Client{
		bucketName: config.BucketName,
		path:       config.Path,
	}, nil
}
