package blobstorage

import "os"

type BasicUploadable struct {
	BucketName        string
	UploadDestination string
	Filename          string
}

func (b BasicUploadable) GetBucketName() string {
	return b.BucketName
}

func (b BasicUploadable) GetFileLocation() string {
	return b.UploadDestination + "/" + b.Filename
}

type PathUploadable struct {
	BucketName string
	Filepath   string
}

func (p PathUploadable) GetBucketName() string {
	return p.BucketName
}

func (p PathUploadable) GetFileLocation() string {
	return p.Filepath
}

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
	return ParseFileLocationFromLink(p.Path)
}
