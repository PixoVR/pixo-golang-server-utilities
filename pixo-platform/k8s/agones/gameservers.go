package agones

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *Client) GetGameServers(labelSelectors *metav1.LabelSelector) (*v1.GameServerList, error) {
	log.Debug().Msg("Fetching game servers")

	var options metav1.ListOptions

	if labelSelectors == nil {
		options = metav1.ListOptions{}
	} else {
		options = metav1.ListOptions{LabelSelector: labelSelectors.String()}
	}

	gameservers, err := a.Clientset.AgonesV1().GameServers(a.Namespace).List(context.Background(), options)

	if err != nil {
		log.Error().Err(err).Msg("Error fetching label information")
		return nil, err
	}

	return gameservers, err
}

func (a *Client) GetGameServer(name string) (*v1.GameServer, error) {
	log.Debug().Msgf("Fetching game server %s", name)

	gameserver, err := a.Clientset.
		AgonesV1().
		GameServers(a.Namespace).
		Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get game server %s", name)
		return nil, err
	}

	return gameserver, err
}

func (a *Client) CreateGameServer(gameserver *v1.GameServer) (*v1.GameServer, error) {
	log.Debug().Msg("Creating game server")

	newGameServer, err := a.Clientset.
		AgonesV1().
		GameServers(a.Namespace).
		Create(context.TODO(), gameserver, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to create a game server")
		return nil, err
	}

	return newGameServer, err
}

func (a *Client) DeleteGameServer(name string) error {
	log.Debug().Msgf("Deleting game server %s", name)

	err := a.Clientset.
		AgonesV1().
		GameServers(a.Namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete game server %s", name)
		return err
	}

	return nil
}
