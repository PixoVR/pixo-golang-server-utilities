package aws

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
)

func (c Client) SanitizeFilename(filename string, timestamp int64) string {
	return blobstorage.SanitizeFilename(filename, timestamp)
}
