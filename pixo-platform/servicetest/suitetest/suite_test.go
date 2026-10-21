package suitetest_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server test suite", func() {
	It("fails a strict run when a step matches no definition", func() {
		status := run(`Feature: undefined step
  Scenario: typo
    Given I create a random integr
`)

		Expect(status).NotTo(Equal(0))
	})

	It("passes a strict run that uses the item-at-index not-null step", func() {
		status := run(`Feature: item not null
  Scenario: item has the property
    Given the response is '[{"id": 1}]'
    Then the response should contain an item at index 0 with "id" that is not null
`)

		Expect(status).To(Equal(0))
	})

	It("runs a repo step instead of the library step with the same expression", func() {
		overrideCalled = false

		status := run(`Feature: override
  Scenario: repo step wins
    Given I create a random integer
`)

		Expect(status).To(Equal(0))
		Expect(overrideCalled).To(BeTrue())
	})
})
