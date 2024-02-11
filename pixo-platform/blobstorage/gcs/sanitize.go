package gcs

import (
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"time"
)

func (c Client) SanitizeFilename(filename string) string {
	return client.SanitizeFilename(time.Now().Unix(), filename)
}
