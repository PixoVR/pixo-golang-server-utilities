package base

import (
	"bytes"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	corev1 "k8s.io/api/core/v1"
)

func GetPodLogs(podName string, containerName string) (string, error) {

	podLogOpts := corev1.PodLogOptions{
		Container: containerName,
		TailLines: &[]int64{100}[0],
	}

	req := K8sClient.CoreV1().Pods(Namespace).GetLogs(podName, &podLogOpts)

	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Err(err).Msgf("Error in opening stream: %s", err)
		return "", err
	}

	defer func(podLogs io.ReadCloser) {
		err = podLogs.Close()
		if err != nil {
			log.Err(err).Msgf("Error in closing stream: %s", err)
		}
	}(podLogs)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	str := buf.String()

	return str, err
}
