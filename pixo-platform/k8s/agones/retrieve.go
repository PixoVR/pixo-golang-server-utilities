package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

func (c Client) IsGameServerReady(namespace, name string) bool {
	gameservers, err := c.Clientset.AgonesV1().GameServers(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Err(err).Msgf("Error getting game servers by name: %v", name)
		return false
	}

	for _, gs := range gameservers.Items {
		isReady := gs.Status.State == agonesv1.GameServerStateReady
		if gs.GetName() == name {
			log.Info().Msgf("Game Server: (%s) %v, %v", name, gs.Name, gs.Status.State)
			if isReady {
				return true
			}
		}
	}

	return false
}

func (c Client) GetGameServerByName(namespace, name string) (*agonesv1.GameServer, error) {
	gameservers, err := c.Clientset.AgonesV1().GameServers(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting game servers by name: %v", name)
		return nil, err
	}

	for _, gs := range gameservers.Items {
		isReadyOrAllocated := gs.Status.State == agonesv1.GameServerStateReady ||
			gs.Status.State == agonesv1.GameServerStateAllocated

		if strings.Contains(gs.Name, name) && isReadyOrAllocated {
			return &gs, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("game server %s not found", name))
}

func (c Client) GetGameServersBySelectors(namespace string, selectors labels.Set) ([]agonesv1.GameServer, error) {
	gameservers, err := c.Clientset.AgonesV1().GameServers(namespace).List(context.TODO(),
		metav1.ListOptions{LabelSelector: selectors.String()})
	if err != nil {
		log.Err(err).Msgf("Unable to get gameservers with labels: %v", selectors)
		return nil, err
	}

	return gameservers.Items, nil
}
