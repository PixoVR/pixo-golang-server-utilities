package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c Client) GetPublicURL(object blobstorage.UploadableObject) string {
	bucketName := c.getBucketName(object)
	fileLocation := object.GetFileLocation()
	if fileLocation == "" {
		return ""
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileLocation)
}

func (c Client) GetSignedURL(ctx context.Context, object blobstorage.UploadableObject, options ...blobstorage.Option) (string, error) {
	s3Client, err := c.getClient(ctx)
	if err != nil {
		return "", err
	}

	bucketName := c.getBucketName(object)
	destination := object.GetFileLocation()

	input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &destination,
	}

	presignClient := s3.NewPresignClient(s3Client)

	resp, err := presignClient.PresignGetObject(ctx, input, func(options *s3.PresignOptions) { options.Expires = time.Second * 3600 })
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
