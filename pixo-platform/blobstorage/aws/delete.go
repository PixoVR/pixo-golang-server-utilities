package aws

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c Client) DeleteFile(ctx context.Context, object blobstorage.UploadableObject) error {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	fileLocation := object.GetFileLocation()
	deleteObjectInput := s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    &fileLocation,
	}
	_, err = s3Client.DeleteObject(ctx, &deleteObjectInput)
	if err != nil {
		return err
	}

	return nil
}
