package agones_test

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Agones", func() {

	var (
		agonesClient *agones.Client
		namespace    = "dev-multiplayer"
	)

	BeforeEach(func() {
		var err error
		agonesClient, err = agones.NewLocalAgonesClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(agonesClient).To(Not(BeNil()))
	})

	It("can get the list of gameservers", func() {
		gameservers, err := agonesClient.GetGameServers(namespace, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(gameservers).NotTo(BeNil())
	})

	It("can create, get, and delete a game server", func() {
		sampleGameServer := &agonesv1.GameServer{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-gameserver-",
			},
			Spec: agonesv1.GameServerSpec{
				Container: "simplegameserver",
				Ports: []agonesv1.GameServerPort{{
					Name:          agones.DefaultGameServerPortName,
					ContainerPort: agones.DefaultGameServerPort,
					Protocol:      corev1.ProtocolUDP,
					PortPolicy:    agonesv1.Dynamic,
				}},
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Tolerations: agones.GameServerTolerations,
						Containers: []corev1.Container{
							{
								Name:            agones.DefaultGameServerContainerName,
								Image:           agones.SimpleGameServerImage,
								ImagePullPolicy: corev1.PullAlways,
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{},
								},
							},
						},
					},
				},
			},
		}

		gameserver, err := agonesClient.CreateGameServer(namespace, sampleGameServer)

		Expect(err).NotTo(HaveOccurred())
		Expect(gameserver).NotTo(BeNil())

		newGameserver, err := agonesClient.GetGameServer(namespace, gameserver.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(newGameserver).NotTo(BeNil())

		err = agonesClient.DeleteGameServer(namespace, gameserver.Name)
		Expect(err).NotTo(HaveOccurred())
	})

})
