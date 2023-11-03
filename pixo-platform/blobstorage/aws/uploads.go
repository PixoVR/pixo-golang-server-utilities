package aws

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"io"
)

func (c Client) UploadFile(ctx context.Context, object client.UploadableObject, fileReader io.Reader) (string, error) {
	log.Debug().Msgf("Uploading %s/%s", c.bucketName, object.GetUploadDestination())

	s3Client := s3.New(s3.Options{
		Region: "us-east-1",
	})

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(object.GetUploadDestination()),
	})

	signedURL, err := c.GetSignedURL(ctx, object)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (c Client) InitResumableUpload(ctx context.Context, object client.UploadableObject) (*client.ResumableUploadResponse, error) {
	log.Debug().Msgf("Initializing resumable upload for %s/%s", c.bucketName, object.GetUploadDestination())

	presignClient := s3.NewPresignClient(s3.New(s3.Options{
		Region: "us-east-1",
	}))

	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(object.GetUploadDestination()),
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
