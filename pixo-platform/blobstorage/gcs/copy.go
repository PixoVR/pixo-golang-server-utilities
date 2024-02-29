package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

func (c Client) Copy(ctx context.Context, src, dst blobstorage.UploadableObject) error {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer storageClient.Close()

	srcBucketName := c.getBucketName(src)
	srcFileLocation := src.GetFileLocation()

	dstBucketName := c.getBucketName(dst)
	dstFileLocation := dst.GetFileLocation()

	srcObject := storageClient.Bucket(srcBucketName).Object(srcFileLocation)
	dstObject := storageClient.Bucket(dstBucketName).Object(dstFileLocation)

	if _, err = dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
		return err
	}

	return nil
}
