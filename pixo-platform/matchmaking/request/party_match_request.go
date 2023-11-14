package request

import (
	"encoding/json"
	"fmt"
)

type PartyMatchRequest struct {
	BaseTicketRequest
	PartyCode string `json:"partyCode"`
}

func (p *PartyMatchRequest) GetLabel() string {
	return fmt.Sprintf("p-%s", p.PartyCode)
}

func (p *PartyMatchRequest) MarshalJSON() ([]byte, error) {
	type Alias PartyMatchRequest
	return json.Marshal(&struct{ *Alias }{
		Alias: (*Alias)(p),
	})
}

func (p *PartyMatchRequest) UnmarshalJSON(data []byte) error {
	type Alias PartyMatchRequest
	aux := &struct{ *Alias }{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}
