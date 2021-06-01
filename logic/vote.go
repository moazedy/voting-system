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

// VoteLogic is interface of vote entity in logic layer of system
type VoteLogic interface {
	// SaveNewVote saves new vote data and returns id of the vote
	SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) (*models.VoteId, error)
	// ReadVoteData returns data of given voteId
	ReadVoteData(ctx context.Context, voteId, requesterId string, requestedByAdmin bool) (*models.Vote, error)
	// DeleteVote deletes requested voteId and returns an error if any problem happens during the operation
	DeleteVote(ctx context.Context, voteId, requesterId string, requestedByAdmin bool) error
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
