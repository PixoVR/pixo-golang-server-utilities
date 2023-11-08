package agones

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"errors"
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

func (c Client) GetGameServer(ctx context.Context, namespace, name string) (*v1.GameServer, error) {
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

func (c Client) CreateGameServer(ctx context.Context, namespace string, gameserver *v1.GameServer) (*v1.GameServer, error) {
	log.Debug().Msg("Creating game server")

	newGameServer, err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Create(ctx, gameserver, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to create a game server")
		return nil, err
	}

	maxWaitSeconds := 30
	for {
		if maxWaitSeconds == 0 {
			return nil, errors.New("timed out waiting for game server to be available")
		}

		if !c.IsGameServerAvailable(ctx, namespace, newGameServer.Name) {
			maxWaitSeconds--
		} else {
			break
		}
	}

	return newGameServer, err
}

func (c Client) DeleteGameServer(ctx context.Context, namespace, name string) error {
	log.Debug().Msgf("Deleting game server %s", name)

	gameserver, err := c.GetGameServer(ctx, namespace, name)
	if err != nil {
		return err
	}

	if _, err = c.AddLabelToGameServer(ctx, gameserver, "deleted", "true"); err != nil {
		return err
	}

	err = c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete game server %s", name)
		return err
	}

	return nil
}

func (c Client) AddLabelToGameServer(ctx context.Context, gameserver *v1.GameServer, key, value string) (*v1.GameServer, error) {
	log.Debug().Msgf("Adding label %s=%s to game server %s", key, value, gameserver.GetName())

	if gameserver == nil {
		return nil, errors.New("gameserver is nil")
	}

	retrievedGameserver, err := c.GetGameServer(ctx, gameserver.Namespace, gameserver.Name)
	if err != nil {
		return nil, err
	}

	if retrievedGameserver.Labels == nil {
		retrievedGameserver.Labels = make(map[string]string)
	} else {
		retrievedGameserver.Labels[key] = value
	}

	updatedGameserver, err := c.Clientset.
		AgonesV1().
		GameServers(retrievedGameserver.Namespace).
		Update(ctx, retrievedGameserver, metav1.UpdateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to update game server labels for gameserver: %s", updatedGameserver.GetName())
		return nil, err
	}

	return updatedGameserver, nil
}

func (c Client) IsGameServerAvailable(ctx context.Context, namespace, name string) bool {
	log.Debug().Msgf("Checking if game server %s is terminating", name)

	gameserver, err := c.GetGameServer(ctx, namespace, name)
	if err != nil {
		return false
	}

	pod, err := c.BaseClient.GetPod(ctx, namespace, name)
	if err != nil || pod == nil {
		return false
	}

	return pod.Status.Phase == v12.PodRunning && gameserver.Labels["deleted"] != "true"
}
