package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

func (c Client) DeleteFile(ctx context.Context, object client.UploadableObject) error {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	if err = storageClient.Bucket(c.getBucketName(object)).Object(object.GetFileLocation()).Delete(ctx); err != nil {
		return err
	}

	return nil
}
