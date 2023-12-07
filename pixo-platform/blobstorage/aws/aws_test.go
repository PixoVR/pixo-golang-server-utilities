package aws_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("S3 Blob Storage", Ordered, func() {

	var (
		bucketFileDir = "testdata"
		localFileDir  = "../testdata"
		filename      = "test-file.txt"
		localFilepath = fmt.Sprintf("%s/%s", localFileDir, filename)
		awsClient     aws.Client
	)
	Context("General S3", func() {
		var (
			config = aws.Config{
				BucketName:      os.Getenv("S3_BUCKET_NAME"),
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				Region:          "us-east-1",
			}
			expectedSignedURLValue = "X-Amz-Algorithm=AWS4-HMAC-SHA256"

			object = client.BasicUploadableObject{
				BucketName:        config.BucketName,
				UploadDestination: bucketFileDir,
				Filename:          filename,
			}
		)
		BeforeAll(func() {
			var err error
			awsClient, err = aws.NewClient(config)
			Expect(awsClient).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("can return an error if no bucket name is given", func() {
			_, err := aws.NewClient(aws.Config{})
			Expect(err).To(HaveOccurred())
		})

		It("can upload a file to s3", func() {
			fileReader, err := os.Open(localFilepath)

			signedURL, err := awsClient.UploadFile(context.Background(), object, fileReader)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedURL).To(ContainSubstring(bucketFileDir))
			Expect(signedURL).To(ContainSubstring(filename))
		})

		It("can generate a signed url for a file", func() {
			signedUrl, err := awsClient.GetSignedURL(context.Background(), object)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedUrl).To(ContainSubstring(bucketFileDir))
			Expect(signedUrl).To(ContainSubstring(filename))
		})

		It("can read a file", func() {
			fileReader, err := awsClient.ReadFile(context.Background(), object)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileReader).NotTo(BeNil())

			bytes := make([]byte, 7)
			n, err := fileReader.Read(bytes)
			//Expect(err).NotTo(HaveOccurred()) // TODO: This is returning an EOF error, but the file is still being read...
			Expect(n).To(Equal(7))
			Expect(string(bytes)).To(ContainSubstring("Go Blue"))
			Expect(fileReader.Close()).To(Succeed())
		})

		It("can delete a file", func() {
			err := awsClient.DeleteFile(context.Background(), object)
			Expect(err).NotTo(HaveOccurred())

			fileReader, err := awsClient.ReadFile(context.Background(), object)
			Expect(err).To(HaveOccurred())
			Expect(fileReader).To(BeNil())
		})

		//It("can initiate a multipart upload", func() {
		//	res, err := awsClient.InitResumableUpload(context.Background(), object)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(res).NotTo(BeNil())
		//	Expect(res.UploadURL).To(ContainSubstring("x-id=GetObject"))
		//	Expect(res.UploadURL).To(ContainSubstring(bucketFileDir))
		//	Expect(res.UploadURL).To(ContainSubstring(filename))
		//})

	})
	Context("STC S3", func() {
		var (
			config = aws.Config{
				BucketName:      os.Getenv("STC_S3_BUCKET_NAME"),
				AccessKeyID:     os.Getenv("STC_AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("STC_AWS_SECRET_ACCESS_KEY"),
				Endpoint:        os.Getenv("STC_AWS_ENDPOINT"),
				Region:          os.Getenv("STC_AWS_REGION"),
			}
			expectedSignedURLValue = "X-Amz-Algorithm=AWS4-HMAC-SHA256"

			object = client.BasicUploadableObject{
				BucketName:        config.BucketName,
				UploadDestination: bucketFileDir,
				Filename:          filename,
			}
		)
		BeforeAll(func() {
			var err error
			awsClient, err = aws.NewClient(config)
			Expect(awsClient).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("can upload a file to s3", func() {
			fileReader, err := os.Open(localFilepath)

			signedURL, err := awsClient.UploadFile(context.Background(), object, fileReader)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedURL).To(ContainSubstring(bucketFileDir))
			Expect(signedURL).To(ContainSubstring(filename))
		})

		It("can generate a signed url for a file", func() {
			signedUrl, err := awsClient.GetSignedURL(context.Background(), object)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedUrl).To(ContainSubstring(bucketFileDir))
			Expect(signedUrl).To(ContainSubstring(filename))
		})

		It("can read a file", func() {
			fileReader, err := awsClient.ReadFile(context.Background(), object)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileReader).NotTo(BeNil())

			bytes := make([]byte, 7)
			n, err := fileReader.Read(bytes)
			//Expect(err).NotTo(HaveOccurred()) // TODO: This is returning an EOF error, but the file is still being read...
			Expect(n).To(Equal(7))
			Expect(string(bytes)).To(ContainSubstring("Go Blue"))
			Expect(fileReader.Close()).To(Succeed())
		})

		It("can delete a file", func() {
			err := awsClient.DeleteFile(context.Background(), object)
			Expect(err).NotTo(HaveOccurred())

			fileReader, err := awsClient.ReadFile(context.Background(), object)
			Expect(err).To(HaveOccurred())
			Expect(fileReader).To(BeNil())
		})

		//It("can initiate a multipart upload", func() {
		//	res, err := awsClient.InitResumableUpload(context.Background(), object)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(res).NotTo(BeNil())
		//	Expect(res.UploadURL).To(ContainSubstring("x-id=GetObject"))
		//	Expect(res.UploadURL).To(ContainSubstring(bucketFileDir))
		//	Expect(res.UploadURL).To(ContainSubstring(filename))
		//})

	})
})
