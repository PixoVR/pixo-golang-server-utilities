package gcs_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/gcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Blob Storage", Ordered, func() {

	var (
		localTestDir   = "../testdata"
		filename       = "test-file.txt"
		localFilepath  = fmt.Sprintf("%s/%s", localTestDir, filename)
		bucketName     = "dev-apex-api-modules"
		bucketFilepath = "testdata"

		gcsClient              gcs.Client
		expectedSignedURLValue = "X-Goog-Algorithm=GOOG4-RSA-SHA256"

		uploadableObject = client.BasicUploadableObject{
			BucketName:        bucketName,
			UploadDestination: bucketFilepath,
			Filename:          filename,
		}
	)

	BeforeAll(func() {
		var err error
		gcsClient, err = gcs.NewClient(gcs.Config{BucketName: bucketName})
		Expect(err).NotTo(HaveOccurred())
		Expect(gcsClient).NotTo(BeNil())
	})

	It("can upload a file", func() {
		fileReader, err := os.Open(localFilepath)
		Expect(err).NotTo(HaveOccurred())

		signedURL, err := gcsClient.UploadFile(context.Background(), uploadableObject, fileReader)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring(bucketFilepath))
		Expect(signedURL).To(ContainSubstring(filename))
	})

	It("can get the signed url for the previously uploaded file", func() {
		signedUrl, err := gcsClient.GetSignedURL(context.Background(), uploadableObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
		Expect(signedUrl).To(ContainSubstring(bucketName))
		Expect(signedUrl).To(ContainSubstring(bucketFilepath))
		Expect(signedUrl).To(ContainSubstring(filename))
	})

	It("can read a file", func() {
		fileReader, err := gcsClient.ReadFile(context.Background(), uploadableObject)
		Expect(err).NotTo(HaveOccurred())
		Expect(fileReader).NotTo(BeNil())

		bytes := make([]byte, 7)
		n, err := fileReader.Read(bytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(7))
		Expect(string(bytes)).To(ContainSubstring("Go Blue"))
		Expect(fileReader.Close()).To(Succeed())
	})

	It("can initiate a multipart upload", func() {
		res, err := gcsClient.InitResumableUpload(context.Background(), uploadableObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.UploadURL).To(ContainSubstring("upload_id="))
		Expect(res.UploadURL).To(ContainSubstring(bucketName))
		Expect(res.UploadURL).To(ContainSubstring(bucketFilepath))
		Expect(res.UploadURL).To(ContainSubstring(filename))
	})

})
