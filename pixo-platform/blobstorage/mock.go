package blobstorage

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

type MockStorageClient struct {
	GetPublicURLNumTimesCalled int
	GetPublicURLError          error
	GetPublicURLObjects        []UploadableObject

	GetSignedURLNumTimesCalled int
	GetSignedURLError          error
	GetSignedURLOptions        [][]Option
	GetSignedURLObjects        []UploadableObject

	UploadFileNumTimesCalled int
	UploadFileError          error
	UploadFileObjects        []UploadableObject

	FileExistsNumTimesCalled int
	FileShouldExist          bool
	FileExistsError          error
	FileExistsObjects        []UploadableObject

	CopyNumTimesCalled int
	CopyError          error
	CopySrcObjects     []UploadableObject
	CopyDestObjects    []UploadableObject

	ReadFileNumTimesCalled int
	ReadFileError          error
	ReadFileObjects        []UploadableObject
	ReadFilePath           string

	FindFilesWithNameNumTimesCalled int
	FindFilesWithNameEmpty          bool
	FindFilesWithNameError          error
	FindFilesWithNameQueries        [][]string

	DeleteFileNumTimesCalled int
	DeleteFileError          error
	DeleteFileObjects        []UploadableObject

	InitResumableUploadNumTimesCalled int
	InitResumableUploadError          error
	InitResumableUploadObjects        []UploadableObject
}

var _ StorageClient = (*MockStorageClient)(nil)

func NewMockStorageClient() *MockStorageClient {
	return &MockStorageClient{FileShouldExist: true}
}

func (f *MockStorageClient) Reset() {
	f.GetPublicURLNumTimesCalled = 0
	f.GetPublicURLError = nil
	f.GetPublicURLObjects = nil

	f.GetSignedURLNumTimesCalled = 0
	f.GetSignedURLError = nil
	f.GetSignedURLOptions = nil
	f.GetSignedURLObjects = nil

	f.UploadFileNumTimesCalled = 0
	f.UploadFileError = nil
	f.UploadFileObjects = nil

	f.FileExistsNumTimesCalled = 0
	f.FileShouldExist = true
	f.FileExistsError = nil
	f.FileExistsObjects = nil

	f.CopyNumTimesCalled = 0
	f.CopyError = nil
	f.CopySrcObjects = nil
	f.CopyDestObjects = nil

	f.ReadFileNumTimesCalled = 0
	f.ReadFileError = nil
	f.ReadFileObjects = nil
	f.ReadFilePath = ""

	f.FindFilesWithNameNumTimesCalled = 0
	f.FindFilesWithNameError = nil

	f.DeleteFileNumTimesCalled = 0
	f.DeleteFileError = nil
	f.DeleteFileObjects = nil

	f.InitResumableUploadNumTimesCalled = 0
	f.InitResumableUploadError = nil
	f.InitResumableUploadObjects = nil
}

func (f *MockStorageClient) GetPublicURL(object UploadableObject) string {
	f.GetPublicURLNumTimesCalled++
	f.GetPublicURLObjects = append(f.GetPublicURLObjects, object)

	if f.GetPublicURLError != nil {
		return ""
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", object.GetBucketName(), object.GetFileLocation())
}

func (f *MockStorageClient) GetSignedURL(ctx context.Context, object UploadableObject, options ...Option) (string, error) {
	f.GetSignedURLNumTimesCalled++
	f.GetSignedURLObjects = append(f.GetSignedURLObjects, object)
	f.GetSignedURLOptions = append(f.GetSignedURLOptions, options)

	if f.GetSignedURLError != nil {
		return "", f.GetSignedURLError
	}

	signedURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s?X-Goog-Algorithm=GOOG4-RSA-SHA256&X-Goog-Credential=credential&X-Goog-Date=20210101T000000Z&X-Goog-Expires=3600&X-Goog-SignedHeaders=host&X-Goog-Signature=signature", object.GetBucketName(), object.GetFileLocation())
	return signedURL, nil
}

func (f *MockStorageClient) GetChecksum(ctx context.Context, object UploadableObject) (string, error) {
	return fmt.Sprint(md5.New().Sum([]byte(object.GetFileLocation()))), nil
}

func (f *MockStorageClient) UploadFile(ctx context.Context, object UploadableObject, fileReader io.Reader) (string, error) {
	f.UploadFileNumTimesCalled++
	f.UploadFileObjects = append(f.UploadFileObjects, object)

	if f.UploadFileError != nil {
		return "", f.UploadFileError
	}

	return object.GetFileLocation(), nil
}

func (f *MockStorageClient) FileExists(ctx context.Context, object UploadableObject) (bool, error) {
	f.FileExistsNumTimesCalled++
	f.FileExistsObjects = append(f.FileExistsObjects, object)

	if f.FileExistsError != nil {
		return false, f.FileExistsError
	}

	return f.FileShouldExist, nil
}

func (f *MockStorageClient) FindFilesWithName(ctx context.Context, bucketName, prefix, filename string) ([]string, error) {
	f.FindFilesWithNameNumTimesCalled++
	f.FindFilesWithNameQueries = append(f.FindFilesWithNameQueries, []string{bucketName, prefix, filename})

	if f.FindFilesWithNameError != nil {
		return nil, f.FindFilesWithNameError
	}

	if f.FindFilesWithNameEmpty {
		return nil, nil
	}

	return []string{"one/" + filename, "two/" + filename}, nil
}

func (f *MockStorageClient) Copy(ctx context.Context, src, dest UploadableObject) error {
	f.CopyNumTimesCalled++
	f.CopyDestObjects = append(f.CopyDestObjects, dest)
	f.CopySrcObjects = append(f.CopySrcObjects, src)

	if f.CopyError != nil {
		return f.CopyError
	}

	return nil
}

// ReadFile reads a file from the storage client as a mock
//
// If ReadFilePath is set, it will read the file from the path specified
// and return an open io.ReadCloser to the file. Note: make sure to close the io.ReadCloser
func (f *MockStorageClient) ReadFile(ctx context.Context, object UploadableObject) (io.ReadCloser, error) {
	f.ReadFileNumTimesCalled++
	f.ReadFileObjects = append(f.ReadFileObjects, object)

	if f.ReadFileError != nil {
		return nil, f.ReadFileError
	}

	if f.ReadFilePath != "" {
		file, err := os.Open(f.ReadFilePath)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	reader := strings.NewReader("test")
	readCloser := io.NopCloser(reader)

	return readCloser, nil
}

func (f *MockStorageClient) DeleteFile(ctx context.Context, object UploadableObject) error {
	f.DeleteFileNumTimesCalled++
	f.DeleteFileObjects = append(f.DeleteFileObjects, object)

	if f.DeleteFileError != nil {
		return f.DeleteFileError
	}

	return nil
}

func (f *MockStorageClient) InitResumableUpload(ctx context.Context, object UploadableObject) (*ResumableUploadResponse, error) {
	f.InitResumableUploadNumTimesCalled++
	f.InitResumableUploadObjects = append(f.InitResumableUploadObjects, object)

	if f.InitResumableUploadError != nil {
		return nil, f.InitResumableUploadError
	}

	return nil, nil
}
