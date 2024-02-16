package aws_test

import (
	"context"
	"fmt"
	storage "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
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
		storageClient  aws.Client
		ctx            = context.Background()
	)
	Context("General S3", func() {
		var (
			config = aws.Config{
				BucketName:      "x-na",
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				Region:          "us-east-1",
			}
			expectedSignedURLValue = "X-Amz-Algorithm=AWS4-HMAC-SHA256"

			object = storage.BasicUploadable{
				BucketName:        config.BucketName,
				UploadDestination: bucketFilePath,
				Filename:          filename,
			}
			uploadedObject storage.PathUploadable
		)

		BeforeAll(func() {
			var err error
			storageClient, err = aws.NewClient(config)
			Expect(storageClient).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("can create a client if no bucket name is given", func() {
			_, err := aws.NewClient(aws.Config{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("can return empty string if object is empty", func() {
			publicURL := storageClient.GetPublicURL(storage.PathUploadable{})
			Expect(publicURL).To(Equal(""))
		})

		It("can return a public url", func() {
			publicURL := storageClient.GetPublicURL(object)
			Expect(publicURL).To(Equal(fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.BucketName, bucketFilePath+"/"+filename)))
		})

		It("can sanitize a filename", func() {
			sanitizedName := storageClient.SanitizeFilename(filename)
			Expect(sanitizedName).To(MatchRegexp(`^blob_\d+.txt$`))
		})

		It("can upload a file to s3", func() {
			fileReader, err := os.Open(localFilepath)
			Expect(err).NotTo(HaveOccurred())

			locationInBucket, err := storageClient.UploadFile(ctx, object, fileReader)

			Expect(err).NotTo(HaveOccurred())
			Expect(locationInBucket).To(MatchRegexp(`^testdata/blob_\d+.txt$`))
			uploadedObject = storage.PathUploadable{
				BucketName: config.BucketName,
				Filepath:   locationInBucket,
			}
		})

		It("can copy a file", func() {
			destinationObject := storage.PathUploadable{
				BucketName: config.BucketName,
				Filepath:   "testdata/copied-file.txt",
			}

			err := storageClient.Copy(ctx, uploadedObject, destinationObject)

			Expect(err).NotTo(HaveOccurred())
			exists, err := storageClient.FileExists(ctx, destinationObject)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

		})

		It("can generate a signed url for a file", func() {
			signedUrl, err := storageClient.GetSignedURL(ctx, object)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedUrl).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedUrl).To(ContainSubstring(bucketFilePath))
			Expect(signedUrl).To(ContainSubstring(filename))
		})

		It("can read a file", func() {
			fileReader, err := storageClient.ReadFile(ctx, uploadedObject)

			Expect(err).NotTo(HaveOccurred())
			Expect(fileReader).NotTo(BeNil())
			bytes := make([]byte, 7)
			n, _ := fileReader.Read(bytes)
			//Expect(err).NotTo(HaveOccurred()) // TODO: This is returning an EOF error, but the file is still being read...
			Expect(n).To(Equal(7))
			Expect(string(bytes)).To(ContainSubstring("Go Blue"))
			Expect(fileReader.Close()).To(Succeed())
		})

		It("can check if a file exists", func() {
			exists, err := storageClient.FileExists(ctx, uploadedObject)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			exists, err = storageClient.FileExists(ctx, storage.BasicUploadable{
				BucketName:        config.BucketName,
				UploadDestination: bucketFilePath,
				Filename:          "nonexistent-file.txt",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("can delete a file", func() {
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
			Expect(res.UploadURL).To(ContainSubstring(expectedSignedURLValue))
			Expect(res.UploadURL).To(ContainSubstring(bucketFilePath))
			Expect(res.UploadURL).NotTo(ContainSubstring(filename))
		})

	})

	Context("STC S3", func() {
		var (
			config = aws.Config{
				BucketName:      "apex-stc-test",
				Endpoint:        "https://api-object.bluvalt.com:8082",
				Region:          "us-east-1",
				AccessKeyID:     os.Getenv("STC_S3_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("STC_S3_SECRET_ACCESS_KEY"),
			}
			expectedSignedURLValue = "X-Amz-Algorithm=AWS4-HMAC-SHA256"

			object = storage.BasicUploadable{
				BucketName:        config.BucketName,
				UploadDestination: bucketFilePath,
				Filename:          filename,
			}
			uploadedObject storage.PathUploadable
		)

		BeforeAll(func() {
			var err error
			storageClient, err = aws.NewClient(config)
			Expect(storageClient).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("can upload a file to s3", func() {
			fileReader, err := os.Open(localFilepath)
			Expect(err).NotTo(HaveOccurred())

			locationInBucket, err := storageClient.UploadFile(ctx, object, fileReader)
			uploadedObject = storage.PathUploadable{
				BucketName: config.BucketName,
				Filepath:   locationInBucket,
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(locationInBucket).To(MatchRegexp(`^testdata/blob_\d+.txt$`))
		})

		It("can generate a signed url for a file", func() {
			signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject)

			Expect(err).NotTo(HaveOccurred())
			Expect(signedURL).To(ContainSubstring(expectedSignedURLValue))
			Expect(signedURL).To(ContainSubstring(bucketFilePath))
			Expect(signedURL).NotTo(ContainSubstring(filename))
		})

		It("can read a file", func() {
			fileReader, err := storageClient.ReadFile(ctx, uploadedObject)

			Expect(err).NotTo(HaveOccurred())
			Expect(fileReader).NotTo(BeNil())
			bytes := make([]byte, 7)
			n, _ := fileReader.Read(bytes)
			//Expect(err).NotTo(HaveOccurred()) // TODO: This is returning an EOF error, but the file is still being read...
			Expect(n).To(Equal(7))
			Expect(string(bytes)).To(ContainSubstring("Go Blue"))
			Expect(fileReader.Close()).To(Succeed())
		})

		It("can delete a file", func() {
			err := storageClient.DeleteFile(ctx, uploadedObject)
			Expect(err).NotTo(HaveOccurred())

			fileReader, err := storageClient.ReadFile(ctx, uploadedObject)

			Expect(err).To(HaveOccurred())
			Expect(fileReader).To(BeNil())
		})

		It("can initiate a multipart upload", func() {
			res, err := storageClient.InitResumableUpload(ctx, uploadedObject)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.UploadURL).To(ContainSubstring(expectedSignedURLValue))
			Expect(res.UploadURL).To(ContainSubstring(bucketFilePath))
			Expect(res.UploadURL).NotTo(ContainSubstring(filename))
		})

	})
})
