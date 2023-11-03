package aws

import (
	"context"
	"errors"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (c Client) GetSignedURL(ctx context.Context, object client.UploadableObject) (string, error) {

	if c.bucketName == "" || object.GetUploadDestination() == "" {
		err := errors.New("bucket or destination is empty")
		log.Err(err).Msg("unable to get presigned url")
		return "", err
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
		log.Fatal().Err(err).Msg("unable to get presigned url")
	}

	awsClient := s3.NewFromConfig(cfg)
	uploadDestination := object.GetUploadDestination()

	input := &s3.GetObjectInput{
		Bucket: &c.bucketName,
		Key:    &uploadDestination,
	}

	psClient := s3.NewPresignClient(awsClient)

	resp, err := GetPresignedURL(ctx, psClient, input)
	if err != nil {
		log.Error().Err(err).Msg("unable to get presigned url")
		return "", err
	}

	log.Debug().Str("url", resp.URL).Msg("presigned url")
	return resp.URL, nil
}
