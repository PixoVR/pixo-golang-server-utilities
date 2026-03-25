package common_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"
)

var _ = Describe("ToJSONString", func() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	It("should convert a struct to a JSON string", func() {
		input := Person{Name: "Alice", Age: 30}
		jsonStr, err := common.ToJSONString(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(jsonStr).To(MatchJSON(`{"name":"Alice","age":30}`))
	})

	It("should convert a map to a JSON string", func() {
		input := map[string]interface{}{
			"foo": "bar",
			"num": 42,
		}
		jsonStr, err := common.ToJSONString(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(jsonStr).To(MatchJSON(`{"foo":"bar", "num":42}`))
	})

	It("should return an error when marshalling an unsupported value", func() {
		input := func() {}

		_, err := common.ToJSONString(input)
		Expect(err).To(HaveOccurred())
	})

	It("should convert a slice of strings to JSON", func() {
		input := []string{"apple", "banana", "cherry"}
		jsonStr, err := common.ToJSONString(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(jsonStr).To(Equal(`["apple","banana","cherry"]`))
	})
})
