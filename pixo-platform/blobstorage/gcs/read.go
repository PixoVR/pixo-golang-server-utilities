package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"io"
)

func (c Client) FileExists(ctx context.Context, object client.UploadableObject) (bool, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return false, err
	}

	_, err = storageClient.Bucket(c.getBucketName(object)).Object(object.GetFileLocation()).Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c Client) ReadFile(ctx context.Context, object client.UploadableObject) (io.ReadCloser, error) {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	rc, err := storageClient.Bucket(c.getBucketName(object)).Object(object.GetFileLocation()).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return rc, nil
}
