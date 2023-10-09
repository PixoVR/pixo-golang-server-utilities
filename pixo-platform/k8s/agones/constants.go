package agones

import corev1 "k8s.io/api/core/v1"

const (
	DefaultGameServerContainerName = "gameserver"
	DefaultGameServerPortName      = "udp"
	DefaultGameServerPort          = 7777
	SimpleGameServerImage          = "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
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
)
