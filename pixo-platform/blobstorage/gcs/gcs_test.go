package gcs_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	storage "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/gcs"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Google Cloud Storage", Ordered, func() {
	var (
		filename       = "test-file.txt"
		bucketName     = "dev-apex-primary-api-modules"
		bucketFilepath = "testdata"

		storageClient          gcs.Client
		expectedSignedURLValue = "X-Goog-Algorithm=GOOG4-RSA-SHA256"

		object = storage.BasicUploadable{
			BucketName:        bucketName,
			UploadDestination: bucketFilepath,
			Filename:          filename,
		}
		ctx = context.Background()

		uploadedObject storage.PathUploadable
	)

	BeforeAll(func() {
		file, err := os.Create(filename)
		Expect(err).NotTo(HaveOccurred())
		_, err = file.WriteString("Go Blue")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		storageClient, err = gcs.NewClient(gcs.Config{BucketName: bucketName})
		Expect(err).NotTo(HaveOccurred())
		Expect(storageClient).NotTo(BeNil())
	})

	AfterAll(func() {
		Expect(os.Remove(filename)).NotTo(HaveOccurred())
	})

	It("can return empty string if object is empty", func() {
		publicURL := storageClient.GetPublicURL(storage.PathUploadable{})
		Expect(publicURL).To(Equal(""))
	})

	It("can format a public url", func() {
		publicURL := storageClient.GetPublicURL(object)
		Expect(publicURL).To(Equal("https://storage.googleapis.com/dev-apex-primary-api-modules/testdata/test-file.txt"))
	})

	It("can sanitize a filename", func() {
		Expect(storageClient.SanitizeFilename(filename)).To(MatchRegexp(`^blob_\d+.txt$`))
		Expect(storageClient.SanitizeFilename("model/thumbnails/")).To(MatchRegexp(`^model/thumbnails/blob_\d+$`))
		Expect(storageClient.SanitizeFilename("model/thumbnails/file.txt")).To(MatchRegexp(`^model/thumbnails/blob_\d+.txt$`))
	})

	It("can upload a file", func() {
		fileReader, err := os.Open(filename)
		Expect(err).NotTo(HaveOccurred())

		locationInBucket, err := storageClient.UploadFile(ctx, object, fileReader)

		Expect(err).NotTo(HaveOccurred())
		Expect(locationInBucket).To(MatchRegexp(`^testdata/blob_\d+.txt$`))
		uploadedObject = storage.PathUploadable{
			BucketName: bucketName,
			Filepath:   locationInBucket,
		}
	})

	It("can copy a file", func() {
		destinationObject := storage.PathUploadable{
			BucketName: bucketName,
			Filepath:   "testdata/copied-file.txt",
		}

		err := storageClient.Copy(ctx, uploadedObject, destinationObject)

		Expect(err).NotTo(HaveOccurred())
		exists, err := storageClient.FileExists(ctx, destinationObject)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())
	})

	It("can check if a file exists", func() {
		exists, err := storageClient.FileExists(ctx, uploadedObject)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		exists, err = storageClient.FileExists(ctx, storage.PathUploadable{
			BucketName: bucketName,
			Filepath:   "testdata/does-not-exist.txt",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("can crawl the bucket to find a file by filename", func() {
		filename := storage.GetFilenameFromLocation(uploadedObject.Filepath)

		locations, err := storageClient.FindFilesWithName(ctx, uploadedObject.BucketName, "", filename)

		Expect(err).NotTo(HaveOccurred())
		Expect(locations).To(HaveLen(1))
		Expect(locations[0]).To(Equal(uploadedObject.Filepath))
	})

	It("can get the signed url for the previously uploaded file", func() {
		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring(bucketFilepath))
		Expect(signedURL).NotTo(ContainSubstring(filename))
		httpClient := resty.New()
		response, err := httpClient.R().Get(signedURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeNil())
		Expect(response.StatusCode()).To(Equal(http.StatusOK))
		reader := bytes.NewReader(response.Body())
		Expect(reader).NotTo(BeNil())
		data, err := io.ReadAll(reader)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("Go Blue"))
	})

	It("can read a file", func() {
		fileReader, err := storageClient.ReadFile(ctx, uploadedObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(fileReader).NotTo(BeNil())
		bytes := make([]byte, 7)
		n, err := fileReader.Read(bytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(7))
		Expect(string(bytes)).To(ContainSubstring("Go Blue"))
		Expect(fileReader.Close()).To(Succeed())
	})

	It("can delete a file", func() {
		time.Sleep(1 * time.Second) // wait 1 second to allow for retention policy to be met

		err := storageClient.DeleteFile(ctx, uploadedObject)

		Expect(err).NotTo(HaveOccurred())
		fileReader, err := storageClient.ReadFile(ctx, uploadedObject)
		Expect(err).To(HaveOccurred())
		Expect(fileReader).To(BeNil())
	})

	It("can initiate a multipart upload", func() {
		res, err := storageClient.InitResumableUpload(ctx, object)

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.UploadURL).To(ContainSubstring("upload_id="))
		Expect(res.UploadURL).To(ContainSubstring(bucketName))
		Expect(res.UploadURL).To(ContainSubstring(bucketFilepath))
		Expect(res.UploadURL).To(ContainSubstring(filename))
	})
	It("can set the content disposition in the signed url", func() {
		options := []blobstorage.Option{
			{
				ContentDisposition: "attachment",
			},
		}
		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject, options...)
		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring("attachment"))
	})

	It("can set the expiration for the signed url", func() {
		expireTime := time.Now().Add(5 * time.Hour)
		expectedTime := expireTime.Format("20060102T150405Z")
		options := []blobstorage.Option{
			{
				Expires: &expireTime,
			},
		}
		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject, options...)
		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring(expectedTime))
	})
})
