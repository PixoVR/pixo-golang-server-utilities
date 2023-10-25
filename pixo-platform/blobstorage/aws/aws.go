package aws

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/blob"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3PresignGetObjectAPI defines the interface for the PresignGetObject function.
type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// GetPresignedURL retrieves a presigned URL for an Amazon S3 bucket object.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If successful, the presigned URL for the object and nil.
//	Otherwise, nil and an error from the call to PresignGetObject.
func GetPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input, func(options *s3.PresignOptions) { options.Expires = time.Second * 3600 })
}

var customResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	if os.Getenv("ENV") == "SAUDI" {
		return aws.Endpoint{
			URL:           "https://api-object.bluvalt.com:8082",
			SigningRegion: "us-east-1",
		}, nil
	}

	return aws.Endpoint{}, &aws.EndpointNotFoundError{}
})

func InitResumableUpload(bucketName, filepath string) (blob.ResumableUploadResponse, error) {
	return blob.ResumableUploadResponse{}, nil
}
