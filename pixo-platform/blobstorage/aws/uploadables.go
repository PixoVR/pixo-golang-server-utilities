package aws

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
	return os.Getenv("S3_STORAGE_PUBLIC")
}

func (p DefaultPublicUploadable) GetFileLocation() string {
	return blobstorage.ParseFileLocationFromLink(p.Path)
}

type DefaultPrivateUploadable struct {
	DefaultPublicUploadable
}

func PrivateUploadable(fileLocation string) DefaultPrivateUploadable {
	return DefaultPrivateUploadable{DefaultPublicUploadable{Path: fileLocation}}
}

func (p DefaultPrivateUploadable) GetBucketName() string {
	return os.Getenv("S3_STORAGE_PRIVATE")
}

func (p DefaultPrivateUploadable) GetTimestamp() int64 {
	return 0
}
