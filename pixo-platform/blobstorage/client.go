package blobstorage

import (
	"context"
	"io"
	"net/http"
	"time"
)

type StorageClient interface {
	FindFilesWithName(ctx context.Context, bucketName, prefix, filename string) ([]string, error)
	GetPublicURL(object UploadableObject) string
	GetSignedURL(ctx context.Context, object UploadableObject, options ...Option) (string, error)
	GetChecksum(ctx context.Context, object UploadableObject) (string, error)
	FileExists(ctx context.Context, object UploadableObject) (bool, error)
	UploadFile(ctx context.Context, object UploadableObject, fileReader io.Reader) (string, error)
	Copy(ctx context.Context, src UploadableObject, dest UploadableObject) error
	ReadFile(ctx context.Context, object UploadableObject) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, object UploadableObject) error
	InitResumableUpload(ctx context.Context, object UploadableObject) (*ResumableUploadResponse, error)
}

type Option struct {
	ContentDisposition string
	Lifetime           time.Duration
	Method             string
}

type UploadableObject interface {
	GetBucketName() string
	GetFileLocation() string
	GetTimestamp() int64
}

type SignedURLPartsRequest struct {
	ID        int    `json:"id" binding:"required"`
	Filename  string `json:"filename" binding:"required"`
	NumChunks int    `json:"numChunks" binding:"required"`
}

type ResumableUploadResponse struct {
	UploadURL    string
	Method       string
	SignedHeader http.Header
}
