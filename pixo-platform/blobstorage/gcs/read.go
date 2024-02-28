package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"google.golang.org/api/iterator"
	"io"
)

func (c Client) FileExists(ctx context.Context, object blobstorage.UploadableObject) (bool, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return false, err
	}
	defer storageClient.Close()

	_, err = storageClient.Bucket(c.getBucketName(object)).Object(object.GetFileLocation()).Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c Client) FindFilesWithName(ctx context.Context, bucketName, prefix, filename string) ([]string, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer storageClient.Close()

	var files []string
	it := storageClient.Bucket(bucketName).Objects(ctx, &storage.Query{Prefix: prefix})
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		if blobstorage.GetFilenameFromLocation(attrs.Name) == filename {
			files = append(files, attrs.Name)
		}
	}

	return files, nil
}

func (c Client) ReadFile(ctx context.Context, object blobstorage.UploadableObject) (io.ReadCloser, error) {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer storageClient.Close()

	rc, err := storageClient.Bucket(c.getBucketName(object)).Object(object.GetFileLocation()).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return rc, nil
}
