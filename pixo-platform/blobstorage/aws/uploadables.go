package aws

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"os"
)

var _ blobstorage.UploadableObject = (*DefaultPublicUploadable)(nil)
var _ blobstorage.UploadableObject = (*DefaultPrivateUploadable)(nil)

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

func (p DefaultPublicUploadable) GetTimestamp() int64 {
	return 0
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
