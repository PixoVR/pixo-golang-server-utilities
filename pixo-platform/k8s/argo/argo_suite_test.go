package argo_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestArgo(t *testing.T) {
	RegisterFailHandler(Fail)
	_ = os.Setenv("IS_LOCAL", "true")
	RunSpecs(t, "Argo Client Suite")
}
