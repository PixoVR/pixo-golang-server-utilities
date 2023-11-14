package cache

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/request"
	"time"
)

func (r *GameProfileCacheClient) getCacheDuration() time.Duration {
	return time.Minute * 15
}

func (r *GameProfileCacheClient) getLabel() string {
	return DefaultLabel
}

func (r *GameProfileCacheClient) getCacheKey(req request.MultiplayerMatchProfile) string {
	return fmt.Sprintf("%s:%s", r.getLabel(), req.GetLabel())
}
