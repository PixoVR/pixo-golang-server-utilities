package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"os"

	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ client.StorageClient = (*Client)(nil)

type Config struct {
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string
}

type Client struct {
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	region          string
	endpoint        string
}

func NewClient(config Config) (Client, error) {
	if config.BucketName == "" {
		config.BucketName = os.Getenv("S3_BUCKET_NAME")
	}

	return Client{
		bucketName:      config.BucketName,
		accessKeyID:     config.AccessKeyID,
		secretAccessKey: config.SecretAccessKey,
		region:          config.Region,
		endpoint:        config.Endpoint,
	}, nil
}

func (c Client) getClient(ctx context.Context) (*s3.Client, error) {
	if c.bucketName == "" {
		return nil, errors.New("bucket is empty")
	}

	if c.region == "" {
		c.region = "us-east-1"
	}

	var customResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if c.endpoint != "" {
			return aws.Endpoint{
				URL:           c.endpoint,
				SigningRegion: c.region,
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(c.region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKeyID, c.secretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func (c Client) getBucketName(object client.UploadableObject) string {
	bucketName := object.GetBucketName()
	if bucketName != "" {
		return bucketName
	}

	if c.bucketName != "" {
		return c.bucketName
	}

	return os.Getenv("S3_BUCKET_NAME")
}
