package request

import (
	"encoding/json"
	"fmt"
)

type MatchRequest struct {
	BaseTicketRequest
	ModuleID      int    `json:"moduleId"`
	OrgID         int    `json:"orgId"`
	ServerVersion string `json:"serverVersion"`
}

func (m *MatchRequest) GetLabel() string {
	return fmt.Sprintf("o-%d-m-%d-v-%s", m.OrgID, m.ModuleID, m.ServerVersion)
}

func (m *MatchRequest) MarshalJSON() ([]byte, error) {
	type Alias MatchRequest
	return json.Marshal(&struct{ *Alias }{
		Alias: (*Alias)(m),
	})
}

func (m *MatchRequest) UnmarshalJSON(data []byte) error {
	type Alias MatchRequest
	aux := &struct{ *Alias }{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}
