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

	FindFilesWithNameNumTimesCalled int
	FindFilesWithNameEmpty          bool
	FindFilesWithNameError          error

	DeleteFileNumTimesCalled int
	DeleteFileError          error

	InitResumableUploadNumTimesCalled int
	InitResumableUploadError          error
}

var _ StorageClient = (*MockStorageClient)(nil)

func NewMockStorageClient() *MockStorageClient {
	return &MockStorageClient{FileShouldExist: true}
}

func (f *MockStorageClient) Reset() {
	f.GetPublicURLNumTimesCalled = 0
	f.GetPublicURLError = nil

	f.GetSignedURLNumTimesCalled = 0
	f.GetSignedURLError = nil

	f.UploadFileNumTimesCalled = 0
	f.UploadFileError = nil

	f.FileExistsNumTimesCalled = 0
	f.FileShouldExist = true
	f.FileExistsError = nil

	f.CopyNumTimesCalled = 0
	f.CopyError = nil

	f.ReadFileNumTimesCalled = 0
	f.ReadFileError = nil

	f.FindFilesWithNameNumTimesCalled = 0
	f.FindFilesWithNameError = nil

	f.DeleteFileNumTimesCalled = 0
	f.DeleteFileError = nil

	f.InitResumableUploadNumTimesCalled = 0
	f.InitResumableUploadError = nil
}

func (f *MockStorageClient) GetPublicURL(object UploadableObject) string {
	f.GetPublicURLNumTimesCalled++

	if f.GetPublicURLError != nil {
		return ""
	}

	return "https://storage.googleapis.com/bucket/test-file.txt"
}

func (f *MockStorageClient) GetSignedURL(ctx context.Context, object UploadableObject, options ...Option) (string, error) {
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

func (f *MockStorageClient) FindFilesWithName(ctx context.Context, bucketName, prefix, filename string) ([]string, error) {
	f.FindFilesWithNameNumTimesCalled++

	if f.FindFilesWithNameError != nil {
		return nil, f.FindFilesWithNameError
	}

	if f.FindFilesWithNameEmpty {
		return nil, nil
	}

	return []string{"one/" + filename, "two/" + filename}, nil
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
