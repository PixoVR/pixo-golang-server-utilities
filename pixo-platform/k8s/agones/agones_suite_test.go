package agones_test

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAgones(t *testing.T) {
	RegisterFailHandler(Fail)
	_ = os.Setenv("IS_LOCAL", "true")
	RunSpecs(t, "Agones Client Suite")
}

var (
	namespace    = "dev-multiplayer"
	fleetName    = "test-fleet"
	agonesClient *agones.Client
)

var _ = BeforeSuite(func() {
	baseClient, err := base.NewLocalClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(baseClient).To(Not(BeNil()))

	agonesClient, err = agones.NewLocalAgonesClient(*baseClient)
	Expect(err).NotTo(HaveOccurred())
	Expect(agonesClient).To(Not(BeNil()))

	sampleFleetObject := &agonesv1.Fleet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fleetName,
			Namespace: namespace,
		},
		Spec: agonesv1.FleetSpec{
			Replicas: 1,
			Template: agonesv1.GameServerTemplateSpec{
				Spec: agonesv1.GameServerSpec{
					Container: agones.DefaultGameServerContainerName,
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
})

var _ = AfterSuite(func() {
	err := agonesClient.DeleteFleet(namespace, fleetName)
	Expect(err).NotTo(HaveOccurred())
})
