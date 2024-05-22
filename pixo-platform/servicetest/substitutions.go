package servicetest

import (
	"fmt"
	"reflect"
	"strings"
)

func (s *ServerTestFeature) PerformSubstitutions(data []byte) []byte {
	substitutions := map[string]string{
		"$ID":          fmt.Sprint(s.GraphQLResponse["id"]),
		"$USER_ID":     fmt.Sprint(s.UserID),
		"$RANDOM_ID":   generateRandomID(),
		"$RANDOM_UUID": generateRandomUUID(),
	}

	for key, value := range substitutions {
		data = []byte(strings.ReplaceAll(string(data), key, value))
	}

	for key, value := range s.staticSubstitutions {
		data = []byte(strings.ReplaceAll(string(data), key, value))
	}

	for key, value := range s.dynamicSubstitutions {
		if value != nil && reflect.ValueOf(value).Kind() == reflect.Func {
			data = []byte(strings.ReplaceAll(string(data), key, string(value(data))))
		}
	}

	return data
}
