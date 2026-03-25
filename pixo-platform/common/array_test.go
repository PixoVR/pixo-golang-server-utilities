package common_test

import (
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Array", func() {

	Context("Is In Slice", func() {
		It("can check if an integer value is in a slice of integers", func() {
			slice := []int{1, 2, 3, 4, 5}
			Expect(IsInSlice(slice, 3)).To(BeTrue())
			Expect(IsInSlice(slice, 6)).To(BeFalse())
		})

		It("can check if a string value is in a slice of strings", func() {
			slice := []string{"one", "two", "three", "four", "five"}
			Expect(IsInSlice(slice, "three")).To(BeTrue())
			Expect(IsInSlice(slice, "six")).To(BeFalse())
		})

		It("can check if a float value is in a slice of floats", func() {
			slice := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
			Expect(IsInSlice(slice, 3.3)).To(BeTrue())
			Expect(IsInSlice(slice, 6.6)).To(BeFalse())
		})

		It("can check if a boolean value is in a slice of booleans", func() {
			slice := []bool{true, false, true, false, true}
			Expect(IsInSlice(slice, true)).To(BeTrue())
			Expect(IsInSlice(slice, false)).To(BeTrue())
			Expect(IsInSlice(slice, nil)).To(BeFalse())
		})
	})

	Context("Remove Nil Pointers", func() {
		It("can remove nil pointers from a slice", func() {
			slice := []*int{nil, GetPointer(1), nil, GetPointer(2), nil}
			newSlice := RemoveNilPointers(slice)
			Expect(newSlice).To(Equal([]int{1, 2}))
		})

		It("can return nil if the slice is nil", func() {
			newSlice := RemoveNilPointers(nil)
			Expect(newSlice).To(BeNil())
		})

		It("can return nil if the param is not a slice", func() {
			newSlice := RemoveNilPointers(1)
			Expect(newSlice).To(BeNil())
		})
	})

	Context("Remove Nil Values", func() {

		Context("when given a nil slice", func() {
			It("returns nil", func() {
				var input []DummyStruct = nil
				Expect(RemoveNilValues(input)).To(BeNil())
			})
		})

		Context("when given an empty slice", func() {
			It("returns an empty slice", func() {
				input := make([]DummyStruct, 0)
				Expect(RemoveNilValues(input)).To(BeEmpty())
			})
		})

		Context("when given a slice with all non-nil values", func() {
			It("returns the same slice", func() {
				d1 := DummyStruct{}
				d2 := DummyStruct{}
				input := []DummyStruct{d1, d2}
				result := RemoveNilValues(input)
				Expect(result).To(Equal([]DummyStruct{d1, d2}))
			})
		})

		Context("when given a slice with some nil values", func() {
			It("returns a slice with only non-nil values", func() {
				d1 := DummyStruct{Value: 1}
				d2 := DummyStruct{Value: 2}
				input := []Dummy{nil, d1, nil, d2, nil}
				result := RemoveNilValues(input)
				Expect(result).To(Equal([]Dummy{d1, d2}))
			})
		})

		Context("when given a slice with all nil values", func() {
			It("returns an empty slice", func() {
				input := []Dummy{nil, nil}
				Expect(RemoveNilValues(input)).To(BeEmpty())
			})
		})
	})

	Context("ToPointerArray", func() {
		It("can convert a value array of type T to a pointer array of type T", func() {
			values := []int{1, 2, 3, 4, 5}
			pointers := ToPointerArray(values)
			Expect(pointers).To(Equal([]*int{GetPointer(1), GetPointer(2), GetPointer(3), GetPointer(4), GetPointer(5)}))
		})

		It("can return nil if the values array is nil", func() {
			pointers := ToPointerArray(nil)
			Expect(pointers).To(BeNil())
		})

		It("can return nil if the param is not a slice", func() {
			pointers := ToPointerArray(1)
			Expect(pointers).To(BeNil())
		})
	})

	Context("ToStringArray", func() {
		It("can convert a value array of type T to a string array", func() {
			values := []int{1, 2, 3, 4, 5}
			strings := ToStringArray(values)
			Expect(strings).To(Equal([]string{"1", "2", "3", "4", "5"}))
		})
	})

	Context("Get distinct values from slice", func() {
		type TestObject struct {
			StringProp string
			IntProp    int
		}

		var testSlice []TestObject
		BeforeEach(func() {
			testSlice = []TestObject{
				{"String1", 30},
				{"String2", 25},
				{"String1", 28},
				{"String3", 30},
			}
		})

		It("should return distinct names from a slice of Person structs", func() {
			distinctStrings := GetDistinctValues(testSlice, func(t TestObject) string {
				return t.StringProp
			})

			expectedStrings := []string{"String1", "String2", "String3"}
			Expect(distinctStrings).To(ConsistOf(expectedStrings))
		})

		It("should return distinct ages from a slice of Person structs", func() {
			distinctIntegers := GetDistinctValues(testSlice, func(t TestObject) int {
				return t.IntProp
			})

			expectedIntegers := []int{30, 25, 28}
			Expect(distinctIntegers).To(ConsistOf(expectedIntegers))
		})

		It("should handle an empty slice gracefully", func() {
			distinctStrings := GetDistinctValues([]TestObject{}, func(t TestObject) string {
				return t.StringProp
			})

			Expect(distinctStrings).To(BeEmpty())
		})

		It("should return distinct values for a slice of integers", func() {
			numbers := []int{1, 2, 3, 1, 2, 4}

			distinctNumbers := GetDistinctValues(numbers, func(n int) int {
				return n
			})

			expectedNumbers := []int{1, 2, 3, 4}
			Expect(distinctNumbers).To(ConsistOf(expectedNumbers))
		})
	})
})
