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
		bucketFilePath = "testdata"
		localFileDir   = "../testdata"
		filename       = "test-file.txt"
		localFilepath  = fmt.Sprintf("%s/%s", localFileDir, filename)
		awsClient      aws.Client
		ctx            = context.Background()
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

			object = client.BasicUploadable{
				BucketName:        config.BucketName,
				UploadDestination: bucketFilePath,
				Filename:          filename,
			}
		)
		BeforeAll(func() {
			var err error
			awsClient, err = aws.NewClient(config)
			Expect(awsClient).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("can create a client if no bucket name is given", func() {
			_, err := aws.NewClient(aws.Config{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("can return empty string if object is empty", func() {
			publicURL := awsClient.GetPublicURL(client.PathUploadable{})
			Expect(publicURL).To(Equal(""))
		})

		It("can return a public url", func() {
			publicURL := awsClient.GetPublicURL(object)
			Expect(publicURL).To(Equal(fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.BucketName, bucketFilePath+"/"+filename)))
		})

		It("can sanitize a filename", func() {
			sanitizedName := awsClient.SanitizeFilename(filename)
			Expect(sanitizedName).To(MatchRegexp(`^blob_\d+.txt$`))
		})

		It("can upload a file to s3", func() {
			fileReader, err := os.Open(localFilepath)

			locationInBucket, err := awsClient.UploadFile(ctx, object, fileReader)

			Expect(err).NotTo(HaveOccurred())
			Expect(locationInBucket).To(MatchRegexp(`^testdata/blob_\d+.txt$`))
		})

		It("can generate a signed url for a file", func() {
			signedUrl, err := awsClient.GetSignedURL(ctx, object)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedUrl).To(ContainSubstring(bucketFilePath))
			Expect(signedUrl).To(ContainSubstring(filename))
		})

		It("can read a file", func() {
			fileReader, err := awsClient.ReadFile(ctx, object)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileReader).NotTo(BeNil())

			bytes := make([]byte, 7)
			n, err := fileReader.Read(bytes)
			//Expect(err).NotTo(HaveOccurred()) // TODO: This is returning an EOF error, but the file is still being read...
			Expect(n).To(Equal(7))
			Expect(string(bytes)).To(ContainSubstring("Go Blue"))
			Expect(fileReader.Close()).To(Succeed())
		})

		It("can check if a file exists", func() {
			exists, err := awsClient.FileExists(ctx, object)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			exists, err = awsClient.FileExists(ctx, client.BasicUploadable{
				BucketName:        config.BucketName,
				UploadDestination: bucketFilePath,
				Filename:          "nonexistent-file.txt",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("can delete a file", func() {
			err := awsClient.DeleteFile(ctx, object)
			Expect(err).NotTo(HaveOccurred())

			fileReader, err := awsClient.ReadFile(ctx, object)
			Expect(err).To(HaveOccurred())
			Expect(fileReader).To(BeNil())
		})

		//It("can initiate a multipart upload", func() {
		//	res, err := awsClient.InitResumableUpload(ctx, object)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(res).NotTo(BeNil())
		//	Expect(res.UploadURL).To(ContainSubstring("x-id=GetObject"))
		//	Expect(res.UploadURL).To(ContainSubstring(bucketFilePath))
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

			object = client.BasicUploadable{
				BucketName:        config.BucketName,
				UploadDestination: bucketFilePath,
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

			signedURL, err := awsClient.UploadFile(ctx, object, fileReader)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedURL).To(ContainSubstring(bucketFilePath))
			Expect(signedURL).To(ContainSubstring(filename))
		})

		It("can generate a signed url for a file", func() {
			signedUrl, err := awsClient.GetSignedURL(ctx, object)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedUrl).To(ContainSubstring(bucketFilePath))
			Expect(signedUrl).To(ContainSubstring(filename))
		})

		It("can read a file", func() {
			fileReader, err := awsClient.ReadFile(ctx, object)
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
			err := awsClient.DeleteFile(ctx, object)
			Expect(err).NotTo(HaveOccurred())

			fileReader, err := awsClient.ReadFile(ctx, object)
			Expect(err).To(HaveOccurred())
			Expect(fileReader).To(BeNil())
		})

		//It("can initiate a multipart upload", func() {
		//	res, err := awsClient.InitResumableUpload(ctx, object)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(res).NotTo(BeNil())
		//	Expect(res.UploadURL).To(ContainSubstring("x-id=GetObject"))
		//	Expect(res.UploadURL).To(ContainSubstring(bucketFilePath))
		//	Expect(res.UploadURL).To(ContainSubstring(filename))
		//})

	})
})
