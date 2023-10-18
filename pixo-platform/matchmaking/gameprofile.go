package matchmaking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MatchRequestParams struct {
	ModuleID      int    `json:"moduleId"`
	OrgID         int    `json:"orgId"`
	ClientVersion string `json:"clientVersion"`
}

type TicketRequestParams struct {
	MatchRequestParams
	Engine        string `json:"engine"`
	ServerVersion string `json:"serverVersion"`
	ImageRegistry string `json:"imageRegistry"`
	Status        string `json:"status"`
	Capacity      int    `json:"capacity"`
}

type GameProfileRepository struct {
	Client *redis.Client
	ctx    context.Context
}

func getProfileKey(ticketRequest TicketRequestParams) string {
	return fmt.Sprintf("profile:%d%d%s", ticketRequest.OrgID, ticketRequest.ModuleID, ticketRequest.ClientVersion)
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

func (r *GameProfileRepository) SaveProfile(ticketRequest TicketRequestParams) error {
	data, err := json.Marshal(ticketRequest)
	if err != nil {
		return err
	}
	key := getProfileKey(ticketRequest)
	err = r.Client.Set(r.ctx, key, data, getDuration()).Err()
	return err
}

func (r *GameProfileRepository) GetAllProfiles() ([]TicketRequestParams, error) {
	var profiles []TicketRequestParams

	keys, _, err := r.Client.Scan(r.ctx, 0, "profile:*", 100).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		data, err := r.Client.Get(r.ctx, key).Result()
		if err != nil {
			return nil, err
		}
		var profile TicketRequestParams
		err = json.Unmarshal([]byte(data), &profile)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
