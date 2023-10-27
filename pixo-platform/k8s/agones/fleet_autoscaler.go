package agones

import (
	autoscaling "agones.dev/agones/pkg/apis/autoscaling/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) GetFleetAutoscaler(namespace string, name string) (*autoscaling.FleetAutoscaler, error) {
	log.Debug().Msgf("Getting fleet autoscaler: %s", name)

	res, err := c.Clientset.
		AutoscalingV1().
		FleetAutoscalers(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get fleet autoscaler: %s", name)
		return nil, err
	}

	log.Debug().Msgf("Got fleet autoscaler: %s", name)
	return res, err
}

func (c Client) CreateFleetAutoscaler(namespace string, autoscaler *autoscaling.FleetAutoscaler) (*autoscaling.FleetAutoscaler, error) {
	log.Debug().Msgf("Creating fleet autoscaler for fleet: %s", autoscaler.Spec.FleetName)

	res, err := c.Clientset.
		AutoscalingV1().
		FleetAutoscalers(namespace).
		Create(context.TODO(), autoscaler, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to create fleet autoscaler for fleet: %s", autoscaler.Spec.FleetName)
		return nil, err
	}

	log.Debug().Msgf("Created fleet autoscaler for fleet: %s", autoscaler.Spec.FleetName)
	return res, err
}

func (c Client) DeleteFleetAutoscaler(namespace string, name string) error {
	log.Debug().Msgf("Deleting fleet autoscaler: %s", name)

	err := c.Clientset.
		AutoscalingV1().
		FleetAutoscalers(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete fleet autoscaler: %s", name)
		return err
	}

	log.Debug().Msgf("Deleted fleet autoscaler: %s", name)
	return nil
}
