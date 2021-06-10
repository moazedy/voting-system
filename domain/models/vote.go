package models

import (
	"errors"
	"time"
	"voting-system/constants"

	"github.com/google/uuid"
)

type Vote struct {
	Id              uuid.UUID `json:"id"`
	VoteTime        time.Time `json:"vote_time"`
	CandidateId     string    `json:"candidate_id"`
	ContributorId   string    `json:"contributor_id,omitempty"`
	ContributorName string    `json:"contributor_name,omitempty"`
	VoteValue       bool      `json:"vote_value"`
	PrivateVoting   bool      `json:"private_voting"`
	Deleted         bool      `json:"deleted,omitempty"`
	DeletedAt       time.Time `json:"deleted_at,omitempty"`
	ElectionId      string    `json:"election_id"`
}

func (v *Vote) Validate() error {
	if err := v.CandidateIdValidate(); err != nil {
		return err
	}

	if err := v.ContributorIdValidate(); err != nil {
		return err
	}

	if err := v.ElectionIdValidate(); err != nil {
		return err
	}

	if v.Deleted {
		return errors.New(constants.DeletedCanNotBeTrueInTheBegining)
	}

	return nil

}

func (v *Vote) CandidateIdValidate() error {
	if v.CandidateId == "" {
		return errors.New(constants.CandidateIdCanNotBeEmpty)
	}

	_, err := uuid.Parse(v.CandidateId)
	if err != nil {
		return errors.New(constants.InvalidCandidateId)
	}

	return nil
}

func (v *Vote) ContributorIdValidate() error {
	if v.ContributorId == "" {
		return errors.New(constants.ContributorIdCanNotBeEmpty)
	}

	_, err := uuid.Parse(v.ContributorId)
	if err != nil {
		return errors.New(constants.InvalidContributorId)
	}

	return nil
}

func (v *Vote) ElectionIdValidate() error {
	if v.ElectionId == "" {
		return errors.New(constants.ElectionIdCanNotBeEmpty)
	}

	_, err := uuid.Parse(v.ElectionId)
	if err != nil {
		return errors.New(constants.InvalidElectionId)
	}
	return nil
}

type VoteId struct {
	VoteId string `json:"vote_id"`
}

type VoteCount struct {
	Count int `json:"count"`
}

// VoteValidationAccordingToElectionType checks vote values validation according to election
func (v *Vote) VoteValidationAccordingToElectionType(electionType ElectionType) error {
	switch electionType {
	case PublicVotersData:
		// TODO : checking dependencies of this type of election
		return nil
	case PrivateVotersData:
		// TODO : checking dependencies of this type of election
		return nil
	case SimpleYesOrNo:
		// TODO : checking dependencies of this type of election
		return nil
	case PrivateLimitedCount:
		// TODO : checking dependencies of this type of election
		return nil
	case PublicLimitedCount:
		// TODO : checking dependencies of this type of election
		return nil
	default:
		return errors.New(constants.InvalidElectionType)
	}
}
