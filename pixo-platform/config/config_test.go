package config_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Config", func() {
	It("can determine the lifecycle of the current environment", func() {
		expectedLifecycle := "test"
		actualLifecycle := config.GetLifecycle()
		Expect(actualLifecycle).To(Equal(expectedLifecycle))
	})

	It("can determine the region of the current environment", func() {
		os.Setenv("REGION", "us-central1")
		expectedLifecycle := "us-central1"
		actualLifecycle := config.GetRegion()
		Expect(actualLifecycle).To(Equal(expectedLifecycle))
	})

	It("can determine the region of the current environment when in saudi", func() {
		os.Setenv("REGION", "me-central2")
		expectedLifecycle := "saudi"
		actualLifecycle := config.GetRegion()
		Expect(actualLifecycle).To(Equal(expectedLifecycle))
	})

	It("can load environment variables", func() {
		config.LoadEnvVars()
		expectedLifecycle := "test"
		Expect(config.GetLifecycle()).To(Equal(expectedLifecycle))
	})

	// WARNING: This test assumes you are running from the root of the pixo-golang-server-utilities repo
	It("can determine the project root directory", func() {
		root := config.GetProjectRoot("../..")
		Expect(root).NotTo(BeNil())
		Expect(root).To(ContainSubstring("pixo-golang-server-utilities"))
		Expect(root).NotTo(ContainSubstring("pixo-platform"))
		Expect(root).NotTo(ContainSubstring("config"))
	})

})
