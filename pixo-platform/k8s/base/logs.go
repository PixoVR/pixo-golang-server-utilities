package base

import (
	"context"
	"io"
	corev1 "k8s.io/api/core/v1"
)

func (c Client) GetPodLogs(ctx context.Context, namespace, podName, containerName string) (io.ReadCloser, error) {

	podLogOpts := corev1.PodLogOptions{
		Container: containerName,
		//TailLines: &[]int64{100}[0],
	}

	req := c.
		CoreV1().
		Pods(namespace).
		GetLogs(podName, &podLogOpts)

	return req.Stream(ctx)
}
