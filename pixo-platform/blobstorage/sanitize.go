package blobstorage

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	filenameOverride = "blob"
)

func SanitizeFilename(originalFilename string, timestamp int64) (sanitizedFilename string) {

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
		sanitizedFilename = formatFilename(originalFilename, timestamp)
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

		sanitizedFilename = fmt.Sprintf("%s.%s", formatFilename(filename, timestamp), fileExtension)
	}

	return sanitizedFilename
}

func formatFilename(filename string, timestamp int64) string {
	if timestamp == 0 {
		return filename
	}

	return fmt.Sprintf("%s_%d", filename, timestamp)
}

func ParseFileLocationFromLink(link string) string {
	if !strings.Contains(link, "https://") {
		return link
	}

	splitLink := strings.Split(link, "https://")
	if len(splitLink) < 2 {
		return ""
	}

	filePath := strings.Split(splitLink[1], "/")
	filePath = filePath[1:]
	return strings.Join(filePath, "/")
}

func GetFilenameFromLocation(location string) string {
	parts := strings.Split(location, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
