package aws

import (
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c Client) DeleteFile(ctx context.Context, object client.UploadableObject) error {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	destination := client.GetFullPath(object)
	deleteObjectInput := s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    &destination,
	}
	_, err = s3Client.DeleteObject(ctx, &deleteObjectInput)
	if err != nil {
		return err
	}

	return nil
}
