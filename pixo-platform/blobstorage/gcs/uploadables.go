package gcs

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"os"
)

type DefaultPublicUploadable struct {
	Path string
}

func PublicUploadable(fileLocation string) DefaultPublicUploadable {
	return DefaultPublicUploadable{Path: fileLocation}
}

func (p DefaultPublicUploadable) GetBucketName() string {
	return os.Getenv("GOOGLE_STORAGE_PUBLIC")
}

func (p DefaultPublicUploadable) GetFileLocation() string {
	return blobstorage.ParseFileLocationFromLink(p.Path)
}

func (p DefaultPublicUploadable) GetTimestamp() int64 {
	return 0
}
