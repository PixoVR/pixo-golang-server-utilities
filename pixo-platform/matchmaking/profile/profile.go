package profile

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	var val wrapperspb.Int32Value
	if err := anypb.UnmarshalTo(p.MatchProfile.GetExtensions()[MaxPlayersExtensionKey], &val, proto.UnmarshalOptions{}); err != nil {
		log.Error().Err(err).Msg("Cannot retrieve number of players per match")
		return 100
	} else {
		log.Debug().Msgf("Max players per match: %d", val.Value)
		return int(val.Value)
	}
}
