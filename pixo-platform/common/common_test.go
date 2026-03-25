package common_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("General", func() {
	Describe("Filter", func() {
		It("returns only the filtered elements", func() {
			numberList := []int{1, 2, 3, 4, 5}
			filterFunc := func(i int) bool {
				return i%2 == 0
			}

			Expect(common.Filter(numberList, filterFunc)).To(Equal([]int{2, 4}))
		})

		It("returns only the first n elements", func() {
			numberList := []int{1, 2, 3, 4, 5}

			Expect(common.Take(numberList, 3)).To(Equal([]int{1, 2, 3}))
			Expect(common.Take(numberList, 999)).To(Equal([]int{1, 2, 3, 4, 5}))
			Expect(common.Take(numberList, 0)).To(Equal([]int{}))
			Expect(common.Take(numberList, -1)).To(Equal([]int{}))

			Expect(common.Take[int](nil, 6)).To(Equal([]int{}))
		})

		It("finds the first element", func() {
			numberList := []int{1, 2, 3, 4, 5}
			findfunc := func(i int) bool { return i%2 == 0 }

			result := common.Find(numberList, findfunc)
			expectedResult := 2

			Expect(result).To(Equal(&expectedResult))
		})

		It("plucks the key", func() {
			numberList := []int{1, 2, 3, 4, 5}
			keyFunc := func(i int) int { return i * 2 }

			Expect(common.Pluck(numberList, keyFunc)).To(Equal([]int{2, 4, 6, 8, 10}))
		})

		It("calculates the average", func() {
			numberList := []int{1, 2, 3, 4, 5}
			keyFunc := func(i int) float64 { return float64(i) }

			Expect(common.Average(numberList, keyFunc)).To(Equal(3.0))
		})
	})
	Describe("Contains", func() {
		It("returns true if the item is in the list", func() {
			numberList := []int{1, 2, 3, 4, 5}
			Expect(common.Contains(numberList, 3)).To(BeTrue())
		})
		It("returns false if the item is not in the list", func() {
			numberList := []int{1, 2, 3, 4, 5}
			Expect(common.Contains(numberList, 6)).To(BeFalse())
		})
	})

	Describe("Get Value from object", func() {
		DescribeTable("returns the value of the key for a struct or pointer to a struct",
			func(obj interface{}, fieldName string, expectedVal *string, expectErr bool, errMsg *string) {
				value, err := common.GetFieldValue(obj, fieldName)

				if expectErr {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(*errMsg))
					Expect(value).To(BeNil())
				} else {
					Expect(err).NotTo(HaveOccurred())
					Expect(value).To(Equal(*expectedVal))
				}
			},
			Entry("returns an error if the object is not a struct", 123, "key", nil, true, common.GetPointer("expected a struct or a pointer to a struct")),
			Entry("return an error if the object is a pointer but not to a struct", common.GetPointer(123), "key", nil, true, common.GetPointer("expected a struct or a pointer to a struct")),
			Entry("returns the value of the key for a struct", struct{ Value string }{Value: "value"}, "Value", common.GetPointer("value"), false, nil),
			Entry("returns the value of the key for a pointer to a struct", &struct{ Value string }{Value: "value"}, "Value", common.GetPointer("value"), false, nil),
			Entry("returns nil if the key is not found", struct{ Value string }{Value: "value"}, "not-key", nil, true, common.GetPointer("no such field")),
		)
	})
	Describe("Values", func() {
		Context("with a non-empty map", func() {
			It("returns all the values for a map with string values", func() {
				m := map[int]string{
					1: "one",
					2: "two",
					3: "three",
				}
				values := common.Values(m)
				Expect(values).To(ConsistOf("one", "two", "three"))
			})

			It("returns all the values for a map with bool values", func() {
				m := map[string]bool{
					"a": true,
					"b": false,
					"c": true,
				}
				values := common.Values(m)
				Expect(values).To(ConsistOf(true, false, true))
			})
		})

		Context("with an empty map", func() {
			It("returns an empty slice", func() {
				m := map[string]int{}
				values := common.Values(m)
				Expect(values).To(BeEmpty())
			})
		})
	})
})
