package cache

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/request"
)

func (r *GameProfileCacheClient) SaveProfile(ctx context.Context, profile request.MultiplayerMatchProfile) error {

	data, err := profile.MarshalJSON()
	if err != nil {
		return err
	}

	key := r.getCacheKey(profile)

	return r.Cache.Set(ctx, key, data, r.getCacheDuration()).Err()
}
