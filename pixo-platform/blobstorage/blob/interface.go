package blob

import "net/http"

type StorageBlob interface {
	Filepath() string
	GetSignedURL() (string, error)
	GetResumableUploadURL() (ResumableUploadResponse, error)
}

type ResumableUploadResponse struct {
	UploadURL    string
	Method       string
	SignedHeader http.Header
}
