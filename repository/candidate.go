package repository

import (
	"context"
	"fmt"
	"log"
	"voting-system/domain/models"
	"voting-system/repository/couchbaseQueries"

	"github.com/couchbase/gocb/v2"
)

// CandidateRepo is interface of candidate  entity in repository layer. other layers of system can interface with it using this
// interface here
type CandidateRepo interface {
	// CreateCandidate is for creating new candidate entitiy in voting system
	CreateCandidate(ctx context.Context, NewCandidate models.Candidate) (*models.Id, error)
	// ReadCandidate reads data of requested candidate id, if it exists in db
	ReadCandidate(ctx context.Context, candidateId string) (*models.Candidate, error)
	// GetListOfSomeElectionCandidates gets list of all election candidates
	GetListOfElectionCandidates(ctx context.Context, electionId, order string, offset, limit int) ([]models.Candidate, error)
	// DeleteCandidate deletes given candidate so that it can not be accessable for voting
	DeleteCandidate(ctx context.Context, candidateId string) error
	// UpdateCandidate updates candidate data using received data
	UpdateCandidate(ctx context.Context, candidateData models.Candidate) error
	// IsCandidateExists checks if given candidate id exists in system or not
	IsCandidateExists(ctx context.Context, candidateId string) (*bool, error)
	// GetAllElectionCandidates gets list of all candidates in an election
	GetAllElectionCandidates(ctx context.Context, electionId string) ([]models.Candidate, error)
}

// candidate is a struct that represents candidate entity in repository layer and its the way we can access to repository methods of
// candidate in this layer
type candidate struct {
}

// NewCandidateRepo is constractor function for CandidateRepo
func NewCandidateRepo() CandidateRepo {
	return new(candidate)
}

func (c *candidate) CreateCandidate(ctx context.Context, NewCandidate models.Candidate) (*models.Id, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.SaveNewCandidateQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{NewCandidate.Id, NewCandidate},
	})
	if err != nil {
		log.Println(" error in saving new candidate, error :", err.Error())
		return nil, err
	}

	var id models.Id
	err = result.One(&id)
	if err != nil {
		log.Println("error in reading candidate id value, error : ", err.Error())
		return nil, err
	}
	return &id, nil
}

func (c *candidate) ReadCandidate(ctx context.Context, candidateId string) (*models.Candidate, error) {
	var can models.Candidate
	result, err := DBS.Couch.Query(couchbaseQueries.ReadCandidateDataQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println("error in reading candidate data, error :", err.Error())
		return nil, err
	}
	err = result.One(&can)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &can, nil
		}
		log.Println("error in reading candidate item, error :", err.Error())
		return nil, err
	}

	return &can, nil
}

func (c *candidate) GetListOfElectionCandidates(ctx context.Context, electionId, order string, offset, limit int) ([]models.Candidate, error) {
	query := fmt.Sprintf(couchbaseQueries.GetElectionCandidatesQuery, order)
	result, err := DBS.Couch.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId, offset, limit},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var candidates []models.Candidate
	for result.Next() {
		var candid models.Candidate
		err := result.Row(&candid)
		if err != nil {
			if err == gocb.ErrNoResult {
				return candidates, nil
			}
			log.Println("error in reading candidate item, error :", err.Error())
			return nil, err
		}

		candidates = append(candidates, candid)
	}

	return candidates, nil
}

func (c *candidate) DeleteCandidate(ctx context.Context, candidateId string) error {
	_, err := DBS.Couch.Query(couchbaseQueries.DeleteCandidateQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println(" error in deleting candidate, error :", err.Error())
		return err
	}
	return nil
}

func (c *candidate) UpdateCandidate(ctx context.Context, candidateData models.Candidate) error {
	_, err := DBS.Couch.Query(couchbaseQueries.UpdateCandidateQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{
			candidateData.Name,
			candidateData.Type,
			candidateData.Descriptions,
			candidateData.ElectionId,
			candidateData.Id,
		},
	})
	if err != nil {
		log.Println(" error in updating candidate, error :", err.Error())
		return err
	}
	return nil
}

func (c *candidate) IsCandidateExists(ctx context.Context, candidateId string) (*bool, error) {
	var exists bool
	result, err := DBS.Couch.Query(couchbaseQueries.IsCandidateExistsQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var count models.CandidatesCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &exists, nil
		}
		log.Println("error in reading candidates count, error :", err.Error())
		return nil, err
	}

	if count.Count > 0 {
		exists = true
	}
	return &exists, nil

}

func (c candidate) GetAllElectionCandidates(ctx context.Context, electionId string) ([]models.Candidate, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetAllElectionCandidatesQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var candidates []models.Candidate
	for result.Next() {
		var candid models.Candidate
		err := result.Row(&candid)
		if err != nil {
			if err == gocb.ErrNoResult {
				return candidates, nil
			}
			log.Println("error in reading candidate item, error :", err.Error())
			return nil, err
		}

		candidates = append(candidates, candid)
	}

	return candidates, nil
}
