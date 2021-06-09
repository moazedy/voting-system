package models

import (
	"errors"
	"time"
	"voting-system/constants"

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

func (c *Contributor) Validate() error {
	if err := c.NameValidation(); err != nil {
		return err
	}

	if err := c.ElectionIdValidation(); err != nil {
		return err
	}

	return nil
}

func (c *Contributor) NameValidation() error {
	if c.Name == "" {
		return errors.New(constants.ContributorNameCanNotBeEmpty)
	}

	if len(c.Name) > constants.MximomContributorNameLength {
		return errors.New(constants.CandidateNameIsLongerThanExpected)
	}
	return nil
}

func (c *Contributor) ElectionIdValidation() error {
	_, err := uuid.Parse(c.ElectionId)
	if err != nil {
		return errors.New(constants.InvalidElectionId)
	}

	return nil
}
