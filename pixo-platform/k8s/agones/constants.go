package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	DefaultGameServerContainerName        = "gameserver"
	DefaultGameServerSidecarContainerName = "agones-gameserver-sidecar"
	DefaultGameServerPortName             = "udp"
	DefaultGameServerPort                 = 7777
	SimpleGameServerImage                 = "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
)

var (
	GameServerTolerations = []corev1.Toleration{
		{
			Effect:   "NoExecute",
			Key:      "gameserver",
			Operator: "Equal",
			Value:    "true",
		},
	}

	SimpleGameServer = agonesv1.GameServer{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-gameserver-",
			Labels: labels.Set{
				"agones.dev/sdk-OrgID":    "1",
				"agones.dev/sdk-ModuleID": "1",
			},
		},
		Spec: agonesv1.GameServerSpec{
			Container: "simplegameserver",
			Ports: []agonesv1.GameServerPort{{
				Name:          DefaultGameServerPortName,
				ContainerPort: DefaultGameServerPort,
				Protocol:      corev1.ProtocolUDP,
				PortPolicy:    agonesv1.Dynamic,
			}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Tolerations: GameServerTolerations,
					Containers: []corev1.Container{
						{
							Name:            DefaultGameServerContainerName,
							Image:           SimpleGameServerImage,
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
)
