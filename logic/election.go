package logic

import (
	"errors"
	"time"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/repository"
)

type ElectionLogic interface {
	CreateNewElection(userId string, electionData models.Election) error
}

type election struct {
	repo repository.ElectionRepo
}

func NewElectionLogic() ElectionLogic {
	return new(election)
}

func (e *election) CreateNewElection(userId string, electionData models.Election) error {
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	electionData.CreatorId = userId
	electionData.CreationTime = time.Now()
}
