package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/rs/zerolog/log"
	"io"
)

func (c Client) ReadFile(ctx context.Context, object client.UploadableObject) (io.ReadCloser, error) {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage client")
		return nil, err
	}

	rc, err := storageClient.Bucket(c.getBucketName(object)).Object(c.getFullPath(object)).NewReader(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage reader")
		return nil, err
	}

	return rc, nil
}
