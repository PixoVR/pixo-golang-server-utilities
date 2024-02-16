package blobstorage

import (
	"context"
	"io"
)

type MockStorageClient struct {
	GetPublicURLNumTimesCalled int
	GetPublicURLError          error

	GetSignedURLNumTimesCalled int
	GetSignedURLError          error

	UploadFileNumTimesCalled int
	UploadFileError          error

	FileExistsNumTimesCalled int
	FileShouldExist          bool
	FileExistsError          error

	CopyNumTimesCalled int
	CopyError          error

	ReadFileNumTimesCalled int
	ReadFileError          error

	DeleteFileNumTimesCalled int
	DeleteFileError          error

	InitResumableUploadNumTimesCalled int
	InitResumableUploadError          error
}

func NewMockStorageClient() *MockStorageClient {
	return &MockStorageClient{}
}

func (f *MockStorageClient) GetPublicURL(object UploadableObject) string {
	f.GetPublicURLNumTimesCalled++

	if f.GetPublicURLError != nil {
		return ""
	}

	return "https://storage.googleapis.com/bucket/test-file.txt"
}

func (f *MockStorageClient) GetSignedURL(ctx context.Context, object UploadableObject) (string, error) {
	f.GetSignedURLNumTimesCalled++

	if f.GetSignedURLError != nil {
		return "", f.GetSignedURLError
	}

	return "fake-signed-url", nil
}

func (f *MockStorageClient) UploadFile(ctx context.Context, object UploadableObject, fileReader io.Reader) (string, error) {
	f.UploadFileNumTimesCalled++

	if f.UploadFileError != nil {
		return "", f.UploadFileError
	}

	return "fake-url", nil
}

func (f *MockStorageClient) FileExists(ctx context.Context, object UploadableObject) (bool, error) {
	f.FileExistsNumTimesCalled++

	if f.FileExistsError != nil {
		return false, f.FileExistsError
	}

	return f.FileShouldExist, nil
}

func (f *MockStorageClient) Copy(ctx context.Context, source, destination UploadableObject) error {
	f.CopyNumTimesCalled++

	if f.CopyError != nil {
		return f.CopyError
	}

	return nil
}

func (f *MockStorageClient) ReadFile(ctx context.Context, object UploadableObject) (io.ReadCloser, error) {
	f.ReadFileNumTimesCalled++

	if f.ReadFileError != nil {
		return nil, f.ReadFileError
	}

	return nil, nil
}

func (f *MockStorageClient) DeleteFile(ctx context.Context, object UploadableObject) error {
	f.DeleteFileNumTimesCalled++

	if f.DeleteFileError != nil {
		return f.DeleteFileError
	}

	return nil
}

func (f *MockStorageClient) InitResumableUpload(ctx context.Context, object UploadableObject) (*ResumableUploadResponse, error) {
	f.InitResumableUploadNumTimesCalled++

	if f.InitResumableUploadError != nil {
		return nil, f.InitResumableUploadError
	}

	return nil, nil
}
