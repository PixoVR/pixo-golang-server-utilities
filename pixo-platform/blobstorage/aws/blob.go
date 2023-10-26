package aws

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/blob"
)

type Blob struct {
	blob blob.Blob
}

func (b *Blob) Filepath() string {
	return b.blob.Filepath()
}

func (b *Blob) GetSignedURL() (string, error) {
	return GetSignedURL(b.blob.BucketName(), b.blob.Filepath())
}

func (b *Blob) GetResumableUploadURL() (blob.ResumableUploadResponse, error) {
	return InitResumableUpload(b.blob.BucketName(), b.blob.Filepath())
}
