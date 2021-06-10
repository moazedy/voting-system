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
	SaveNewContributor(ctx context.Context, newContributor models.Contributor) (*models.Id, error)
	// GetListOfContributorsInAnElection gets list of all contributors in an election
	GetListOfContributorsInAnElection(ctx context.Context, electionId, order string, offset, limit int) ([]models.Contributor, error)
	// ReadContributorData reads a contributor's data from db using receiving id
	ReadContributorData(ctx context.Context, contributorMetaId string) (*models.Contributor, error)
	// DeleteContributor deletes a contributor using given id
	DeleteContributor(ctx context.Context, contributorMetaId string) error
	// IsContributorExists checks for a specific participation of some contributor existance
	IsContributorExists(ctx context.Context, contributorId, electionId string) (*bool, error)
	// IsContributionExists checks on some contribution existance using meta id
	IsContributionExists(ctx context.Context, contributorMetaId string) (*bool, error)
}

// contributor is a struct that represents contributor entity in repository layer and its the way we can access to repository methods of
// contributor in this layer
type contributor struct {
}

// NewContributorRepo is constractor function for ContributorRepo
func NewContributorRepo() ContributorRepo {
	return new(contributor)
}

func (c *contributor) SaveNewContributor(ctx context.Context, newContributor models.Contributor) (*models.Id, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.SaveContributorQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newContributor.MetaId, newContributor},
	})
	if err != nil {
		log.Println(" error in saving new contributor, error :", err.Error())
		return nil, err
	}

	var id models.Id
	err = result.One(&id)
	if err != nil {
		log.Println("error on reading contributor id value, error: ", err.Error())
		return nil, err
	}
	return &id, nil
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

func (c *contributor) ReadContributorData(ctx context.Context, contributorMetaId string) (*models.Contributor, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ReadContributorQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{contributorMetaId},
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

func (c *contributor) DeleteContributor(ctx context.Context, contributorMetaId string) error {
	_, err := DBS.Couch.Query(couchbaseQueries.DeleteContributorQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{contributorMetaId},
	})
	if err != nil {
		log.Println(" error in deleting contributor, error :", err.Error())
		return err
	}
	return nil
}

func (c contributor) IsContributorExists(ctx context.Context, contributorId, electionId string) (*bool, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ContributorExistanceQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{contributorId, electionId},
	})
	if err != nil {
		log.Println("error in query execution in contributor existance check. error :", err.Error())
		return nil, err
	}

	var count models.ContributorsCount
	var exists bool
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &exists, nil
		}
		log.Println("error in reading contributor count data. error :", err.Error())
		return nil, err
	}

	if count.Count > 0 {
		exists = true
	}

	return &exists, nil
}

func (c contributor) IsContributionExists(ctx context.Context, contributorMetaId string) (*bool, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ContributionExistanceQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{contributorMetaId},
	})
	if err != nil {
		log.Println("error in query execution in contributor existance check. error :", err.Error())
		return nil, err
	}

	var count models.ContributorsCount
	var exists bool
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &exists, nil
		}
		log.Println("error in reading contributor count data. error :", err.Error())
		return nil, err
	}

	if count.Count > 0 {
		exists = true
	}

	return &exists, nil
}
