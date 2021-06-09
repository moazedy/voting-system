package logic

import (
	"errors"
	"voting-system/constants"

	"github.com/google/uuid"
)

func IdValidation(Id string) error {
	if Id == "" {
		return errors.New(constants.IdCanNotBeEmpty)
	}

	_, err := uuid.Parse(Id)
	if err != nil {
		return errors.New(constants.InvalidId)
	}

	return nil
}
