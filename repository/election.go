package repository

import (
	"context"
	"log"
	"voting-system/domain/models"
	"voting-system/repository/couchbaseQueries"

	"github.com/couchbase/gocb/v2"
)

//  ElectionRepo is interface of  Election entity in repository layer. other layers of system can interface with it using this
// interface here
type ElectionRepo interface {
	// SaveNewElection saves a new election in db
	SaveNewElection(ctx context.Context, newElection models.Election) error
	// ReadElectionData reads some election's data in db using given Id
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	// DeleteElection deletes given election
	DeleteElection(ctx context.Context, electionId string) error
	// UpdateElection updates some election using received new election data
	UpdateElection(ctx context.Context, electionData models.Election) error
	// GetElectionContributorsCount gets count of given election's contributors
	GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error)
}

// election is a struct that represents election entity in repository layer and its the way we can access to repository methods of
// election in this layer
type election struct {
}

// NewEelectionRepo is constractor fucntion for ElectionRepo
func NewElectionRepo() ElectionRepo {
	return new(election)
}

func (e *election) SaveNewElection(ctx context.Context, newElection models.Election) error {
	_, err := DBS.Couch.Query(couchbaseQueries.SaveElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newElection.Id, newElection},
	})
	if err != nil {
		log.Println(" error in saving new election, error :", err.Error())
		return err
	}
	return nil
}

func (e *election) ReadElectionData(ctx context.Context, electionId string) (*models.Election, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ReadElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in reading election data, error :", err.Error())
		return nil, err
	}
	var elec models.Election
	err = result.One(&elec)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &elec, nil
		}
		log.Println("error in reading election item, error :", err.Error())
		return nil, err
	}

	return &elec, nil
}

func (e *election) DeleteElection(ctx context.Context, electionId string) error {
	_, err := DBS.Couch.Query(couchbaseQueries.DeleteElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println(" error in deleting election, error :", err.Error())
		return err
	}
	return nil
}

func (e *election) UpdateElection(ctx context.Context, electionData models.Election) error {
	_, err := DBS.Couch.Query(couchbaseQueries.UpdateElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{
			electionData.Title,
			electionData.StartTime,
			electionData.EndTime,
			electionData.HasEnded,
			electionData.Type,
			electionData.CandidatesCountLimit,
			electionData.CreatorId,
			electionData.Id,
		},
	})
	if err != nil {
		log.Println(" error in updating election, error :", err.Error())
		return err
	}
	return nil
}

func (e *election) GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetElectionContributorsCountQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var count models.ContributorsCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &count, nil
		}

		log.Println("error in reading contributors count, error :", err.Error())
		return nil, err
	}

	return &count, nil
}
