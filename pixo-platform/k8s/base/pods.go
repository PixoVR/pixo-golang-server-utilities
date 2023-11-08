package base

import (
	"context"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) GetPods(ctx context.Context, namespace string) (*v1.PodList, error) {
	log.Debug().Msg("Fetching pods")

	pods, err := c.Clientset.
		CoreV1().
		Pods(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		log.Error().Err(err).Msg("Error fetching pods")
		return nil, err
	}

	return pods, err
}

func (c Client) GetPod(ctx context.Context, namespace, name string) (*v1.Pod, error) {
	log.Debug().Msgf("Fetching pod %s in namespace %s", name, namespace)
	log.Debug().Msgf("Clientset: %v", c.Clientset)

	pod, err := c.Clientset.
		CoreV1().
		Pods(namespace).
		Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get pod %s", name)
		return nil, err
	}

	return pod, err
}
