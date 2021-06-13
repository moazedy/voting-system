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
	SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) (*models.ContributionData, error)
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
	GetCandidatePositiveVotes(ctx context.Context, candidateId, requesterId string, pagination *helper.Pagination, requestedByAdmin bool) ([]models.Vote, error)
	// GetCandidateNegativeVotes gets list of negative votes of some candidate
	GetCandidateNegativeVotes(ctx context.Context, candidateId, requesterId string, pagination *helper.Pagination, requestedByAmdin bool) ([]models.Vote, error)
}

// vote is a struct that is way to access vote methods in logic layer
type vote struct {
	repo             repository.VoteRepo
	electionLogic    ElectionLogic
	candidateLogic   CandidateLogic
	contributorLogic ContributorLogic
}

// NewVoteLogic is constractor fucntion of VoteLogic
func NewVoteLogic() VoteLogic {
	return new(vote)
}

func (v vote) SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) (*models.ContributionData, error) {
	// implementation of singleton design pattern
	if v.repo == nil {
		v.repo = repository.NewVoteRepo()
	}
	if v.electionLogic == nil {
		v.electionLogic = NewElectionLogic()
	}
	if v.candidateLogic == nil {
		v.candidateLogic = NewCandidateLogic()
	}
	if v.contributorLogic == nil {
		v.contributorLogic = NewContributorLogic()
	}

	// TODO: validation of vote data should be according to the election type

	// validation of received data as a vote
	if err := voteData.Validate(); err != nil {
		return nil, err
	}

	// reading the election data
	theElection, err := v.electionLogic.ReadElectionData(ctx, voteData.ElectionId)
	if err != nil {
		return nil, err
	}

	// check for election termination
	if theElection.TerminationCheck() {
		return nil, errors.New(constants.ElectionHasTerminated)
	}

	// validation of voteData according to the election type ,, TODO : the method is not being implemented yet
	if err := voteData.VoteValidationAccordingToElectionType(theElection.Type); err != nil {
		return nil, err
	}

	// check on selected conadidate existance
	if _, err := v.candidateLogic.CandidateExistanceCheck(ctx, voteData.CandidateId); err != nil {
		return nil, err
	}

	// this part of code checks if this reuqester has participated in the requested election, could not vote any more
	// because every one can vote just one time in an election
	contributorExists, err := v.contributorLogic.ContributorExists(ctx, requesterId, voteData.ElectionId)
	if err != nil {
		return nil, err
	}
	if *contributorExists {
		return nil, errors.New(constants.ContributorAlredyExists)
	}

	// TODO : checking for contributor access on this specific election if election has it's own access levels
	voteData.Id = uuid.New()
	voteData.ContributorId = requesterId
	voteData.VoteTime = time.Now()

	// saving vote into db
	voteId, err := v.repo.SaveVote(ctx, voteData)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	// creating contributor data, using received data for voting
	contributorData := models.Contributor{
		Name:           voteData.ContributorName,
		ContributeTime: time.Now(),
		ElectionId:     voteData.ElectionId,
		VotedAt:        time.Now(),
	}
	// saving contributor data which has participated in the election with his/shes vote
	contributorId, err := v.contributorLogic.SaveNewContributor(ctx, requesterId, contributorData)
	if err != nil {
		// if contributor could not be saved in db, the given vote should be deleted
		er := v.repo.DeleteVote(ctx, voteId.VoteId)
		if er != nil {
			return nil, errors.New(constants.InternalServerError)
		}
		return nil, err
	}

	return &models.ContributionData{
		ContributorId: contributorId.Id,
		VoteId:        voteId.VoteId,
		ElectionId:    voteData.ElectionId,
	}, nil
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
		return errors.New(constants.InvalidId)
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

func (v vote) GetCandidatePositiveVotes(ctx context.Context, candidateId, requesterId string,
	pagination *helper.Pagination, requestedByAdmin bool) ([]models.Vote, error) {
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

	var votes []models.Vote
	if pagination == nil {
		votes, err = v.repo.GetAllCandidatePositiveVotes(ctx, candidateId)
		if err != nil {
			return nil, errors.New(constants.InternalServerError)
		}
	}

	Pagination := *pagination
	votes, err = v.repo.GetCandidatePositiveVotes(
		ctx,
		candidateId,
		Pagination.GetOrder(),
		Pagination.GetOffset(),
		Pagination.GetLimit(),
	)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}

func (v vote) GetCandidateNegativeVotes(ctx context.Context, candidateId, requesterId string,
	pagination *helper.Pagination, requestedByAdmin bool) ([]models.Vote, error) {
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

	var votes []models.Vote
	if pagination == nil {
		votes, err = v.repo.GetAllCandidateNegativeVotes(ctx, candidateId)
		if err != nil {
			return nil, errors.New(constants.InternalServerError)
		}
	}

	Pagination := *pagination
	votes, err = v.repo.GetCandidateNegativeVotes(
		ctx,
		candidateId,
		Pagination.GetOrder(),
		Pagination.GetOffset(),
		Pagination.GetLimit(),
	)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return votes, nil
}
