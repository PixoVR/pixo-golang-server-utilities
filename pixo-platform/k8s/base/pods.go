package base

import (
	"context"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) GetPods(namespace string) (*v1.PodList, error) {
	log.Debug().Msg("Fetching pods")

	pods, err := c.Clientset.
		CoreV1().
		Pods(namespace).
		List(context.Background(), metav1.ListOptions{})

	if err != nil {
		log.Error().Err(err).Msg("Error fetching pods")
		return nil, err
	}

	return pods, err
}

func (c Client) GetPod(namespace, name string) (*v1.Pod, error) {
	log.Debug().Msgf("Fetching pod %s", name)

	pod, err := c.Clientset.
		CoreV1().
		Pods(namespace).
		Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get pod %s", name)
		return nil, err
	}

	return pod, err
}
