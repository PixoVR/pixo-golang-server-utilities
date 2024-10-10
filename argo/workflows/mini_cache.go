package workflows

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewMiniCache() (*redis.Client, *miniredis.Miniredis, error) {

	s, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	c := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	return c, s, nil
}
