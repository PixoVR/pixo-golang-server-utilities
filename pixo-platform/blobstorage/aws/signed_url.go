package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"os"
)

func GetSignedURL(b, k string) (string, error) {
	awsID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if b == "" || k == "" {
		err := errors.New("bucket or key is empty")
		log.Err(err).Msg("unable to get presigned url")
		return "", err
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsID, awsSecret, "")),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to get presigned url")
		panic("configuration error, " + err.Error())
	}

	bucket := &b
	key := &k

	client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	}

	psClient := s3.NewPresignClient(client)

	resp, err := GetPresignedURL(context.TODO(), psClient, input)

	if err != nil {
		log.Error().Err(err).Msg("unable to get presigned url")
		return "", err
	}

	log.Debug().Str("url", resp.URL).Msg("presigned url")
	return resp.URL, nil
}
