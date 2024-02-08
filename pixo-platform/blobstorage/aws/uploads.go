package aws

import (
	"context"
	"io"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c Client) UploadFile(ctx context.Context, object client.UploadableObject, fileReader io.Reader) (string, error) {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return "", err
	}

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.getBucketName(object)),
		Key:    aws.String(client.GetFullPath(object)),
		Body:   fileReader,
	})

	if err != nil {
		return "", err
	}

	signedURL, err := c.GetSignedURL(ctx, object)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (c Client) InitResumableUpload(ctx context.Context, object client.UploadableObject) (*client.ResumableUploadResponse, error) {

	presignClient := s3.NewPresignClient(s3.New(s3.Options{
		Region: "us-east-1",
	}))

	input := &s3.GetObjectInput{
		Bucket: aws.String(c.getBucketName(object)),
		Key:    aws.String(client.GetFullPath(object)),
	}

	presignedResponse, err := GetPresignedURL(ctx, presignClient, input)
	if err != nil {
		return nil, err
	}

	res := &client.ResumableUploadResponse{
		UploadURL:    presignedResponse.URL,
		Method:       presignedResponse.Method,
		SignedHeader: presignedResponse.SignedHeader,
	}

	return res, nil
}
