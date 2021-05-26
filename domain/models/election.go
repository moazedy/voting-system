package models

import (
	"errors"
	"time"
	"voting-system/constants"

	"github.com/google/uuid"
)

type Election struct {
	Id                   uuid.UUID    `json:"id"`
	Title                string       `json:"title"`
	CreationTime         time.Time    `json:"creation_time"`
	StartTime            time.Time    `json:"start_time"`
	EndTime              time.Time    `json:"end_time"`
	HasEnded             bool         `json:"has_ended"`
	Type                 ElectionType `json"type"`
	CandidatesCountLimit int          `json:"candidate_count_limit"`
	CreatorId            string       `json:"creator_id"`
}

type ElectionType int

const (
	PublicVotersData ElectionType = iota
	PrivateVotersData
	SimpleYesOrNo
)

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
	if !((e.Type == PublicVotersData) || (e.Type == PrivateVotersData) || (e.Type == SimpleYesOrNo)) {
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
