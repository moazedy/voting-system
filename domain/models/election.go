package models

import (
	"errors"
	"time"
	"voting-system/constants"

	"github.com/google/uuid"
)

type Election struct {
	// Id is election unique id that every election could be known by it
	Id uuid.UUID `json:"id"`
	// Title is the title that creator chose for election
	Title string `json:"title"`
	// CreationTime is time of the moment wich election created in it
	CreationTime time.Time `json:"creation_time"`
	// StartTime is the time that creator sets as time of election begining
	StartTime time.Time `json:"start_time"`
	// EndTime is the time that creator sets as election termination
	EndTime time.Time `json:"end_time"`
	// HasEnded determines that the election has ended or not
	HasEnded bool `json:"has_ended"`
	// Type specifies type of the election
	Type ElectionType `json"type"`
	// CandidatesCountLimit determines the maximom of candidates in the election
	CandidatesCountLimit int `json:"candidate_count_limit"`
	// CreatorId is id of the user who created the election
	CreatorId string `json:"creator_id"`
	// CountOfLimitedCandidates is a field that be evaluate in LimitedCount types of elections, this map is
	// created from candidate id as string key and maximom picking time of it as int value of the key
	CountOfLimitedCandidates map[string]int `json:"count_of_limited_candidates,omitempty"`
	// RelatedPersons holds list of persons related to the election
	RelatedPersons []RelatedPerson `json:"related_persons"`
	// RelatedCategories holds list of categories related to the election
	RelatedCategories []RelatedCategory `json:"related_categories,omitempty"`
}

type ElectionType int

const (
	// PublicVotersData : in this type of election, votes data is public and it is clear that who votes to who/what
	PublicVotersData ElectionType = iota
	// PrivateVotesData : in this type of election, votes data is private and results only accessable
	PrivateVotersData
	// SimpleYesOrNo : is a simple questioning about a promblem by vote values as yes or no
	SimpleYesOrNo
	// PrivateLimitedCout : this type of election being used when picking of limited count objects or problems is
	// considerd. it means that every introduced candidate has its own limited count to be picked, like a basket of
	// different.
	PrivateLimitedCount
	// PublicLimitedCount : it is exactly like PrivateLimitedCount but in this type, voters data is public
	PublicLimitedCount
)

// Validate validates data inside of the election that it's being called on
func (e *Election) Validate() error {
	if err := e.TitleValidate(); err != nil {
		return err
	}

	if err := e.TypeValidate(); err != nil {
		return nil
	}

	if err := e.TimesValidate(); err != nil {
		return err
	}

	if err := e.CreatorIdValidate(); err != nil {
		return err
	}

	if err := e.CandidatesCountValidate(); err != nil {
		return err
	}

	return nil
}

func (e *Election) TitleValidate() error {
	if e.Title == "" {
		return errors.New(constants.TitleCanNotBeEmpty)
	}

	if len(e.Title) > constants.MaximomTitleLength {
		return errors.New(constants.TitleLengthIsMoreThanMaximom)
	}
	// TODO : adding more validation statements
	return nil
}

func (e *Election) TypeValidate() error {
	if !((e.Type == PublicVotersData) || (e.Type == PrivateVotersData) || (e.Type == SimpleYesOrNo) ||
		(e.Type == PrivateLimitedCount) || (e.Type == PublicLimitedCount)) {
		return errors.New(constants.InvalidElectionType)
	}
	return nil
}

func (e *Election) TimesValidate() error {
	if e.EndTime.Before(e.StartTime) {
		return errors.New(constants.EndTimeCanNotBeBeforStartTime)
	}

	if e.StartTime.Before(e.CreationTime) {
		return errors.New(constants.StartTimeCanNotBeBeforCreationTime)
	}

	if e.HasEnded {
		return errors.New(constants.ElectionCanNotBeEndedBeforItBegins)
	}

	return nil
}

func (e *Election) CreatorIdValidate() error {
	if e.CreatorId == "" {
		return errors.New(constants.ElectionCreatorIdCanNotBeEmpty)
	}

	_, err := uuid.Parse(e.CreatorId)
	if err != nil {
		return errors.New(constants.InvalidCreatorId)
	}

	return nil
}

func (e *Election) CandidatesCountValidate() error {
	if (e.CandidatesCountLimit <= 0) || (e.CandidatesCountLimit > constants.MaximomCandidatesCount) {
		return errors.New(constants.InvalidCandidatesCountLimit)
	}

	return nil
}

type ElectionsCount struct {
	Count int `json:"count"`
}

/*type RelatedUsers struct {
	Users []string `json:"users"`
}*/

/*type RelatedCategories struct {
	Categories []string `json:"categories"`
}*/

type ElectionId struct {
	ElectionId string `json:"election_id"`
}

type Elections struct {
	Elections []ElectionId `json:"elections"`
}

type RelatedPerson struct {
	Id       string `json:"id"`
	Username string `json:"username,omitempty"`
}

type RelatedCategory struct {
	Id  string `json:"id"`
	Tag string `json:"tag"`
}

type ElectionResults struct {
	Id         uuid.UUID                 `json:"id"`
	ElectionId string                    `json:"election_id"`
	Title      string                    `json:"title"`
	Type       ElectionType              `json:"type"`
	HasEnded   bool                      `json:"has_ended"`
	Results    []CandidateElectionResult `json:"results"`
}

type CandidateElectionResult struct {
	CandidateId        string `json:"cnadidate_id"`
	CandidateName      string `json:"candidate_name"`
	PositiveVotesCount int    `json:"positive_votes_count"`
	NegativeVotesCount int    `json:"negative_votes_count,omitempty"`
}
