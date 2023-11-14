package cache

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewMiniGameProfileCache() (*GameProfileCacheClient, *miniredis.Miniredis, *redis.Client, error) {

	s, err := miniredis.Run()
	if err != nil {
		return nil, nil, nil, err
	}

	c := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	gameProfileCache := NewGameProfileCacheClient(c)

	return gameProfileCache, s, c, nil
}
