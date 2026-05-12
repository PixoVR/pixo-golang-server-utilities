package blobstorage_test

import (
	"bytes"
	"context"
	"errors"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockStorageClient", func() {

	var (
		mock *blobstorage.MockStorageClient
		ctx  = context.Background()
	)

	BeforeEach(func() {
		mock = blobstorage.NewMockStorageClient()
	})

	Describe("UploadRawFile", func() {

		It("returns the exact file location without sanitization", func() {
			object := blobstorage.PathUploadable{
				BucketName: "test-bucket",
				Filepath:   "43/images/textures/test.txt",
			}

			location, err := mock.UploadRawFile(ctx, object, bytes.NewReader([]byte("content")))

			Expect(err).NotTo(HaveOccurred())
			Expect(location).To(Equal("43/images/textures/test.txt"))
			Expect(mock.UploadRawFileNumTimesCalled).To(Equal(1))
			Expect(mock.UploadRawFileObjects).To(HaveLen(1))
			Expect(mock.UploadRawFileObjects[0].GetFileLocation()).To(Equal("43/images/textures/test.txt"))
		})

		It("returns an error when configured to fail", func() {
			mock.UploadRawFileError = errors.New("upload failed")

			location, err := mock.UploadRawFile(ctx, blobstorage.PathUploadable{
				BucketName: "test-bucket",
				Filepath:   "some/path/file.txt",
			}, bytes.NewReader([]byte("content")))

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("upload failed"))
			Expect(location).To(BeEmpty())
			Expect(mock.UploadRawFileNumTimesCalled).To(Equal(1))
		})

		It("tracks multiple calls", func() {
			objects := []blobstorage.PathUploadable{
				{BucketName: "bucket", Filepath: "path/one.txt"},
				{BucketName: "bucket", Filepath: "path/two.txt"},
				{BucketName: "bucket", Filepath: "path/three.txt"},
			}

			for _, obj := range objects {
				_, err := mock.UploadRawFile(ctx, obj, bytes.NewReader([]byte("data")))
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(mock.UploadRawFileNumTimesCalled).To(Equal(3))
			Expect(mock.UploadRawFileObjects).To(HaveLen(3))
		})

		It("is reset when Reset is called", func() {
			mock.UploadRawFileError = errors.New("error")
			_, _ = mock.UploadRawFile(ctx, blobstorage.PathUploadable{
				BucketName: "bucket",
				Filepath:   "file.txt",
			}, bytes.NewReader([]byte("data")))

			mock.Reset()

			Expect(mock.UploadRawFileNumTimesCalled).To(Equal(0))
			Expect(mock.UploadRawFileError).To(BeNil())
			Expect(mock.UploadRawFileObjects).To(BeNil())
		})

		It("preserves nested directory paths", func() {
			object := blobstorage.PathUploadable{
				BucketName: "cas-bucket",
				Filepath:   "42/deeply/nested/directory/structure/file.png",
			}

			location, err := mock.UploadRawFile(ctx, object, bytes.NewReader([]byte("image data")))

			Expect(err).NotTo(HaveOccurred())
			Expect(location).To(Equal("42/deeply/nested/directory/structure/file.png"))
		})

		It("does not interfere with UploadFile tracking", func() {
			object := blobstorage.PathUploadable{
				BucketName: "bucket",
				Filepath:   "path/file.txt",
			}

			_, err := mock.UploadRawFile(ctx, object, bytes.NewReader([]byte("raw")))
			Expect(err).NotTo(HaveOccurred())

			_, err = mock.UploadFile(ctx, object, bytes.NewReader([]byte("sanitized")))
			Expect(err).NotTo(HaveOccurred())

			Expect(mock.UploadRawFileNumTimesCalled).To(Equal(1))
			Expect(mock.UploadFileNumTimesCalled).To(Equal(1))
		})
	})
})
