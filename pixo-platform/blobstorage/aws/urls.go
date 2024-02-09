package aws

import (
	"context"
	"fmt"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c Client) GetPublicURL(object client.UploadableObject) (string, error) {
	bucketName := c.getBucketName(object)
	destination := object.GetFileLocation()
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, destination), nil
}

func (c Client) GetSignedURL(ctx context.Context, object client.UploadableObject) (string, error) {

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

	psClient := s3.NewPresignClient(s3Client)

	resp, err := GetPresignedURL(ctx, psClient, input)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
