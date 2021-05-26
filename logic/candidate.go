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

// CandidateLogic is interface for candidate entity in logic layer
type CandidateLogic interface {
	// CreateNewCandidate creates new candidate
	CreateNewCandidate(ctx context.Context, candidateData models.Candidate) error
	// ReadCandidateData reads data of givne candidate
	ReadCandidateData(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) (*models.Candidate, error)
	// DeleteCandidate deletes given candidateId
	DeleteCandidate(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) error
	// CandidateExistanceCheck checks for some candidate existance in db
	CandidateExistanceCheck(ctx context.Context, candidateId string) (*bool, error)
	// GetListOfSomeElectionCandidates gets list of all election candidates
	GetListOfElectionCandidates(ctx context.Context, electionId, order string, offset, limit int) ([]models.Candidate, error)
}

// candidate is a struct that holds methods for candidate in logic
type candidate struct {
	repo repository.CandidateRepo
}

// NewCandidateLogic is constractor function for CandidateLogic
func NewCandidateLogic() CandidateLogic {
	newCandidate := new(candidate)
	if newCandidate.repo == nil {
		newCandidate.repo = repository.NewCandidateRepo()
	}

	return newCandidate
}

func (c *candidate) CreateNewCandidate(ctx context.Context, candidateData models.Candidate) error {
	if err := candidateData.Validate(); err != nil {
		return err
	}

	candidateData.Id = uuid.New()
	candidateData.Created_At = time.Now()
	electionExists, err := repository.NewElectionRepo().ElectionExists(ctx, candidateData.Id.String())
	if err != nil {
		return errors.New(constants.InternalServerError)
	}
	if !*electionExists {
		return errors.New(constants.ElectionDoesNotExist)
	}

	err = c.repo.CreateCandidate(ctx, candidateData)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (c *candidate) ReadCandidateData(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) (*models.Candidate, error) {
	_, err := c.CandidateExistanceCheck(ctx, candidateId)
	if err != nil {
		return nil, err
	}
	theCandidate, err := c.repo.ReadCandidate(ctx, candidateId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	if !requestedByAdmin {
		if theCandidate.Id.String() != requesterId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	return theCandidate, nil

}

func (c *candidate) DeleteCandidate(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) error {
	_, err := c.ReadCandidateData(ctx, candidateId, requesterId, false)
	if err != nil {
		return err
	}

	err = c.repo.DeleteCandidate(ctx, candidateId)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (c *candidate) CandidateExistanceCheck(ctx context.Context, candidateId string) (*bool, error) {
	if err := candidateIdValidate(candidateId); err != nil {
		return nil, err
	}

	exists, err := c.repo.IsCandidateExists(ctx, candidateId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	if !*exists {
		return exists, errors.New(constants.CandidateDoesNotExist)
	}

	return exists, nil
}

func (c *candidate) GetListOfElectionCandidates(ctx context.Context, electionId, order string, offset, limit int) ([]models.Candidate, error) {

	contributors, err := c.repo.GetListOfElectionCandidates(ctx, electionId, order, offset, limit)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}
	return contributors, nil
}
