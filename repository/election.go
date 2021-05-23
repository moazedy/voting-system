package repository

import (
	"context"
	"voting-system/domain/models"
)

//  ElectionRepo is interface of  Election entity in repository layer. other layers of system can interface with it using this
// interface here
type ElectionRepo interface {
	// NewElection saves a new election in db
	NewElection(ctx context.Context, newElection models.Election) error
	// ReadElectionData reads some election's data in db using given Id
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	// DeleteElection deletes given election
	DeleteElection(ctx context.Context, electionId string) error
	// UpdateElection updates some election using received new election data
	UpdateElection(ctx context.Context, electionData models.Election) error
	// GetElectionContributorsCount gets count of given election's contributors
	GetElectionContributorsCount(ctx context.Context, electionId string) (*int, error)
}

// election is a struct that represents election entity in repository layer and its the way we can access to repository methods of
// election in this layer
type election struct {
}

// NewEelectionRepo is constractor fucntion for ElectionRepo
func NewEelectionRepo() ElectionRepo {
	return new(election)
}

func (e *election) NewElection(ctx context.Context, newElection models.Election) error {
	//TODO
	return nil
}
func (e *election) ReadElectionData(ctx context.Context, electionId string) (*models.Election, error) {
	//TODO
	return nil, nil
}

func (e *election) DeleteElection(ctx context.Context, electionId string) error {
	//TODO
	return nil
}
func (e *election) UpdateElection(ctx context.Context, electionData models.Election) error {
	//TODO
	return nil
}
func (e *election) GetElectionContributorsCount(ctx context.Context, electionId string) (*int, error) {
	//TODO
	return nil, nil
}
