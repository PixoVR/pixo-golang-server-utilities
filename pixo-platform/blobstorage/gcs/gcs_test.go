package gcs_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/alicebob/miniredis/v2"
	"github.com/onsi/gomega/types"
	"github.com/redis/go-redis/v9"
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
		bucketName     = config.GetEnvOrReturn("GCS_BUCKET_NAME", "pixo-test-bucket")
		bucketFilepath = "testdata"

		c                           *redis.Client
		storageClient               gcs.Client
		googleAlgorithm             = "X-Goog-Algorithm=GOOG4-RSA-SHA256"
		googleDateParam             = "X-Goog-Expires="
		googleExpiresAlgorithmParam = "X-Goog-Expires="

		object = storage.BasicUploadable{
			BucketName:        bucketName,
			UploadDestination: bucketFilepath,
			Filename:          filename,
			Timestamp:         time.Now().Unix(),
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

		s, err := miniredis.Run()
		Expect(err).NotTo(HaveOccurred())

		c = redis.NewClient(&redis.Options{
			Addr:     s.Addr(),
			Password: "",
			DB:       0,
		})

		storageClient, err = gcs.NewClient(gcs.Config{
			BucketName: bucketName,
			Cache:      c,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(storageClient).NotTo(BeNil())
	})

	AfterAll(func() {
		Expect(os.Remove(filename)).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		c.Del(context.Background(), storageClient.CacheKey(uploadedObject))
	})

	It("can return empty string if object is empty", func() {
		publicURL := storageClient.GetPublicURL(storage.PathUploadable{})
		Expect(publicURL).To(Equal(""))
	})

	It("can format a public url", func() {
		publicURL := storageClient.GetPublicURL(object)
		Expect(publicURL).To(Equal(fmt.Sprintf("https://storage.googleapis.com/%s/testdata/test-file.txt", bucketName)))
	})

	It("can sanitize a filename", func() {
		Expect(storageClient.SanitizeFilename(filename, time.Now().Unix())).To(MatchRegexp(`^blob_\d+.txt$`))
		Expect(storageClient.SanitizeFilename("model/thumbnails/", time.Now().Unix())).To(MatchRegexp(`^model/thumbnails/blob_\d+$`))
		Expect(storageClient.SanitizeFilename("model/thumbnails/file.txt", time.Now().Unix())).To(MatchRegexp(`^model/thumbnails/blob_\d+.txt$`))
	})

	It("can use the cache to get the signed url if it exists", func() {
		cachedURL := "https://storage.googleapis.com/pixo-test-bucket/testdata/cached-blob.txt"
		c.Set(context.Background(), storageClient.CacheKey(uploadedObject), cachedURL, 0)

		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(Equal(cachedURL))
	})

	It("can upload a file", func() {
		fileReader, err := os.Open(filename)
		Expect(err).NotTo(HaveOccurred())

		locationInBucket, err := storageClient.UploadFile(ctx, object, fileReader)

		Expect(err).NotTo(HaveOccurred())
		Expect(locationInBucket).To(Equal(fmt.Sprintf("testdata/blob_%d.txt", object.Timestamp)))
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
		locationInBucket, err := storageClient.UploadFile(ctx, object, bytes.NewReader([]byte("Go Blue")))
		Expect(err).NotTo(HaveOccurred())
		uploadedObject = storage.PathUploadable{
			BucketName: bucketName,
			Filepath:   locationInBucket,
		}

		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(googleAlgorithm))
		Expect(signedURL).To(ContainSubstring(googleExpiresAlgorithmParam))
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

		cacheRes := c.Get(context.Background(), fmt.Sprintf("signed-url:%s/%s", bucketName, uploadedObject.Filepath))
		Expect(cacheRes.Err()).NotTo(HaveOccurred())
		Expect(cacheRes.Val()).To(Equal(signedURL))
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
		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject, blobstorage.Option{
			ContentDisposition: "attachment",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(googleAlgorithm))
		Expect(signedURL).To(ContainSubstring(googleExpiresAlgorithmParam))
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring("attachment"))
	})

	It("can set the expiration for the signed url", func() {
		lifetime := 5 * time.Hour

		signedURL, err := storageClient.GetSignedURL(ctx, uploadedObject, blobstorage.Option{Lifetime: lifetime})

		Expect(err).NotTo(HaveOccurred())
		Expect(signedURL).To(ContainSubstring(bucketName))
		Expect(signedURL).To(ContainSubstring(googleAlgorithm))
		Expect(signedURL).To(ContainSubstring(googleDateParam))
		Expect(signedURL).To(ContainSubstring(googleExpiresAlgorithmParam))
		Expect(signedURL).To(ContainLifetime(lifetime))
	})

})

func ContainLifetime(lifetime time.Duration) types.GomegaMatcher {
	return &containsLifetimeMatcher{
		expected: lifetime,
	}
}

type containsLifetimeMatcher struct {
	expected time.Duration
}

// checks for expected.Seconds() +- 1 second
func (m *containsLifetimeMatcher) Match(signedURL interface{}) (success bool, err error) {
	signedURLStr := signedURL.(string)
	expectedStr := fmt.Sprintf("%d", int(m.expected.Seconds()))
	expectedMinusOneStr := fmt.Sprintf("%d", int(m.expected.Seconds()-1))
	expectedPlusOneStr := fmt.Sprintf("%d", int(m.expected.Seconds()+1))

	if success, err = ContainSubstring(expectedStr).Match(signedURLStr); success {
		return true, nil
	}

	if success, err = ContainSubstring(expectedMinusOneStr).Match(signedURLStr); success {
		return true, nil
	}

	return ContainSubstring(expectedPlusOneStr).Match(signedURLStr)
}

func (m *containsLifetimeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected %s to contain %s", actual, m.expected)
}

func (m *containsLifetimeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected %s not to contain %s", actual, m.expected)
}
