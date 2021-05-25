package logic

import (
	"context"
	"errors"
	"time"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/repository"

	"github.com/google/uuid"
)

type ElectionLogic interface {
	CreateNewElection(ctx context.Context, userId string, electionData models.Election) error
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error
}

type election struct {
	repo repository.ElectionRepo
}

func NewElectionLogic() ElectionLogic {
	return new(election)
}

func (e *election) CreateNewElection(ctx context.Context, userId string, electionData models.Election) error {
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	electionData.CreatorId = userId
	electionData.CreationTime = time.Now()

	if err := electionData.Validate(); err != nil {
		return err
	}

	if err := e.repo.SaveNewElection(ctx, electionData); err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (e *election) ReadElectionData(ctx context.Context, electionId string) (*models.Election, error) {
	if err := electionIdValidate(electionId); err != nil {
		return nil, err
	}

	wantedElection, err := e.repo.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return wantedElection, nil
}

func (e *election) DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error {
	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return errors.New(constants.AccessDenied)
		}
	}

	if err := e.repo.DeleteElection(ctx, electionId); err != nil {
		return err
	}

	return nil
}

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
