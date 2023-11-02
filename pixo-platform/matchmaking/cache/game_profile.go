package cache

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/request"
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

func (r *GameProfileCacheClient) SaveProfile(ctx context.Context, profile request.MultiplayerMatchProfile) error {
	data, err := profile.MarshalJSON()
	if err != nil {
		return err
	}

	key := r.getCacheKey(profile)

	return r.Cache.Set(ctx, key, data, r.getCacheDuration()).Err()
}
