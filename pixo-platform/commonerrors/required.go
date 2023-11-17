package common_errors

import (
	"errors"
)

func ErrorRequired(fieldName string) error {
	return errors.New(fieldName + " is required")
}
