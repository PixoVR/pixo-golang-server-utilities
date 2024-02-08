package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type StorageClient interface {
	GetSignedURL(ctx context.Context, object UploadableObject) (string, error)
	FileExists(ctx context.Context, object UploadableObject) (bool, error)
	UploadFile(ctx context.Context, object UploadableObject, fileReader io.Reader) (string, error)
	ReadFile(ctx context.Context, object UploadableObject) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, object UploadableObject) error
	InitResumableUpload(ctx context.Context, object UploadableObject) (*ResumableUploadResponse, error)
}

type UploadableObject interface {
	GetBucketName() string
	GetUploadDestination() string
	GetFileLocation() string
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

func GetFullPath(object UploadableObject) string {
	fileDest := object.GetUploadDestination()
	if fileDest == "" {
		return object.GetFileLocation()
	}

	fullPath := fmt.Sprintf("%s/%s", fileDest, object.GetFileLocation())
	return fullPath
}
