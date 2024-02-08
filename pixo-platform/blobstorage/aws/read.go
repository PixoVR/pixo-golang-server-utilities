package aws

import (
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

func (c Client) FileExists(ctx context.Context, object client.UploadableObject) (bool, error) {
	s3Client, err := c.getClient(ctx)
	if err != nil {
		return false, err
	}

	bucketName := c.getBucketName(object)
	destination := object.GetFileLocation()

	if _, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &bucketName, Key: &destination}); err != nil {
		return false, nil
	}

	return true, nil
}

func (c Client) ReadFile(ctx context.Context, object client.UploadableObject) (io.ReadCloser, error) {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	bucketName := c.getBucketName(object)
	destination := object.GetFileLocation()

	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucketName, Key: &destination})
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}
