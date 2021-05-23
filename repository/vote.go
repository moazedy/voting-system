package repository

import (
	"context"
	"voting-system/domain/models"
)

// VoteRepo is interface of vote entity in repository layer. other layers of system can interface with it using this
// interface here
type VoteRepo interface {
	// Poll saves new vote in database
	Poll(ctx context.Context, newVote models.Vote) error
	// ReadSpecificVoteData is to Reading data of a specific vote, stored in database befor
	ReadSpecificVoteData(ctx context.Context, voteId string) (*models.Vote, error)
	// DeleteVote deletes some specific vote using it's id, to not be calculated in results any more
	DeleteVote(ctx context.Context, voteId string) error
	// AgregateOfCandidatePositiveVotes is to reading count of some candidate's positive votes
	AgregateOfCandidatePositiveVotes(ctx context.Context, candidateId string) (*int, error)
	// AgregateOfCandidateNegativeVotes is to reading count of some candidate's negative votes
	AgregateOfCandidateNegativeVotes(ctx context.Context, candidateId string) (*int, error)
	// UpdateVoteDat aupdates some vote's data
	UpdateVoteData(ctx context.Context, newVoteData models.Vote) error
}

// vote is a struct that represents vote entity in repository layer and its the way we can access to repository methods of
// vote in this layer
type vote struct {
}

// NewVoteRepo is Constructor function of vote entity in repositroy layer
func NewVoteRepo() VoteRepo {
	return new(vote)
}

func (v *vote) Poll(ctx context.Context, newVote models.Vote) error {
	// TODO
	return nil
}

func (v *vote) ReadSpecificVoteData(ctx context.Context, voteId string) (*models.Vote, error) {

	// TODO
	return nil, nil
}

func (v *vote) DeleteVote(ctx context.Context, voteId string) error {

	// TODO
	return nil
}

func (v *vote) AgregateOfCandidatePositiveVotes(ctx context.Context, candidateId string) (*int, error) {

	// TODO
	return nil, nil
}

func (v *vote) AgregateOfCandidateNegativeVotes(ctx context.Context, candidateId string) (*int, error) {

	// TODO
	return nil, nil
}

func (v *vote) UpdateVoteData(ctx context.Context, newVoteData models.Vote) error {

	// TODO
	return nil
}
