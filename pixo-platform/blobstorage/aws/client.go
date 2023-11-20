package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
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

func (c Client) getClient(ctx context.Context) (*s3.Client, error) {

	if c.bucketName == "" {
		err := errors.New("bucket is empty")
		log.Err(err).Msg("unable to get presigned url")
		return nil, err
	}

	if c.region == "" {
		c.region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(c.region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKeyID, c.secretAccessKey, "")),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to get presigned url")
		return nil, err
	}

	awsClient := s3.NewFromConfig(cfg)

	return awsClient, nil
}
