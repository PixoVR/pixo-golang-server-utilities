package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"bytes"
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	log "github.com/rs/zerolog/log"
	"io"
)

func (c Client) GetGameServerLogs(ctx context.Context, baseClient *base.Client, namespace string, gs agonesv1.GameServer) (string, string, error) {
	if baseClient == nil {
		return "", "", errors.New("baseClient is nil")
	}

	podName := gs.ObjectMeta.Name

	logs, err := baseClient.GetPodLogs(ctx, namespace, podName, DefaultGameServerContainerName)
	if err != nil {
		log.Err(err).Msgf("error getting pod logs for pod %v", podName)
		return "", "", err
	}

	sidecarLogs, err := baseClient.GetPodLogs(ctx, namespace, podName, DefaultGameServerSidecarContainerName)
	if err != nil {
		log.Err(err).Msgf("error getting pod logs for pod %v", podName)
		return "", "", err
	}

	return ReadLogsFromStream(logs), ReadLogsFromStream(sidecarLogs), nil
}

func ReadLogsFromStream(stream io.ReadCloser) string {

	defer func(podLogs io.ReadCloser) {
		if err := podLogs.Close(); err != nil {
			log.Err(err).Msg("Error in closing stream")
		}
	}(stream)

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, stream); err != nil {
		return "unable to read logs"
	}

	return buf.String()
}
