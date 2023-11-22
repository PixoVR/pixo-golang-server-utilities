package aws

import (
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (c Client) GetSignedURL(ctx context.Context, object client.UploadableObject) (string, error) {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return "", err
	}

	bucketName := c.getBucketName(object)
	destination := client.GetFullPath(object)

	input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &destination,
	}

	psClient := s3.NewPresignClient(s3Client)

	resp, err := GetPresignedURL(ctx, psClient, input)
	if err != nil {
		log.Error().Err(err).Msg("unable to get presigned url")
		return "", err
	}

	log.Debug().Str("url", resp.URL).Msg("presigned url")
	return resp.URL, nil
}
