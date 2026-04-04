package common

import (
	"strings"
	"time"
)

// GetStringValue returns a trimmed string from the input pointer.
// If the input is nil, it returns an empty string.
func GetStringValue(input *string) string {
	if input == nil {
		return ""
	}

	return strings.Trim(*input, " ")
}

// GetUTCStringForDateTime converts a time.Time pointer to a UTC string in RFC3339 format.
// If the input is nil, it returns an empty string.
func GetUTCStringForDateTime(input *time.Time) string {
	if input == nil {
		return ""
	}

	return input.UTC().Format(time.RFC3339)
}

// IsNilOrEmpty checks if a string pointer is nil or if the string is empty.
// Returns true if the string is nil or empty, false otherwise.
func IsNilOrEmpty(input *string) bool {
	return input == nil || IsEmptyString(*input)
}

// IsEmptyString checks if a string is empty after trimming whitespace.
// Returns true if the string is empty, false otherwise.
func IsEmptyString(input string) bool {
	return strings.TrimSpace(input) == ""
}

// IsNilOrZero is a generic function that checks if a numeric pointer is nil or has a zero value.
// It works with int, float32, and float64 types.
// Returns true if the value is nil or zero, false otherwise.
func IsNilOrZero[T int | float32 | float64](input *T) bool {
	return input == nil || *input == 0
}
