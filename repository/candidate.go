package repository

import (
	"context"
	"voting-system/domain/models"
)

type CandidateRepo interface {
	// CreateCandidate is for creating new candidate entitiy in voting system
	CreateCandidate(ctx context.Context, NewCandidate models.Candidate) error
	// ReadCandidate reads data of requested candidate id, if it exists in db
	ReadCandidate(ctx context.Context, cadidateId string) (*models.Candidate, error)

}
