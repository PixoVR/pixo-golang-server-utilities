package aws

import (
	"errors"
)

type Config struct {
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type Client struct {
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	region          string
}

func NewClient(config Config) (Client, error) {

	if config.BucketName == "" {
		return Client{}, errors.New("no bucket name given")
	}

	return Client{
		bucketName:      config.BucketName,
		accessKeyID:     config.AccessKeyID,
		secretAccessKey: config.SecretAccessKey,
		region:          config.Region,
	}, nil
}
