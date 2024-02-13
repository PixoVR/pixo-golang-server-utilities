package client

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	filenameOverride = "blob"
)

func SanitizeFilename(timestamp int64, originalFilename string) (sanitizedFilename string) {

	if originalFilename == "" {
		randomInt := rand.Intn(100000)
		originalFilename = fmt.Sprintf("unnamed_file_%d", randomInt)
	}

	if originalFilename[len(originalFilename)-1] == '/' {
		originalFilename += filenameOverride
	}

	originalFilename = strings.ReplaceAll(originalFilename, " ", "_")
	nameParts := strings.Split(originalFilename, ".")

	if len(nameParts) < 2 {
		sanitizedFilename = fmt.Sprintf("%s_%d", originalFilename, timestamp)
	} else {
		pathParts := strings.Split(nameParts[0], "/")
		path := strings.Join(pathParts[:len(pathParts)-1], "/")

		var filename string
		if len(pathParts) > 1 {
			filename = fmt.Sprintf("%s/%s", path, filenameOverride)
		} else {
			filename = filenameOverride
		}

		fileExtension := nameParts[len(nameParts)-1]

		sanitizedFilename = fmt.Sprintf("%s_%d.%s", filename, timestamp, fileExtension)
	}

	return sanitizedFilename
}
