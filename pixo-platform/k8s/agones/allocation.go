package agones

import (
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) CreateGameServerAllocation(namespace string, allocation *allocationv1.GameServerAllocation) (*allocationv1.GameServerAllocation, error) {
	gsa, err := c.AllocationV1().GameServerAllocations(namespace).Create(context.Background(), allocation, metav1.CreateOptions{})
	if err != nil {
		log.Err(err).Msg("Failed to allocate the game server")
		return nil, err
	}

	return gsa, nil
}
