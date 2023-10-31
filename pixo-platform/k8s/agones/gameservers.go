package agones

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"github.com/rs/zerolog/log"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) GetGameServers(namespace string, labelSelectors *metav1.LabelSelector) (*v1.GameServerList, error) {
	log.Debug().Msg("Fetching game servers")

	var options metav1.ListOptions

	if labelSelectors == nil {
		options = metav1.ListOptions{}
	} else {
		options = metav1.ListOptions{LabelSelector: labelSelectors.String()}
	}

	gameservers, err := c.Clientset.AgonesV1().GameServers(namespace).List(context.Background(), options)

	if err != nil {
		log.Error().Err(err).Msg("Error fetching label information")
		return nil, err
	}

	return gameservers, err
}

func (c Client) GetGameServer(namespace, name string) (*v1.GameServer, error) {
	log.Debug().Msgf("Fetching game server %s", name)

	gameserver, err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get game server %s", name)
		return nil, err
	}

	return gameserver, err
}

func (c Client) CreateGameServer(namespace string, gameserver *v1.GameServer) (*v1.GameServer, error) {
	log.Debug().Msg("Creating game server")

	newGameServer, err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Create(context.TODO(), gameserver, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to create a game server")
		return nil, err
	}

	return newGameServer, err
}

func (c Client) DeleteGameServer(namespace, name string) error {
	log.Debug().Msgf("Deleting game server %s", name)

	err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete game server %s", name)
		return err
	}

	return nil
}

func (c Client) IsGameServerAvailable(namespace, name string) bool {
	log.Debug().Msgf("Checking if game server %s is terminating", name)

	gameserver, err := c.GetGameServer(namespace, name)
	if err != nil {
		return false
	}

	pod, err := c.BaseClient.GetPod(namespace, gameserver.Spec.Template.ObjectMeta.Name)
	if err != nil {
		return false
	}

	return pod.Status.Phase == v12.PodRunning
}
