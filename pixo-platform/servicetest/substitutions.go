package servicetest

import (
	"fmt"
	"strings"
)

func (s *ServerTestFeature) replaceSubstitutions(data []byte) []byte {
	substitutions := map[string]string{
		"$ID":          fmt.Sprint(s.GraphQLResponse["id"]),
		"$USER_ID":     fmt.Sprint(s.UserID),
		"$RANDOM_ID":   generateRandomID(),
		"$RANDOM_UUID": generateRandomUUID(),
	}

	for key, value := range substitutions {
		data = []byte(strings.ReplaceAll(string(data), key, value))
	}

	return data
}
