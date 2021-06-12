package repository

import (
	"context"
	"fmt"
	"log"
	"voting-system/domain/models"
	"voting-system/repository/couchbaseQueries"

	"github.com/couchbase/gocb/v2"
)

// VoteRepo is interface of vote entity in repository layer. other layers of system can interface with it using this
// interface here
type VoteRepo interface {
	// SaveVote saves new vote in database and returns the id
	SaveVote(ctx context.Context, newVote models.Vote) (*models.VoteId, error)
	// ReadSpecificVoteData is to Reading data of a specific vote, stored in database befor
	ReadSpecificVoteData(ctx context.Context, voteId string) (*models.Vote, error)
	// DeleteVote deletes some specific vote using it's id, to not be calculated in results any more
	DeleteVote(ctx context.Context, voteId string) error
	// AgregateOfCandidatePositiveVotes is to reading count of some candidate's positive votes
	AgregateOfCandidatePositiveVotes(ctx context.Context, candidateId string) (*models.CandidateVotesCount, error)
	// AgregateOfCandidateNegativeVotes is to reading count of some candidate's negative votes
	AgregateOfCandidateNegativeVotes(ctx context.Context, candidateId string) (*models.CandidateVotesCount, error)
	// UpdateVoteDat aupdates some vote's data
	UpdateVoteData(ctx context.Context, newVoteData models.Vote) error
	// GetCandidateVotes gets all of given candidate votes in system
	GetCandidateVotes(ctx context.Context, candidateId, order string, offset, limit int) ([]models.Vote, error)
	// GetCandidatePositiveVotes gets list of  positive votes of candidate
	GetCandidatePositiveVotes(ctx context.Context, candidateId, order string, offset, limit int) ([]models.Vote, error)
	// GetCandidateNegativeVotes gets list of  negative votes of candidate
	GetCandidateNegativeVotes(ctx context.Context, candidateId, order string, offset, limit int) ([]models.Vote, error)
	// GetAllCandidatePositiveVotes gets list of all positive votes of candidate
	GetAllCandidatePositiveVotes(ctx context.Context, candidateId string) ([]models.Vote, error)
	// GetAllCandidateNegativeVotes gets list of all negative votes of candidate
	GetAllCandidateNegativeVotes(ctx context.Context, candidateId string) ([]models.Vote, error)
	// VoteExists checks on some specific vote existance
	VoteExists(ctx context.Context, voteId string) (*bool, error)
}

// vote is a struct that represents vote entity in repository layer and its the way we can access to repository methods of
// vote in this layer
type vote struct {
}

// NewVoteRepo is Constructor function of vote entity in repositroy layer
func NewVoteRepo() VoteRepo {
	return new(vote)
}

func (v *vote) SaveVote(ctx context.Context, newVote models.Vote) (*models.VoteId, error) {
	_, err := DBS.Couch.Query(couchbaseQueries.SaveVoteQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newVote.Id, newVote},
	})
	if err != nil {
		log.Println(" error in taking new vote, error :", err.Error())
		return nil, err
	}
	return &models.VoteId{
		VoteId: newVote.Id.String(),
	}, nil
}

func (v *vote) ReadSpecificVoteData(ctx context.Context, voteId string) (*models.Vote, error) {
	var myVote models.Vote
	result, err := DBS.Couch.Query(couchbaseQueries.ReadVoteQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{voteId},
	})
	if err != nil {
		log.Println("error in reading vote data, error :", err.Error())
		return nil, err
	}
	err = result.One(&myVote)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &myVote, nil
		}
		log.Println("error in reading vote item, error :", err.Error())
		return nil, err
	}

	return &myVote, nil
}

func (v *vote) DeleteVote(ctx context.Context, voteId string) error {
	_, err := DBS.Couch.Query(couchbaseQueries.DeleteVoteQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{voteId},
	})
	if err != nil {
		log.Println(" error in deleting vote, error :", err.Error())
		return err
	}
	return nil
}

func (v *vote) AgregateOfCandidatePositiveVotes(ctx context.Context, candidateId string) (*models.CandidateVotesCount, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetCandidatePositiveVotesCount, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var count models.CandidateVotesCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &count, nil
		}
		log.Println("error in reading votes count, error :", err.Error())
		return nil, err
	}

	return &count, nil
}

func (v *vote) AgregateOfCandidateNegativeVotes(ctx context.Context, candidateId string) (*models.CandidateVotesCount, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetCandidateNegativeVotesCount, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var count models.CandidateVotesCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &count, nil
		}
		log.Println("error in reading votes count, error :", err.Error())
		return nil, err
	}

	return &count, nil
}

func (v *vote) UpdateVoteData(ctx context.Context, newVoteData models.Vote) error {

	_, err := DBS.Couch.Query(couchbaseQueries.UpdateVoteQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{
			newVoteData.CandidateId,
			newVoteData.ContributorId,
			newVoteData.VoteValue,
			newVoteData.PrivateVoting,
			newVoteData.ElectionId,
			newVoteData.Id,
		},
	})
	if err != nil {
		log.Println(" error in updating vote, error :", err.Error())
		return err
	}
	return nil
}

func (v *vote) GetCandidateVotes(ctx context.Context, candidateId, order string, offset, limit int) ([]models.Vote, error) {
	query := fmt.Sprintf(couchbaseQueries.GetCandidateVotesQuery, order)
	result, err := DBS.Couch.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId, offset, limit},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var votes []models.Vote
	for result.Next() {
		var aVote models.Vote
		err := result.Row(&aVote)
		if err != nil {
			if err == gocb.ErrNoResult {
				return votes, nil
			}
			log.Println("error in reading vote item, error :", err.Error())
			return nil, err
		}

		votes = append(votes, aVote)
	}

	return votes, nil
}

func (v *vote) GetCandidatePositiveVotes(ctx context.Context, candidateId, order string, offset, limit int) ([]models.Vote, error) {
	query := fmt.Sprintf(couchbaseQueries.GetCandidatePositiveVotesQuery, order)
	result, err := DBS.Couch.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId, offset, limit},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var votes []models.Vote
	for result.Next() {
		var aVote models.Vote
		err := result.Row(&aVote)
		if err != nil {
			if err == gocb.ErrNoResult {
				return votes, nil
			}
			log.Println("error in reading vote item, error :", err.Error())
			return nil, err
		}

		votes = append(votes, aVote)
	}

	return votes, nil
}

func (v *vote) GetCandidateNegativeVotes(ctx context.Context, candidateId, order string, offset, limit int) ([]models.Vote, error) {
	query := fmt.Sprintf(couchbaseQueries.GetCandidateNegativeVotesQuery, order)
	result, err := DBS.Couch.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId, offset, limit},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var votes []models.Vote
	for result.Next() {
		var aVote models.Vote
		err := result.Row(&aVote)
		if err != nil {
			if err == gocb.ErrNoResult {
				return votes, nil
			}
			log.Println("error in reading vote item, error :", err.Error())
			return nil, err
		}

		votes = append(votes, aVote)
	}

	return votes, nil
}

func (v vote) VoteExists(ctx context.Context, voteId string) (*bool, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.VoteExistsQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{voteId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var exists bool
	var count models.VoteCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &exists, nil
		}
		log.Println("error in reading vote count, error :", err.Error())
		return nil, err
	}

	if count.Count > 0 {
		exists = true
		return &exists, nil
	}

	return &exists, nil

}

func (v vote) GetAllCandidatePositiveVotes(ctx context.Context, candidateId string) ([]models.Vote, error) {

	result, err := DBS.Couch.Query(couchbaseQueries.GetAllCandidatePositiveVotesQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var votes []models.Vote
	for result.Next() {
		var aVote models.Vote
		err := result.Row(&aVote)
		if err != nil {
			if err == gocb.ErrNoResult {
				return votes, nil
			}
			log.Println("error in reading vote item, error :", err.Error())
			return nil, err
		}

		votes = append(votes, aVote)
	}

	return votes, nil
}

func (v vote) GetAllCandidateNegativeVotes(ctx context.Context, candidateId string) ([]models.Vote, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetAllCandidateNegativeVotesQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{candidateId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var votes []models.Vote
	for result.Next() {
		var aVote models.Vote
		err := result.Row(&aVote)
		if err != nil {
			if err == gocb.ErrNoResult {
				return votes, nil
			}
			log.Println("error in reading vote item, error :", err.Error())
			return nil, err
		}

		votes = append(votes, aVote)
	}

	return votes, nil
}
