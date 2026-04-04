package common_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"
)

var _ = Describe("Pointer", func() {
	It("return the pointer if it is not nil", func() {
		Expect(common.GetPointerOrDefault(common.GetPointer(1), 1)).To(Equal(common.GetPointer(1)))
	})

	It("return the default if the pointer is nil", func() {
		Expect(common.GetPointerOrDefault(nil, 1)).To(Equal(common.GetPointer(1)))
	})

	Context("GetValueOrDefault", func() {
		It("return the value if it is not nil", func() {
			Expect(common.GetValueOrDefault(common.GetPointer(1), 1)).To(Equal(1))
		})

		It("return the default if the value is nil", func() {
			Expect(common.GetValueOrDefault(nil, 1)).To(Equal(1))
		})
	})

	Context("PointerValuesAreEqual", func() {
		DescribeTable("returns true if pointers are equal", func(a, b *int, isEqual bool) {
			Expect(common.PointerValuesAreEqual(a, b)).To(Equal(isEqual))
		},
			Entry("both nil", nil, nil, true),
			Entry("both non-nil and equal", common.GetPointer(1), common.GetPointer(1), true),
			Entry("both non-nil and not equal", common.GetPointer(1), common.GetPointer(2), false),
			Entry("one nil and one non-nil", nil, common.GetPointer(1), false),
			Entry("one non-nil and one nil", common.GetPointer(1), nil, false),
		)
	})
})
