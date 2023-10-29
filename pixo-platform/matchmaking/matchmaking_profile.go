package matchmaking

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rs/zerolog/log"
	"open-match.dev/open-match/pkg/pb"
)

const (
	MaxPlayersExtensionKey = "max-players"
)

type MultiplayerMatchmakingProfile struct {
	*pb.MatchProfile
}

func NewPixoMatchmakingProfile(p *pb.MatchProfile) *MultiplayerMatchmakingProfile {
	return &MultiplayerMatchmakingProfile{MatchProfile: p}
}

func (p *MultiplayerMatchmakingProfile) GetMaxPlayersPerMatch() int {
	var val wrappers.Int32Value
	if err := ptypes.UnmarshalAny(p.MatchProfile.GetExtensions()[MaxPlayersExtensionKey], &val); err != nil {
		log.Error().Err(err).Msg("Cannot retrieve number of players per match")
		return 1
	} else {
		log.Debug().Msgf("Max players per match: %d", val.Value)
		return int(val.Value)
	}
}
