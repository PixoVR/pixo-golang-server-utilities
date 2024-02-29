package blobstorage_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/aws"
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

	It("can parse the filename from a file location", func() {
		filename := blobstorage.GetFilenameFromLocation("images/modules/distributor/file.png")
		Expect(filename).To(Equal("file.png"))
	})

	It("can parse the path for an image field object", func() {
		pathUploadable := aws.DefaultPublicUploadable{
			Path: "https://bucket.s3.amazonaws.com/images/modules/distributor/image.png",
		}
		Expect(pathUploadable.GetFileLocation()).To(Equal("images/modules/distributor/image.png"))
	})

})
