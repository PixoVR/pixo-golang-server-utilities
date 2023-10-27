package agones_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fleets", Ordered, func() {

	It("can get the list of fleets", func() {
		fleets, err := agonesClient.GetFleetsBySelectors(namespace, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(fleets).NotTo(BeNil())
	})

	It("can get a fleet by name", func() {
		newGameserver, err := agonesClient.GetFleet(namespace, fleetName)
		Expect(err).NotTo(HaveOccurred())
		Expect(newGameserver).NotTo(BeNil())
	})

	It("can delete a fleet", func() {
		err := agonesClient.DeleteFleet(namespace, fleetName)
		Expect(err).NotTo(HaveOccurred())
	})

})
