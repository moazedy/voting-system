package models

import (
	"time"

	"github.com/google/uuid"
)

type Candidate struct {
	Id           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Type         CandidateType `json:"type"`
	Descriptions []string      `json:"descriptions"`
	Created_At   time.Time     `json:"created_at"`
	Deleted      bool          `json:"deleted,omitempty"`
	Deleted_at   time.Time     `json:"deleted_at,omitempty"`
}

type CandidateType int

const (
	Person CandidateType = iota
	Problem
)
