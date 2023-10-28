package agones_test

import (
	autoscaling "agones.dev/agones/pkg/apis/autoscaling/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("Fleet Autoscalers", Ordered, func() {

	It("can create a fleet autoscaler", func() {
		autoscalerObject := autoscaling.FleetAutoscaler{
			ObjectMeta: v1.ObjectMeta{
				Name:      fleetName,
				Namespace: namespace,
			},
			Spec: autoscaling.FleetAutoscalerSpec{
				FleetName: fleetName,
				Policy: autoscaling.FleetAutoscalerPolicy{
					Type: autoscaling.BufferPolicyType,
					Buffer: &autoscaling.BufferPolicy{
						BufferSize: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: 3,
						},
						MinReplicas: 3,
						MaxReplicas: 10,
					},
				},
				Sync: &autoscaling.FleetAutoscalerSync{
					Type: autoscaling.FixedIntervalSyncType,
					FixedInterval: autoscaling.FixedIntervalSync{
						Seconds: 60,
					},
				},
			},
		}
		autoscaler, err := agonesClient.CreateFleetAutoscaler(namespace, &autoscalerObject)
		Expect(err).NotTo(HaveOccurred())
		Expect(autoscaler).NotTo(BeNil())
	})

	It("can get a fleet autoscaler", func() {
		autoscaler, err := agonesClient.GetFleetAutoscaler(namespace, fleetName)
		Expect(err).NotTo(HaveOccurred())
		Expect(autoscaler).NotTo(BeNil())
	})

	It("can delete a fleet autoscaler", func() {
		err := agonesClient.DeleteFleetAutoscaler(namespace, fleetName)
		Expect(err).NotTo(HaveOccurred())
	})

})
