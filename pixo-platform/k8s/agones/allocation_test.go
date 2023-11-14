package agones_test

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var _ = Describe("Allocations", func() {

	var (
		gameserver *agonesv1.GameServer
	)

	BeforeEach(func() {
		var err error
		gameserver, err = agonesClient.CreateGameServer(context.Background(), namespace, &agones.SimpleGameServer)
		Expect(err).NotTo(HaveOccurred())
		Expect(gameserver).NotTo(BeNil())
	})

	AfterEach(func() {
		err := agonesClient.DeleteGameServer(context.Background(), namespace, gameserver.GetName())
		Expect(err).NotTo(HaveOccurred())
	})

	It("can create and allocate a game server", func() {
		sampleGameServerAllocation := &allocationv1.GameServerAllocation{
			Spec: allocationv1.GameServerAllocationSpec{
				Selectors: []allocationv1.GameServerSelector{
					{
						LabelSelector: metav1.LabelSelector{
							MatchLabels: labels.Set{
								"agones.dev/sdk-OrgID":    "1",
								"agones.dev/sdk-ModuleID": "1",
							},
						},
					},
				},
			},
		}

		allocation, err := agonesClient.CreateGameServerAllocation(context.Background(), namespace, sampleGameServerAllocation)

		Expect(err).NotTo(HaveOccurred())
		Expect(allocation).NotTo(BeNil())
	})

})
