package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

func (c Client) GetChecksum(ctx context.Context, object blobstorage.UploadableObject) (string, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	bucketName := c.getBucketName(object)
	fileLocation := object.GetFileLocation()

	obj := storageClient.
		Bucket(bucketName).
		Object(fileLocation)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", attrs.MD5), nil
}
