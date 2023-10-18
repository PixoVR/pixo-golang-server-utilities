package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	log "github.com/rs/zerolog/log"
)

func (c Client) GetGameServerLogs(baseClient base.Client, namespace string, gs agonesv1.GameServer) (string, string, error) {

	podName := gs.ObjectMeta.Name

	logs, err := baseClient.GetPodLogs(namespace, podName, "gameserver")
	if err != nil {
		log.Err(err).Msgf("error getting pod logs for pod %v", podName)
		return "", "", err
	}

	sidecarLogs, err := baseClient.GetPodLogs(namespace, podName, "agones-gameserver-sidecar")
	if err != nil {
		log.Err(err).Msgf("error getting pod logs for pod %v", podName)
		return "", "", err
	}

	return logs, sidecarLogs, nil
}
