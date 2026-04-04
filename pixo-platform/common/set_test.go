package common_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set", func() {
	var set *common.Set[string]

	BeforeEach(func() {
		set = common.NewSet[string]()
	})

	Context("basic operations", func() {
		It("should add and contain items", func() {
			set.Add("apple")
			Expect(set.Contains("apple")).To(BeTrue())
		})

		It("should maintain correct size", func() {
			set.Add("apple")
			Expect(set.Size()).To(Equal(1))

			// Adding duplicate
			set.Add("apple")
			set.Add("apple")
			Expect(set.Size()).To(Equal(1), "size should not increase when adding duplicate")
		})

		It("should remove items", func() {
			set.Add("apple")
			set.Remove("apple")
			Expect(set.Contains("apple")).To(BeFalse())
		})

		It("should clear all items", func() {
			set.Add("apple")
			set.Add("banana")
			set.Clear()
			Expect(set.Size()).To(BeZero())
		})
	})

	Context("ToSlice operation", func() {
		It("should convert set to slice correctly", func() {
			intSet := common.NewSet[int]()
			intSet.Add(1)
			intSet.Add(2)
			intSet.Add(3)

			slice := intSet.ToSlice()
			Expect(slice).To(HaveLen(3))
			Expect(slice).To(ContainElements(1, 2, 3))
		})
	})

	Context("edge cases", func() {
		It("should handle empty set operations", func() {
			Expect(set.Size()).To(BeZero())
			Expect(set.Contains("nonexistent")).To(BeFalse())
			set.Remove("nonexistent") // should not panic
			Expect(set.ToSlice()).To(BeEmpty())
		})
	})
})
