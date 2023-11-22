package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/rs/zerolog/log"
)

func (c Client) DeleteFile(ctx context.Context, object client.UploadableObject) error {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage client")
		return err
	}

	if err = storageClient.Bucket(c.getBucketName(object)).Object(client.GetFullPath(object)).Delete(ctx); err != nil {
		log.Error().Err(err).Msg("unable to delete storage object")
		return err
	}

	return nil
}
