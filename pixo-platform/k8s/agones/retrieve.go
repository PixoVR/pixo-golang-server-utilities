package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (c Client) IsGameServerReady(namespace, name string) bool {
	gameserver, err := c.GetGameServerByName(namespace, name)

	if err != nil {
		log.Err(err).Msgf("Error getting game servers by name: %v", name)
		return false
	}

	return gameserver.Status.State == agonesv1.GameServerStateReady
}

func (c Client) GetGameServerByName(namespace, name string) (*agonesv1.GameServer, error) {
	gameserver, err := c.Clientset.AgonesV1().GameServers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting game servers by name: %v", name)
		return nil, err
	}

	return gameserver, nil
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
