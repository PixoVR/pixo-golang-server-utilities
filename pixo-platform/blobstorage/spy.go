package client

import (
	"context"
	"io"
)

type StorageClientSpy struct {
	GetPublicURLNumTimesCalled        int
	GetSignedURLNumTimesCalled        int
	UploadFileNumTimesCalled          int
	FileExistsNumTimesCalled          int
	CopyNumTimesCalled                int
	ReadFileNumTimesCalled            int
	DeleteFileNumTimesCalled          int
	InitResumableUploadNumTimesCalled int
}

func NewStorageClientSpy() *StorageClientSpy {
	return &StorageClientSpy{}
}

func (f *StorageClientSpy) GetPublicURL(object UploadableObject) string {
	f.GetPublicURLNumTimesCalled++
	return "https://storage.googleapis.com/bucket/test-file.txt"
}

func (f *StorageClientSpy) GetSignedURL(ctx context.Context, object UploadableObject) (string, error) {
	f.GetSignedURLNumTimesCalled++
	return "fake-signed-url", nil
}

func (f *StorageClientSpy) UploadFile(ctx context.Context, object UploadableObject, fileReader io.Reader) (string, error) {
	f.UploadFileNumTimesCalled++
	return "fake-url", nil
}

func (f *StorageClientSpy) FileExists(ctx context.Context, object UploadableObject) (bool, error) {
	f.FileExistsNumTimesCalled++
	return true, nil
}

func (f *StorageClientSpy) Copy(ctx context.Context, source, destination UploadableObject) error {
	f.CopyNumTimesCalled++
	return nil
}

func (f *StorageClientSpy) ReadFile(ctx context.Context, object UploadableObject) (io.ReadCloser, error) {
	f.ReadFileNumTimesCalled++
	return nil, nil
}

func (f *StorageClientSpy) DeleteFile(ctx context.Context, object UploadableObject) error {
	f.DeleteFileNumTimesCalled++
	return nil
}

func (f *StorageClientSpy) InitResumableUpload(ctx context.Context, object UploadableObject) (*ResumableUploadResponse, error) {
	f.InitResumableUploadNumTimesCalled++
	return nil, nil
}
