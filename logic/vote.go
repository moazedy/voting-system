package logic

import (
	"context"
	"errors"
	"time"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/helper"
	"voting-system/repository"

	"github.com/google/uuid"
)

// VoteLogic is interface of vote entity in logic layer of system
type VoteLogic interface {
	// SaveNewVote saves new vote data and returns id of the vote
	SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) (*models.VoteId, error)
	// ReadVoteData returns data of given voteId
	ReadVoteData(ctx context.Context, voteId, requesterId string, requestedByAdmin bool) (*models.Vote, error)
	// DeleteVote deletes requested voteId and returns an error if any problem happens during the operation
	DeleteVote(ctx context.Context, voteId, requesterId string, requestedByAdmin bool) error
	// AgregateOfCandidatePositiveVotes is to reading count of some candidate's positive votes
	AgregateOfCandidatePositiveVotes(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) (*models.CandidateVotesCount, error)
	// AgregateOfCandidateNegativeVotes is to reading count of some candidate's negative votes
	AgregateOfCandidateNegativeVotes(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) (*models.CandidateVotesCount, error)
	// UpdateVoteData updates data of some specific vote
	UpdateVoteData(ctx context.Context, voteId, requesterId string, voteData models.Vote, requestedByAdmin bool) error
	// GetCandidateVotes gets all of a specific candidate votes
	GetCandidateVotes(ctx context.Context, candidateId, requesterId string, pagination helper.Pagination, requestedByAdmin bool) ([]models.Vote, error)
	// GetCandidatePositiveVotes gets list of positive votes of some candidate
	GetCandidatePositiveVotes(ctx context.Context, candidateId, requesterId string, pagination helper.Pagination, requestedByAdmin bool) ([]models.Vote, error)
	// GetCandidateNegativeVotes gets list of negative votes of some candidate
	GetCandidateNegativeVotes(ctx context.Context, candidateId, requesterId string, pagination helper.Pagination, requestedByAmdin bool) ([]models.Vote, error)
}

// vote is a struct that is way to access vote methods in logic layer
type vote struct {
	repo           repository.VoteRepo
	electionLogic  ElectionLogic
	candidateLogic CandidateLogic
}

// NewVoteLogic is constractor fucntion of VoteLogic
func NewVoteLogic() VoteLogic {
	return new(vote)
}

func (v *vote) SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) (*models.VoteId, error) {
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}
	if v.electionLogic == nil {
		v.electionLogic = NewElectionLogic()
	}
	if v.candidateLogic == nil {
		v.candidateLogic = NewCandidateLogic()
	}

	if err := voteData.Validate(); err != nil {
		return nil, err
	}

	if _, err := v.electionLogic.CheckElectionExistance(ctx, voteData.ElectionId); err != nil {
		return nil, err
	}

	if _, err := v.candidateLogic.CandidateExistanceCheck(ctx, voteData.CandidateId); err != nil {
		return nil, err
	}

	// TODO : checking for contributor access on this specific election
	voteData.Id = uuid.New()
	voteData.ContributorId = requesterId
	voteData.VoteTime = time.Now()

	id, err := v.repo.SaveVote(ctx, voteData)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return id, nil
}

func (v vote) ReadVoteData(ctx context.Context, voteId, requesterId string, requestedByAdmin bool) (*models.Vote, error) {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}

	// check for vote existance
	exists, err := v.repo.VoteExists(ctx, voteId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	if !*exists {
		return nil, errors.New(constants.VoteDoesNotExist)
	}

	// reading the vote
	theVote, err := v.repo.ReadSpecificVoteData(ctx, voteId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	// access checking
	if !requestedByAdmin {
		if requesterId != theVote.ContributorId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	return theVote, nil
}

func (v vote) DeleteVote(ctx context.Context, voteId, requesterId string, requestedByAdmin bool) error {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}

	// reading the vote
	theVote, err := v.repo.ReadSpecificVoteData(ctx, voteId)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	// access checking
	if !requestedByAdmin {
		if requesterId != theVote.ContributorId {
			return errors.New(constants.AccessDenied)
		}
	}

	err = v.repo.DeleteVote(ctx, voteId)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (v vote) AgregateOfCandidatePositiveVotes(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) (*models.CandidateVotesCount, error) {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}
	if v.candidateLogic == nil {
		v.candidateLogic = NewCandidateLogic()
	}

	// calling ReadingCandidateData on some candidateId, without considering returning data, checks on all validations and
	// access levels
	_, err := v.candidateLogic.ReadCandidateData(ctx, candidateId, requesterId, requestedByAdmin)
	if err != nil {
		return nil, err
	}

	votes, err := v.repo.AgregateOfCandidatePositiveVotes(ctx, candidateId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}

func (v vote) AgregateOfCandidateNegativeVotes(ctx context.Context, candidateId, requesterId string, requestedByAdmin bool) (*models.CandidateVotesCount, error) {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}
	if v.candidateLogic == nil {
		v.candidateLogic = NewCandidateLogic()
	}

	// calling ReadingCandidateData on some candidateId, without considering returning data, checks on all validations and
	// access levels
	_, err := v.candidateLogic.ReadCandidateData(ctx, candidateId, requesterId, requestedByAdmin)
	if err != nil {
		return nil, err
	}

	votes, err := v.repo.AgregateOfCandidateNegativeVotes(ctx, candidateId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}

func (v vote) UpdateVoteData(ctx context.Context, voteId, requesterId string, voteData models.Vote, requestedByAdmin bool) error {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}

	if err := IdValidation(voteId); err != nil {
		return err
	}

	if err := voteData.Validate(); err != nil {
		return err
	}

	uid, err := uuid.Parse(voteId)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}
	voteData.Id = uid

	err = v.repo.UpdateVoteData(ctx, voteData)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (v vote) GetCandidateVotes(ctx context.Context, candidateId, requesterId string, pagination helper.Pagination, requestedByAdmin bool) ([]models.Vote, error) {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}

	// calling ReadingCandidateData on some candidateId, without considering returning data, checks on all validations and
	// access levels
	_, err := v.candidateLogic.ReadCandidateData(ctx, candidateId, requesterId, requestedByAdmin)
	if err != nil {
		return nil, err
	}

	votes, err := v.repo.GetCandidateVotes(
		ctx,
		candidateId,
		pagination.GetOrder(),
		pagination.GetOffset(),
		pagination.GetLimit(),
	)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}

func (v vote) GetCandidatePositiveVotes(ctx context.Context, candidateId, requesterId string, pagination helper.Pagination, requestedByAdmin bool) ([]models.Vote, error) {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}

	// calling ReadingCandidateData on some candidateId, without considering returning data, checks on all validations and
	// access levels
	_, err := v.candidateLogic.ReadCandidateData(ctx, candidateId, requesterId, requestedByAdmin)
	if err != nil {
		return nil, err
	}

	votes, err := v.repo.GetCandidatePositiveVotes(
		ctx,
		candidateId,
		pagination.GetOrder(),
		pagination.GetOffset(),
		pagination.GetLimit(),
	)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}

func (v vote) GetCandidateNegativeVotes(ctx context.Context, candidateId, requesterId string, pagination helper.Pagination, requestedByAdmin bool) ([]models.Vote, error) {
	// singlton design pattern ...
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}

	// calling ReadingCandidateData on some candidateId, without considering returning data, checks on all validations and
	// access levels
	_, err := v.candidateLogic.ReadCandidateData(ctx, candidateId, requesterId, requestedByAdmin)
	if err != nil {
		return nil, err
	}

	votes, err := v.repo.GetCandidateNegativeVotes(
		ctx,
		candidateId,
		pagination.GetOrder(),
		pagination.GetOffset(),
		pagination.GetLimit(),
	)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}
