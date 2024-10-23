package blobstorage

import "time"

type BasicUploadable struct {
	BucketName        string
	UploadDestination string
	Filename          string
	Timestamp         int64
}

func (b BasicUploadable) GetBucketName() string {
	return b.BucketName
}

func (b BasicUploadable) GetFileLocation() string {
	return b.UploadDestination + "/" + b.Filename
}

func (b BasicUploadable) GetTimestamp() int64 {
	if b.Timestamp != 0 {
		return b.Timestamp
	}
	return time.Now().Unix()
}

type PathUploadable struct {
	BucketName string
	Filepath   string
	Timestamp  *time.Time
}

func (p PathUploadable) GetBucketName() string {
	return p.BucketName
}

func (p PathUploadable) GetFileLocation() string {
	return p.Filepath
}

func (p PathUploadable) GetTimestamp() int64 {
	if p.Timestamp != nil {
		return p.Timestamp.Unix()
	}
	return 0
}
