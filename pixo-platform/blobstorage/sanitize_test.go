package blobstorage_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sanitize", func() {

	It("can handle degenerate cases", func() {
		Expect(blobstorage.ParseFileLocationFromLink("")).To(Equal(""))
		Expect(blobstorage.ParseFileLocationFromLink("null")).To(Equal("null"))
	})

	It("can parse the path from the link", func() {
		link := "https://bucket.s3.amazonaws.com/images/modules/distributor/file.png"
		expectedFileLocation := "images/modules/distributor/file.png"
		actualFileLocation := blobstorage.ParseFileLocationFromLink(link)

		Expect(actualFileLocation).To(Equal(expectedFileLocation))

		By("being idempotent", func() {
			Expect(blobstorage.ParseFileLocationFromLink(actualFileLocation)).To(Equal(expectedFileLocation))
		})
	})

	It("can parse the path from the stc link", func() {
		link := "https://bucketname.api-object.bluvalt.com:8082/org/module/file.zip"
		expectedFileLocation := "org/module/file.zip"
		actualFileLocation := blobstorage.ParseFileLocationFromLink(link)

		Expect(actualFileLocation).To(Equal(expectedFileLocation))

		By("being idempotent", func() {
			Expect(blobstorage.ParseFileLocationFromLink(actualFileLocation)).To(Equal(expectedFileLocation))
		})
	})

	It("can ignore the timestamp if its 0", func() {
		link := "images/modules/distributor/file.png"
		expectedFileLocation := "images/modules/distributor/blob.png"
		actualFileLocation := blobstorage.SanitizeFilename(link, 0)

		Expect(actualFileLocation).To(Equal(expectedFileLocation))

		By("being idempotent", func() {
			Expect(blobstorage.ParseFileLocationFromLink(actualFileLocation)).To(Equal(expectedFileLocation))
		})
	})

	It("can parse the filename from a file location", func() {
		filename := blobstorage.GetFilenameFromLocation("images/modules/distributor/file.png")
		Expect(filename).To(Equal("file.png"))
	})

})
