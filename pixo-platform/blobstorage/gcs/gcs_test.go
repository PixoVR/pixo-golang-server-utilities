package gcs_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/gcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("Blob Storage", Ordered, func() {

	var (
		localTestDir   = "../testdata"
		filename       = "test-file.txt"
		localFilepath  = fmt.Sprintf("%s/%s", localTestDir, filename)
		bucketName     = "dev-apex-primary-api-modules"
		bucketFilepath = "testdata"

		gcsClient              gcs.Client
		expectedSignedURLValue = "X-Goog-Algorithm=GOOG4-RSA-SHA256"

		object = client.BasicUploadableObject{
			BucketName:        bucketName,
			UploadDestination: bucketFilepath,
			Filename:          filename,
		}
		ctx = context.Background()
	)

	BeforeAll(func() {
		var err error
		gcsClient, err = gcs.NewClient(gcs.Config{BucketName: bucketName})
		Expect(err).NotTo(HaveOccurred())
		Expect(gcsClient).NotTo(BeNil())
	})

	It("can format a public url", func() {
		publicURL := gcsClient.GetPublicURL(object)
		Expect(publicURL).To(Equal("https://storage.googleapis.com/dev-apex-primary-api-modules/testdata/test-file.txt"))
	})

	It("can upload a file", func() {
		fileReader, err := os.Open(localFilepath)
		Expect(err).NotTo(HaveOccurred())

		signedURL, err := gcsClient.UploadFile(ctx, object, fileReader)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring(bucketFilepath))
		Expect(signedURL).To(ContainSubstring(filename))
	})

	It("can check if a file exists", func() {
		exists, err := gcsClient.FileExists(ctx, object)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		exists, err = gcsClient.FileExists(ctx, client.BasicUploadableObject{
			BucketName:        bucketName,
			UploadDestination: bucketFilepath,
			Filename:          "nonexistent-file.txt",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("can get the signed url for the previously uploaded file", func() {
		signedUrl, err := gcsClient.GetSignedURL(ctx, object)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedUrl).To(ContainSubstring(bucketName))
		Expect(signedUrl).To(ContainSubstring(bucketFilepath))
		Expect(signedUrl).To(ContainSubstring(filename))
	})

	It("can read a file", func() {
		fileReader, err := gcsClient.ReadFile(ctx, object)
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
		err := gcsClient.DeleteFile(ctx, object)
		Expect(err).NotTo(HaveOccurred())
		fileReader, err := gcsClient.ReadFile(ctx, object)
		Expect(err).To(HaveOccurred())
		Expect(fileReader).To(BeNil())
	})

	It("can initiate a multipart upload", func() {
		res, err := gcsClient.InitResumableUpload(ctx, object)

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.UploadURL).To(ContainSubstring("upload_id="))
		Expect(res.UploadURL).To(ContainSubstring(bucketName))
		Expect(res.UploadURL).To(ContainSubstring(bucketFilepath))
		Expect(res.UploadURL).To(ContainSubstring(filename))
	})

})
