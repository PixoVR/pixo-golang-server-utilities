package agones_test

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Fleets", func() {

	It("can get the list of fleets", func() {
		fleets, err := agonesClient.GetFleetsBySelectors(namespace, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(fleets).NotTo(BeNil())
	})

	It("can create, get, and delete a fleet", func() {

		sampleFleetObject := &agonesv1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-fleet",
				Namespace: namespace,
			},
			Spec: agonesv1.FleetSpec{
				Replicas: 1,
				Template: agonesv1.GameServerTemplateSpec{
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
				},
			},
		}

		fleet, err := agonesClient.CreateFleet(namespace, sampleFleetObject)

		Expect(err).NotTo(HaveOccurred())
		Expect(fleet).NotTo(BeNil())

		newGameserver, err := agonesClient.GetFleet(namespace, fleet.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(newGameserver).NotTo(BeNil())

		err = agonesClient.DeleteFleet(namespace, fleet.Name)
		Expect(err).NotTo(HaveOccurred())
	})

})
