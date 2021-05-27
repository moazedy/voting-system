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

type VoteLogic interface {
	SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) error
}

type vote struct {
	repo           repository.VoteRepo
	electionLogic  ElectionLogic
	candidateLogic CandidateLogic
}

func NewVoteLogic() VoteLogic {
	return new(vote)
}

func (v *vote) SaveNewVote(ctx context.Context, voteData models.Vote, requesterId string) error {
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
		return err
	}

	if _, err := v.electionLogic.CheckElectionExistance(ctx, voteData.ElectionId); err != nil {
		return err
	}

	if _, err := v.candidateLogic.CandidateExistanceCheck(ctx, voteData.CandidateId); err != nil {
		return err
	}

	// TODO : checking for contributor access on this specific election
	voteData.Id = uuid.New()
	voteData.ContributorId = requesterId
	voteData.VoteTime = time.Now()

	err := v.repo.SaveVote(ctx, voteData)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}
