package client

import (
	"io"
	"net/http"
)

type StorageClient interface {
	UploadFile(uploadableObject UploadableObject, fileReader io.Reader) (string, error)
	InitResumableUpload(uploadableObject UploadableObject) (*ResumableUploadResponse, error)
}

type UploadableObject interface {
	GetBucketName() string
	GetUploadDestination() string
	GetFilename() string
}

type ResumableUploadResponse struct {
	UploadURL    string
	Method       string
	SignedHeader http.Header
}
