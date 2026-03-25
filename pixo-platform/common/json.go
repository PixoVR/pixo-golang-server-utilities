package common

import (
	"encoding/json"
	"fmt"
)

func ToJSONString[T any](v T) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}
