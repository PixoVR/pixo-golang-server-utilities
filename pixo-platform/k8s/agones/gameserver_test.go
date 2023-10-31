package agones_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gameservers", func() {

	It("can get the list of gameservers", func() {
		gameservers, err := agonesClient.GetGameServers(namespace, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(gameservers).NotTo(BeNil())
	})

	It("can create, get, and delete a game server", func() {
		gameserver, err := agonesClient.CreateGameServer(namespace, &agones.SimpleGameServer)

		Expect(err).NotTo(HaveOccurred())
		Expect(gameserver).NotTo(BeNil())

		newGameserver, err := agonesClient.GetGameServer(namespace, gameserver.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(newGameserver).NotTo(BeNil())

		err = agonesClient.DeleteGameServer(namespace, gameserver.Name)
		Expect(err).NotTo(HaveOccurred())

		Expect(agonesClient.IsGameServerAvailable(namespace, gameserver.Name)).To(BeFalse())
	})

})
