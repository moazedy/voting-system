package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	Id            uuid.UUID `json:"id"`
	VoteTime      time.Time `json:"vote_time"`
	CandidateId   string    `json:"candidate_id"`
	VoterId       string    `json:"voter_id, omitempty"`
	VoteValue     bool      `json:"vote_value"`
	PrivateVoting bool      `json:"private_voting"`
	Deleted       bool      `json:"deleted,omitempty"`
	Deleted_at    time.Time `json:"deleted_at,omitempty"`
}
