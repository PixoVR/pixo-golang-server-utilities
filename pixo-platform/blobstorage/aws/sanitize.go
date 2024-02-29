package aws

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"time"
)

func (c Client) SanitizeFilename(filename string) string {
	return blobstorage.SanitizeFilename(time.Now().Unix(), filename)
}
