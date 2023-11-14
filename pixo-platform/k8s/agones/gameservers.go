package agones

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (c Client) CreateGameServer(ctx context.Context, namespace string, gameserver *v1.GameServer) (*v1.GameServer, error) {
	log.Debug().Msg("Creating game server")

	if gameserver == nil {
		return nil, errors.New("gameserver is nil")
	}

	if gameserver.Labels == nil {
		gameserver.Labels = make(map[string]string)
	}

	gameserver.Labels[DeletedGameServerLabel] = "false"

	newGameServer, err := c.Clientset.
		AgonesV1().
		GameServers(namespace).
		Create(ctx, gameserver, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to create a game server")
		return nil, err
	}

	maxWaitSeconds := 60
	for {
		if maxWaitSeconds == 0 {
			return nil, errors.New("timed out waiting for game server reach ready state")
		}

		if c.IsGameServerAvailable(ctx, namespace, newGameServer.Name) {
			break
		}

		maxWaitSeconds--
		time.Sleep(time.Second * 1)
	}

	return newGameServer, nil
}

func (c Client) DeleteGameServer(ctx context.Context, namespace, name string) error {
	log.Debug().Msgf("Deleting game server %s", name)

	gameserver, err := c.GetGameServer(ctx, namespace, name)
	if err != nil {
		return nil
	}

	if _, err = c.AddLabelToGameServer(ctx, gameserver, DeletedGameServerLabel, "true"); err != nil {
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

	return gameserver.Labels[DeletedGameServerLabel] != "true"
}
