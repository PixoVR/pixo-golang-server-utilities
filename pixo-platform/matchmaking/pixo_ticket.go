package matchmaking

import (
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rs/zerolog/log"
	"open-match.dev/open-match/pkg/pb"
)

type PixoTicket struct {
	*pb.Ticket
}

const (
	TicketMatchAttemptExtensionKey = "ticket-match-attempt-count"
	ErrNoTicketProvided            = "no ticket provided"
)

func NewPixoTicket(ticket *pb.Ticket) *PixoTicket {
	return &PixoTicket{Ticket: ticket}
}

func (p *PixoTicket) SetMatchmakingAttemptCount(inputVal int32) error {
	if p.PersistentField == nil {
		p.PersistentField = make(map[string]*any.Any)
	}

	val, err := ptypes.MarshalAny(&wrappers.Int32Value{Value: inputVal})
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshal matchmaking attempt count")
		return err
	}

	p.PersistentField[TicketMatchAttemptExtensionKey] = val

	log.Debug().Msgf("Attempt count set to %v", inputVal)
	return nil
}

func (p *PixoTicket) GetMatchmakingAttemptCount() (int32, error) {
	if p == nil {
		err := errors.New(ErrNoTicketProvided)
		return 0, err
	}

	if p.PersistentField != nil {
		if rawVal, ok := p.PersistentField[TicketMatchAttemptExtensionKey]; ok {
			var val wrappers.Int32Value
			if err := ptypes.UnmarshalAny(rawVal, &val); err != nil {
				log.Error().Err(err).Msg("Unable to unmarshal open slots")
				return 0, err
			}

			log.Debug().Msgf("open slots: %v", val.Value)
			return val.Value, nil
		}
	}

	return 0, nil
}
