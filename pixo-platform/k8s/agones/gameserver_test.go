package agones_test

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gameservers", Ordered, func() {

	var (
		gameserver *v1.GameServer
		ctx        = context.Background()
	)

	BeforeAll(func() {
		var err error
		gameserver, err = agonesClient.CreateGameServer(ctx, namespace, &agones.SimpleGameServer)
		Expect(err).NotTo(HaveOccurred())
		Expect(gameserver).NotTo(BeNil())
		Expect(gameserver.Labels[agones.DeletedGameServerLabel]).To(Equal("false"))
	})

	AfterAll(func() {
		err := agonesClient.DeleteGameServer(ctx, namespace, gameserver.Name)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can get the list of gameservers", func() {
		gameservers, err := agonesClient.GetGameServers(ctx, namespace, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(gameservers).NotTo(BeNil())
		Expect(len(gameservers.Items)).To(BeNumerically(">", 0))
	})

	It("can get a gameserver and add a label to it", func() {
		Expect(gameserver.Labels["test"]).To(BeEmpty())

		updatedGameserver, err := agonesClient.AddLabelToGameServer(ctx, gameserver, "test", "test")

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedGameserver).NotTo(BeNil())
		Expect(updatedGameserver.Labels["test"]).To(Equal("test"))
	})

	It("can delete a game server and then tell its unavailable", func() {
		isAvailable := agonesClient.IsGameServerAvailable(ctx, namespace, gameserver.GetName())
		Expect(isAvailable).To(BeTrue())

		newGameserver, err := agonesClient.GetGameServer(ctx, namespace, gameserver.GetName())
		Expect(err).NotTo(HaveOccurred())
		Expect(newGameserver).NotTo(BeNil())

		Expect(agonesClient.DeleteGameServer(ctx, namespace, gameserver.GetName())).To(Succeed())

		isAvailable = agonesClient.IsGameServerAvailable(ctx, namespace, gameserver.GetName())
		Expect(isAvailable).To(BeFalse())

		gameservers, err := agonesClient.GetGameServers(ctx, namespace, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(gameservers).NotTo(BeNil())

		foundGameserver := false
		for _, gs := range gameservers.Items {
			if gs.GetName() == gameserver.GetName() {
				foundGameserver = true
				break
			}
		}
		Expect(foundGameserver).To(BeFalse())
	})

})
