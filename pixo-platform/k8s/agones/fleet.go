package agones

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) GetFleetsBySelectors(namespace string, labelSelectors *metav1.LabelSelector) (*v1.FleetList, error) {
	log.Debug().Msg("Fetching fleets")

	var options metav1.ListOptions

	if labelSelectors == nil {
		options = metav1.ListOptions{}
	} else {
		options = metav1.ListOptions{LabelSelector: labelSelectors.String()}
	}

	fleets, err := c.Clientset.AgonesV1().Fleets(namespace).List(context.Background(), options)

	if err != nil {
		log.Error().Err(err).Msg("Error fetching label information")
		return nil, err
	}

	return fleets, err
}

func (c Client) GetFleet(namespace, name string) (*v1.Fleet, error) {
	log.Debug().Msgf("Fetching fleet %s", name)

	gameserver, err := c.Clientset.
		AgonesV1().
		Fleets(namespace).
		Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get fleet %s", name)
		return nil, err
	}

	return gameserver, err
}

func (c Client) CreateFleet(ctx context.Context, namespace string, fleet *v1.Fleet) (*v1.Fleet, error) {
	log.Debug().Msg("Creating fleet")

	newFleet, err := c.Clientset.
		AgonesV1().
		Fleets(namespace).
		Create(ctx, fleet, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to create a fleet")
		return nil, err
	}

	return newFleet, err
}

func (c Client) DeleteFleet(ctx context.Context, namespace, name string) error {
	log.Debug().Msgf("Deleting fleet %s", name)

	err := c.Clientset.
		AgonesV1().
		Fleets(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete fleet %s", name)
		return err
	}

	return nil
}
