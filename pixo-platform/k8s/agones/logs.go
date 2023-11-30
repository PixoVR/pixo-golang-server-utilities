package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	log "github.com/rs/zerolog/log"
)

func (c Client) GetGameServerLogs(ctx context.Context, baseClient *base.Client, namespace string, gs agonesv1.GameServer) (string, string, error) {
	if baseClient == nil {
		return "", "", errors.New("baseClient is nil")
	}

	podName := gs.ObjectMeta.Name

	logs, err := baseClient.GetPodLogs(ctx, namespace, podName, DefaultGameServerContainerName, false)
	if err != nil {
		log.Err(err).Msgf("error getting pod logs for pod %v", podName)
		return "", "", err
	}

	sidecarLogs, err := baseClient.GetPodLogs(ctx, namespace, podName, DefaultGameServerSidecarContainerName, false)
	if err != nil {
		log.Err(err).Msgf("error getting pod logs for pod %v", podName)
		return "", "", err
	}

	return base.ReadLogsFromStream(logs), base.ReadLogsFromStream(sidecarLogs), nil
}
