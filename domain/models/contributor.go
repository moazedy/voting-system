package models

import (
	"time"

	"github.com/google/uuid"
)

type Contributor struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name,omitempty"`
	ContributeTime time.Time `json:"contribute_time"`
	Deleted        bool      `json:"deleted"`
	DeletedAt      time.Time `json:"deleted_at,omitempty"`
	ElectionId     string    `json:"election_id"`
	VotedAt        time.Time `json:"voted_at"`
}

type ContributorsCount struct {
	Count int `json:"count"`
}
