package logic

import (
	"errors"
	"voting-system/constants"

	"github.com/google/uuid"
)

func electionIdValidate(Id string) error {
	if Id == "" {
		return errors.New(constants.InvalidElectionId)
	}

	_, err := uuid.Parse(Id)
	if err != nil {
		return errors.New(constants.InvalidElectionId)
	}

	return nil
}
