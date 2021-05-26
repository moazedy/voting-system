package logic

import (
	"context"
	"errors"
	"time"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/repository"
)

type ElectionLogic interface {
	CreateNewElection(ctx context.Context, userId string, electionData models.Election) error
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error
	GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error)
	UpdateElection(ctx context.Context, electionData models.Election, requesterId string, requestedByAdmin bool) error
}

type election struct {
	repo repository.ElectionRepo
}

func NewElectionLogic() ElectionLogic {
	newElection := new(election)

	if newElection.repo == nil {
		newElection.repo = repository.NewElectionRepo()
	}
	return newElection
}

func (e *election) CreateNewElection(ctx context.Context, userId string, electionData models.Election) error {
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
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (e *election) GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error) {
	if err := electionIdValidate(electionId); err != nil {
		return nil, err
	}

	count, err := e.repo.GetElectionContributorsCount(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return count, nil
}

func (e *election) UpdateElection(ctx context.Context, electionData models.Election, requesterId string, requestedByAdmin bool) error {
	theElection, err := e.ReadElectionData(ctx, electionData.Id.String())
	if err != nil {
		return err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return errors.New(constants.AccessDenied)
		}
	}

	if err := electionData.Validate(); err != nil {
		return err
	}

	err = e.repo.UpdateElection(ctx, electionData)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}
