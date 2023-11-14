package cache

import (
	"github.com/redis/go-redis/v9"
)

type GameProfileCacheClient struct {
	Cache *redis.Client
}

func NewGameProfileCacheClient(cache *redis.Client) *GameProfileCacheClient {
	return &GameProfileCacheClient{
		Cache: cache,
	}
}
