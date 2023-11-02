package cache

import (
	"context"
	"fmt"
)

func (r *GameProfileCacheClient) GetAllProfiles(ctx context.Context) ([]string, error) {
	var profiles []string

	keys, _, err := r.Cache.Scan(ctx, 0, fmt.Sprintf("%s:*", r.getLabel()), 100).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {

		data, err := r.Cache.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, data)
	}

	return profiles, nil
}
