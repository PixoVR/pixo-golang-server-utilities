package aws

import (
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"io"
)

func (c Client) ReadFile(ctx context.Context, object client.UploadableObject) (io.ReadCloser, error) {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	bucketName := c.getBucketName(object)
	destination := client.GetFullPath(object)

	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucketName, Key: &destination})
	if err != nil {
		log.Error().Err(err).Msg("unable to get object")
		return nil, err
	}

	return output.Body, nil
}
