package client

import (
	"fmt"
	"math/rand"
	"strings"
)

func SanitizeFilename(timestamp int64, originalFilename string) (sanitizedFilename string) {
	if originalFilename == "" {
		randomInt := rand.Intn(100000)
		originalFilename = fmt.Sprintf("unnamed_file_%d", randomInt)
	}

	originalFilename = strings.ReplaceAll(originalFilename, " ", "_")
	nameParts := strings.Split(originalFilename, ".")

	if len(nameParts) < 2 {
		sanitizedFilename = fmt.Sprintf("%s_%d", originalFilename, timestamp)
	} else {

		filename := strings.Join(nameParts[:len(nameParts)-1], ".")
		fileExtension := nameParts[len(nameParts)-1]

		sanitizedFilename = fmt.Sprintf("%s_%d.%s", filename, timestamp, fileExtension)
	}

	return sanitizedFilename
}
