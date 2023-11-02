package request

import "encoding/json"

type MultiplayerMatchProfile interface {
	json.Marshaler
	json.Unmarshaler
	GetLabel() string
}
