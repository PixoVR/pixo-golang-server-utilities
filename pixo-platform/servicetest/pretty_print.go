package servicetest

import (
	"bytes"
	"encoding/json"
	"strings"
)

// prettify formats the json string
func prettify(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "  ", " ")
	s = strings.TrimSpace(s)

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(s), "", "  "); err != nil {
		return s
	}

	return prettyJSON.String()
}
