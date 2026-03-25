package common_test

import (
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Strings", func() {
	It("can determine if a string is nil or empty", func() {
		Expect(IsNilOrEmpty(nil)).To(BeTrue())
		Expect(IsNilOrEmpty(GetPointer(""))).To(BeTrue())
		Expect(IsNilOrEmpty(GetPointer(" "))).To(BeTrue())
		Expect(IsNilOrEmpty(GetPointer("test"))).To(BeFalse())
	})

	It("can determine if a string is empty", func() {
		Expect(IsEmptyString("")).To(BeTrue())
		Expect(IsEmptyString(" ")).To(BeTrue())
		Expect(IsEmptyString("test")).To(BeFalse())
	})
})
