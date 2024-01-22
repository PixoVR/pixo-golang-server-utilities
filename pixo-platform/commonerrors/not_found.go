package common_errors

import (
	"errors"
)

const DefaultRecordNotFoundMessage = "record not found"

func ErrorNotFound(objectType string) error {
	return errors.New(objectType + " not found")
}
