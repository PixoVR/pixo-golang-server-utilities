package agones

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (c Client) GetGameServers(ctx context.Context, namespace string, labelSelectors labels.Set) (*agonesv1.GameServerList, error) {
	log.Debug().Msg("Fetching game servers")

	if labelSelectors == nil {
		labelSelectors = labels.Set{}
	}

	labelSelectors[DeletedGameServerLabel] = "false"

	options := metav1.ListOptions{LabelSelector: labelSelectors.String()}

	gameservers, err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		List(ctx, options)

	if err != nil {
		log.Error().Err(err).Msg("Error fetching game servers")
		return nil, err
	}

	return gameservers, err
}

func (c Client) GetGameServer(ctx context.Context, namespace, name string) (*agonesv1.GameServer, error) {
	log.Debug().Msgf("Fetching game server %s", name)

	gameserver, err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get game server %s", name)
		return nil, err
	}

	return gameserver, err
}

func (c Client) GetGameServerByName(ctx context.Context, namespace, name string) (*agonesv1.GameServer, error) {

	gameserver, err := c.
		AgonesV1().
		GameServers(namespace).
		Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		log.Err(err).Msgf("Error getting game servers by name: %v", name)
		return nil, err
	}

	return gameserver, nil
}

func (c Client) IsGameServerReady(ctx context.Context, namespace, name string) bool {
	gameserver, err := c.GetGameServerByName(ctx, namespace, name)

	if err != nil {
		log.Err(err).Msgf("Error getting game servers by name: %v", name)
		return false
	}

	return gameserver.Status.State == agonesv1.GameServerStateReady && c.IsGameServerAvailable(ctx, namespace, name)
}
