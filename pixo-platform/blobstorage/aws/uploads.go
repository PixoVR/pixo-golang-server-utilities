package aws

import (
	"context"
	storage "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
)

func (c Client) UploadFile(ctx context.Context, object storage.UploadableObject, fileReader io.Reader) (string, error) {

	s3Client, err := c.getClient(ctx)
	if err != nil {
		return "", err
	}

	sanitizedFileLocation := c.SanitizeFilename(object.GetFileLocation())

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.getBucketName(object)),
		Key:    aws.String(sanitizedFileLocation),
		Body:   fileReader,
	})

	if err != nil {
		return "", err
	}

	return sanitizedFileLocation, nil
}

func (c Client) InitResumableUpload(ctx context.Context, object storage.UploadableObject) (*storage.ResumableUploadResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	sanitizedFileLocation := c.SanitizeFilename(object.GetFileLocation())

	res, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.getBucketName(object)),
		Key:    aws.String(sanitizedFileLocation),
	})
	if err != nil {
		panic("failed to presign request, " + err.Error())
	}

	var uploadRes storage.ResumableUploadResponse

	uploadRes.SignedHeader = res.SignedHeader
	uploadRes.UploadURL = res.URL
	uploadRes.Method = res.Method

	return &uploadRes, nil
}
