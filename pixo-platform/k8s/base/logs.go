package base

import (
	"bytes"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	corev1 "k8s.io/api/core/v1"
)

func (c Client) GetPodLogs(ctx context.Context, namespace, podName, containerName string, follow bool) (io.ReadCloser, error) {

	podLogOpts := corev1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
	}

	req := c.
		CoreV1().
		Pods(namespace).
		GetLogs(podName, &podLogOpts)

	return req.Stream(ctx)
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
