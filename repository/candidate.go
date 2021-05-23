package repository

import (
	"context"
	"voting-system/domain/models"
)

// CandidateRepo is interface of candidate  entity in repository layer. other layers of system can interface with it using this
// interface here
type CandidateRepo interface {
	// CreateCandidate is for creating new candidate entitiy in voting system
	CreateCandidate(ctx context.Context, NewCandidate models.Candidate) error
	// ReadCandidate reads data of requested candidate id, if it exists in db
	ReadCandidate(ctx context.Context, cadidateId string) (*models.Candidate, error)
	// GetListOfSomeElectionCandidates gets list of all election candidates
	GetListOfSomeElectionCandidates(ctx context.Context, electionId, order string, offset, limit int) ([]models.Candidate, error)
	// DeleteCandidate deletes given candidate so that it can not be accessable for voting
	DeleteCandidate(ctx context.Context, candidateId string) error
	// UpdateCandidate updates candidate data using received data
	UpdateCandidate(ctx context.Context, candidateData models.Candidate) error
	// IsCandidateExists checks if given candidate id exists in system or not
	IsCandidateExists(ctx context.Context, candidateId string) (bool, error)
}

// candidate is a struct that represents candidate entity in repository layer and its the way we can access to repository methods of
// candidate in this layer
type candidate struct {
}

// NewCandidateRepo is constractor function for CandidateRepo
func NewCandidateRepo() CandidateRepo {
	return new(candidate)
}

func (c *candidate) CreateCandidate(ctx context.Context, NewCandidate models.Candidate) error {
	// TODO
	return nil
}

func (c *candidate) ReadCandidate(ctx context.Context, cadidateId string) (*models.Candidate, error) {
	// TODO
	return nil, nil

}

func (c *candidate) GetListOfSomeElectionCandidates(ctx context.Context, electionId, order string, offset, limit int) ([]models.Candidate, error) {
	// TODO
	return nil, nil

}

func (c *candidate) DeleteCandidate(ctx context.Context, candidateId string) error {
	// TODO
	return nil

}

func (c *candidate) UpdateCandidate(ctx context.Context, candidateData models.Candidate) error {
	// TODO
	return nil

}

func (c *candidate) IsCandidateExists(ctx context.Context, candidateId string) (bool, error) {
	// TODO
	return false, nil

}
