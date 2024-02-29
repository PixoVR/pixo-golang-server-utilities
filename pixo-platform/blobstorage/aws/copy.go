package aws

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c Client) Copy(ctx context.Context, src, dst blobstorage.UploadableObject) error {
	s3Client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cpSrc := fmt.Sprintf("%s/%s", c.getBucketName(src), src.GetFileLocation())

	dstBucketName := c.getBucketName(dst)
	dstFileLocation := dst.GetFileLocation()

	_, err = s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &dstBucketName,
		CopySource: &cpSrc,
		Key:        &dstFileLocation,
	})
	if err != nil {
		return err
	}

	return nil
}
