package blobstorage_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBlobStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BlobStorage Suite")
}
