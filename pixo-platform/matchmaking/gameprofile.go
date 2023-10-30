package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MatchRequest struct {
	ModuleID      int    `json:"moduleId"`
	OrgID         int    `json:"orgId"`
	ServerVersion string `json:"serverVersion"`
}

type TicketRequest struct {
	MatchRequest
	Capacity int `json:"capacity"`
}

type GameProfileRepository struct {
	ctx    context.Context
	Client *redis.Client
}

func getProfileKey(req MatchRequest) string {
	return fmt.Sprintf("profile:%d%d%s", req.OrgID, req.ModuleID, req.ServerVersion)
}

func getDuration() time.Duration {
	return time.Minute * 15
}

func NewGameProfileRepository(redisAddr, redisPassword string) *GameProfileRepository {
	return &GameProfileRepository{
		Client: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
		}),
		ctx: context.Background(),
	}
}

func (r *GameProfileRepository) SaveProfile(req MatchRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	key := getProfileKey(req)

	return r.Client.Set(r.ctx, key, data, getDuration()).Err()
}

func (r *GameProfileRepository) GetAllProfiles() ([]MatchRequest, error) {
	var profiles []MatchRequest

	keys, _, err := r.Client.Scan(r.ctx, 0, "profile:*", 100).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {

		data, err := r.Client.Get(r.ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var profile MatchRequest
		if err = json.Unmarshal([]byte(data), &profile); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
