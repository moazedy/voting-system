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

// ElectionLogic is election entity interface in logic layer and other layers of program can intract with it
// through this interface here
type ElectionLogic interface {
	// CreateNewElection craates new election using receieved election data
	CreateNewElection(ctx context.Context, userId string, electionData models.Election) error
	// ReadElectionData reads data of given election id
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	// DeleteElection deletes the election with received election id
	DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error
	// GetElectionContributorsCount gets number of persons who contributed in given election id
	GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error)
	// UpdateElection updates data of some election using received electionData
	UpdateElection(ctx context.Context, electionData models.Election, requesterId string, requestedByAdmin bool) error
	// CheckElectionExistance checks on election id existance in db
	CheckElectionExistance(ctx context.Context, electionId string) (*bool, error)
	// GetListOfRelatedUsers gets list of user Ids being added to the election as related users
	GetListOfRelatedUsers(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*models.RelatedUsers, error)
}

// election struct, is holder of election metods in logic layer
type election struct {
	// repo is the way this layer of program can interface with repository layer
	repo repository.ElectionRepo
}

// NewElectionLogic is constractor function of ElectionLogic
func NewElectionLogic() ElectionLogic {
	return new(election)
}

func (e *election) CreateNewElection(ctx context.Context, userId string, electionData models.Election) error {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	electionData.Id = uuid.New()
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
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if err := electionIdValidate(electionId); err != nil {
		return nil, err
	}

	_, err := e.CheckElectionExistance(ctx, electionId)
	if err != nil {
		return nil, err
	}

	wantedElection, err := e.repo.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return wantedElection, nil
}

func (e *election) DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

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
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

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
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

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

func (e *election) CheckElectionExistance(ctx context.Context, electionId string) (*bool, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if err := electionIdValidate(electionId); err != nil {
		return nil, err
	}

	exists, err := e.repo.ElectionExists(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	if !*exists {
		return exists, errors.New(constants.ElectionDoesNotExist)
	}

	return exists, nil

}

func (e *election) GetListOfRelatedUsers(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*models.RelatedUsers, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}
	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	users, err := e.repo.GetListOfRelatedUsers(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return users, nil
}
