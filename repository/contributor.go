package repository

import (
	"context"
	"fmt"
	"log"
	"voting-system/domain/models"
	"voting-system/repository/couchbaseQueries"

	"github.com/couchbase/gocb/v2"
)

// ContributorRepo is interface of contributor  entity in repository layer. other layers of system can interface with it using this
// interface here
type ContributorRepo interface {
	// SaveNewContributor is creating new contributor entity in the db
	SaveNewContributor(ctx context.Context, newContributor models.Contributor) error
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

func (c *contributor) SaveNewContributor(ctx context.Context, newContributor models.Contributor) error {
	_, err := DBS.Couch.Query(couchbaseQueries.SaveContributorQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newContributor.Id, newContributor},
	})
	if err != nil {
		log.Println(" error in saving new contributor, error :", err.Error())
		return err
	}
	return nil
}

func (c *contributor) GetListOfContributorsInAnElection(ctx context.Context, electionId, order string, offset, limit int) ([]models.Contributor, error) {
	query := fmt.Sprintf(couchbaseQueries.GetElectionContributorsQuery, order)
	result, err := DBS.Couch.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId, offset, limit},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var contributors []models.Contributor
	for result.Next() {
		var con models.Contributor
		err := result.Row(&con)
		if err != nil {
			if err == gocb.ErrNoResult {
				return contributors, nil
			}
			log.Println("error in reading contributor item, error :", err.Error())
			return nil, err
		}

		contributors = append(contributors, con)
	}

	return contributors, nil
}

func (c *contributor) ReadContributorData(ctx context.Context, contributorId string) (*models.Contributor, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ReadContributorQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{contributorId},
	})
	if err != nil {
		log.Println("error in reading contributor data, error :", err.Error())
		return nil, err
	}
	var con models.Contributor
	err = result.One(&con)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &con, nil
		}
		log.Println("error in reading contributor item, error :", err.Error())
		return nil, err
	}

	return &con, nil
}

func (c *contributor) DeleteContributor(ctx context.Context, contributorId string) error {
	_, err := DBS.Couch.Query(couchbaseQueries.DeleteContributorQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{contributorId},
	})
	if err != nil {
		log.Println(" error in deleting contributor, error :", err.Error())
		return err
	}
	return nil
}
