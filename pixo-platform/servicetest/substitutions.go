package servicetest

import (
	"fmt"
	"strings"
)

func (s *ServerTestFeature) replaceSubstitutions(data []byte) []byte {
	substitutions := map[string]string{
		"$ID":      fmt.Sprint(s.GraphQLResponse["id"]),
		"$USER_ID": fmt.Sprint(s.UserID),
	}

	for key, value := range substitutions {
		data = []byte(strings.ReplaceAll(string(data), key, value))
	}

	return data
}

func ReplaceRandomID(data []byte) []byte {
	randomID := GetRandomID()
	formattedString := strings.ReplaceAll(string(data), "$RANDOM_ID", randomID)
	return []byte(formattedString)
}
