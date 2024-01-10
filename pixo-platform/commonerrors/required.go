package common_errors

import (
	"errors"
)

func ErrorRequired(fieldName string) error {
	return errors.New(fieldName + " is required")
}

const DefaultRecordNotFoundMessage = "record not found"

func ErrorNotFound(objectType string) error {
	return errors.New(objectType + " not found")
}
