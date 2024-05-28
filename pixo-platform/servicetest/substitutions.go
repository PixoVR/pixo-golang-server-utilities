package servicetest

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func (s *ServerTestFeature) PerformSubstitutions(data []byte) []byte {
	substitutions := map[string]string{
		"ID":          fmt.Sprint(s.ID),
		"USER_ID":     fmt.Sprint(s.UserID),
		"RANDOM_INT":  fmt.Sprint(s.RandomInt),
		"RANDOM_ID":   generateRandomID(),
		"RANDOM_UUID": generateRandomUUID(),
	}

	for key, value := range substitutions {
		data = replace(data, key, value)
	}

	for key, value := range s.staticSubstitutions {
		data = replace(data, key, value)
	}

	for key, value := range s.dynamicSubstitutions {
		if value != nil {
			data = replace(data, key, value(data))
		}
	}

	var configMap []map[string]string
	if err := viper.UnmarshalKey("substitutions", &configMap); err == nil {
		for _, pair := range configMap {
			data = replace(data, pair["key"], pair["value"])
		}
	}

	return data
}

func replace(data []byte, key, value string) []byte {
	return []byte(strings.ReplaceAll(string(data), fmt.Sprintf("$%s", key), value))
}
