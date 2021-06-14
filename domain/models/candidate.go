package models

import (
	"errors"
	"time"
	"voting-system/constants"

	"github.com/google/uuid"
)

type Candidate struct {
	MetaId       uuid.UUID     `json:"meta_id,omitempty"`
	CandidateId  string        `json:"candidate_id"`
	Name         string        `json:"name"`
	Type         CandidateType `json:"type"`
	Descriptions []string      `json:"descriptions"`
	Created_At   time.Time     `json:"created_at"`
	Deleted      bool          `json:"deleted,omitempty"`
	Deleted_at   time.Time     `json:"deleted_at,omitempty"`
	ElectionId   string        `json:"election_id"`
}

type CandidateType int

const (
	Person CandidateType = iota
	Problem
)

type CandidateVotesCount struct {
	Count int `json:"count"`
}

type CandidatesCount struct {
	Count int `json:"count"`
}

func (c *Candidate) Validate() error {
	if err := c.NameValidate(); err != nil {
		return err
	}

	if err := c.TypeValidate(); err != nil {
		return err
	}

	if err := c.ElectionIdValidate(); err != nil {
		return err
	}

	if c.Deleted {
		return errors.New(constants.DeletedCanNotBeTrueInTheBegining)
	}

	return nil
}

func (c *Candidate) NameValidate() error {

	if c.Name == "" {
		return errors.New(constants.CandidateNameCanNotBeEmpty)
	}

	if len(c.Name) > constants.MaximomCandidateNameLength {
		return errors.New(constants.CandidateNameIsLongerThanExpected)
	}
	// TODO adding more name validations

	return nil
}

func (c *Candidate) TypeValidate() error {
	if !((c.Type == Person) || (c.Type == Problem)) {
		return errors.New(constants.InvalidCandidateType)
	}
	return nil
}

func (c *Candidate) ElectionIdValidate() error {
	if c.ElectionId == "" {
		return errors.New(constants.ElectionIdCanNotBeEmpty)
	}

	_, err := uuid.Parse(c.ElectionId)
	if err != nil {
		return errors.New(constants.InvalidElectionId)
	}
	return nil
}
