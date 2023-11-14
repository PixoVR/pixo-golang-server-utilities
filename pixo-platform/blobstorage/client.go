package client

import (
	"context"
	"io"
	"net/http"
)

type StorageClient interface {
	GetSignedURL(ctx context.Context, uploadableObject UploadableObject) (string, error)
	UploadFile(ctx context.Context, uploadableObject UploadableObject, fileReader io.Reader) (string, error)
	InitResumableUpload(ctx context.Context, uploadableObject UploadableObject) (*ResumableUploadResponse, error)
}

type UploadableObject interface {
	GetBucketName() string
	GetUploadDestination() string
	GetFilename() string
}

type SignedURLPartsRequest struct {
	ID        int    `json:"id,required"`
	Filename  string `json:"filename,required"`
	NumChunks int    `json:"numChunks,omitempty"`
}

type ResumableUploadResponse struct {
	UploadURL    string
	Method       string
	SignedHeader http.Header
}
