package profile

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rs/zerolog/log"
	"open-match.dev/open-match/pkg/pb"
)

const (
	MaxPlayersExtensionKey = "max-players"
)

type OpenMatchProfile struct {
	*pb.MatchProfile
}

func NewOpenMatchProfile(p *pb.MatchProfile) *OpenMatchProfile {
	return &OpenMatchProfile{MatchProfile: p}
}

func (p *OpenMatchProfile) GetMaxPlayersPerMatch() int {
	var val wrappers.Int32Value
	if err := ptypes.UnmarshalAny(p.MatchProfile.GetExtensions()[MaxPlayersExtensionKey], &val); err != nil {
		log.Error().Err(err).Msg("Cannot retrieve number of players per match")
		return 100
	} else {
		log.Debug().Msgf("Max players per match: %d", val.Value)
		return int(val.Value)
	}
}
