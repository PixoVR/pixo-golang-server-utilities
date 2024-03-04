package gcs

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"os"
	"strings"
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
	return strings.ReplaceAll(blobstorage.ParseFileLocationFromLink(p.Path), p.GetBucketName(), "")
}
