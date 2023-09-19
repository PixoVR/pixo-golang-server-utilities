package auth_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Auth", func() {

	var (
		originalLifecycle string
	)

	BeforeEach(func() {
		originalLifecycle = os.Getenv("SECRET_KEY")
		os.Setenv("SECRET_KEY", "test")
	})

	AfterEach(func() {
		os.Setenv("SECRET_KEY", originalLifecycle)
	})

	It("can determine if a secret key is valid", func() {
		isSecretKey := auth.IsValidSecretKey("test")
		Expect(isSecretKey).To(BeTrue())
	})

	It("can determine if a secret key is invalid", func() {
		isNotSecretKey := auth.IsValidSecretKey("")
		Expect(isNotSecretKey).To(BeFalse())
	})
})
