package repository

import (
	"context"
	"voting-system/domain/models"
)

// ContributorRepo is interface of contributor  entity in repository layer. other layers of system can interface with it using this
// interface here
type ContributorRepo interface {
	// NewContributor is creating new contributor entity in the db
	NewContributor(ctx context.Context, newContributor models.Contributor) error
	// GetListOfContributorsInAnElection gets list of all contributors in an election
	GetListOfContributorsInAnElection(ctx context.Context, electionId, order string, offset, limit int) ([]models.Contributor, error)
	// ReadContributorData reads a contributor's data from db using receiving id
	ReadContributorData(ctx context.Context, contributorId string) (*models.Contributor, error)
	// DeleteContributor deletes a contributor using given id
	DeleteContributor(ctx context.Context, contributorId string) error
}

// contributor is a struct that represents contributor entity in repository layer and its the way we can access to repository methods of
// contributor in this layer
type contributor struct {
}

// NewContributorRepo is constractor function for ContributorRepo
func NewContributorRepo() ContributorRepo {
	return new(contributor)
}

func (c *contributor) NewContributor(ctx context.Context, newContributor models.Contributor) error {
	// TODO
	return nil
}

func (c *contributor) GetListOfContributorsInAnElection(ctx context.Context, electionId, order string, offset, limit int) ([]models.Contributor, error) {
	// TODO
	return nil, nil

}

func (c *contributor) ReadContributorData(ctx context.Context, contributorId string) (*models.Contributor, error) {
	// TODO
	return nil, nil

}

func (c *contributor) DeleteContributor(ctx context.Context, contributorId string) error {
	// TODO
	return nil

}
