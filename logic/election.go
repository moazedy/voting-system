package logic

import (
	"context"
	"errors"
	"time"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/repository"

	"github.com/google/uuid"
)

// ElectionLogic is election entity interface in logic layer and other layers of program can intract with it
// through this interface here
type ElectionLogic interface {
	// CreateNewElection craates new election using receieved election data
	CreateNewElection(ctx context.Context, userId string, electionData models.Election) (*models.Id, error)
	// ReadElectionData reads data of given election id
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	// DeleteElection deletes the election with received election id
	DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error
	// GetElectionContributorsCount gets number of persons who contributed in given election id
	GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error)
	// UpdateElection updates data of some election using received electionData
	UpdateElection(ctx context.Context, electionData models.Election, requesterId string, requestedByAdmin bool) error
	// CheckElectionExistance checks on election id existance in db
	CheckElectionExistance(ctx context.Context, electionId string) (*bool, error)
	// GetListOfRelatedUsers gets list of user Ids being added to the election as related users
	GetListOfRelatedUsers(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) ([]models.RelatedPerson, error)
	// GetListOfRelatedCategories gets list of category Ids wich being added to the election as related categories
	GetListOfRelatedCategories(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) ([]models.RelatedCategory, error)
	// GetUserRelatedElections gets list of election Ids wich are related to a specific user
	GetUserRelatedElections(ctx context.Context, userId, requesterId string, requestedByAdmin bool) (*models.Elections, error)
	// GetCategoryRelatedElections gets list of elections that are related to given categoryId
	GetCategoryRelatedElections(ctx context.Context, categoryId, requesterId string, requestedByAdmin bool) (*models.Elections, error)
	// ConcurrentCalculationElectionResults gets results of an election and stores it in db (using concurrency in process)
	ConcurrentCalculationElectionResults(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*models.ElectionResults, map[string]error)
	// CalculationElectionResults gets results of an election and stores it in db
	CalculationElectionResults(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*models.ElectionResults, error)
	// ChangeElectionTerminationStatus changes state of termination in given election
	ChangeElectionTerminationStatus(ctx context.Context, electionId, requesterId string, status, requestedByAdmin bool) error
	// GetListOfStartedElections gets list of not ended elections
	GetListOfStartedElections(ctx context.Context, requestedByAdmin bool) ([]models.Election, error)
	// ElectionResultExists checks for existance of a specific election results
	ElectionResultExists(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*bool, error)
}

// election struct, is holder of election metods in logic layer
type election struct {
	// repo is the way this layer of program can interface with repository layer
	repo           repository.ElectionRepo
	candidateLogic CandidateLogic
	voteLogic      VoteLogic
}

// NewElectionLogic is constractor function of ElectionLogic
func NewElectionLogic() ElectionLogic {
	return new(election)
}

func (e *election) CreateNewElection(ctx context.Context, userId string, electionData models.Election) (*models.Id, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	electionData.Id = uuid.New()
	electionData.CreatorId = userId
	electionData.CreationTime = time.Now()

	if err := electionData.Validate(); err != nil {
		return nil, err
	}

	id, err := e.repo.SaveNewElection(ctx, electionData)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return id, nil
}

func (e *election) ReadElectionData(ctx context.Context, electionId string) (*models.Election, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if err := IdValidation(electionId); err != nil {
		return nil, err
	}

	_, err := e.CheckElectionExistance(ctx, electionId)
	if err != nil {
		return nil, err
	}

	wantedElection, err := e.repo.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return wantedElection, nil
}

func (e *election) DeleteElection(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) error {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return errors.New(constants.AccessDenied)
		}
	}

	if err := e.repo.DeleteElection(ctx, electionId); err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (e *election) GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if err := IdValidation(electionId); err != nil {
		return nil, err
	}

	count, err := e.repo.GetElectionContributorsCount(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return count, nil
}

func (e *election) UpdateElection(ctx context.Context, electionData models.Election, requesterId string, requestedByAdmin bool) error {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	theElection, err := e.ReadElectionData(ctx, electionData.Id.String())
	if err != nil {
		return err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return errors.New(constants.AccessDenied)
		}
	}

	if err := electionData.Validate(); err != nil {
		return err
	}

	err = e.repo.UpdateElection(ctx, electionData)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (e *election) CheckElectionExistance(ctx context.Context, electionId string) (*bool, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if err := IdValidation(electionId); err != nil {
		return nil, err
	}

	exists, err := e.repo.ElectionExists(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	if !*exists {
		return exists, errors.New(constants.ElectionDoesNotExist)
	}

	return exists, nil

}

func (e *election) GetListOfRelatedUsers(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) ([]models.RelatedPerson, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}
	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	users, err := e.repo.GetListOfRelatedUsers(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return users, nil
}

func (e *election) GetListOfRelatedCategories(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) ([]models.RelatedCategory, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}
	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, err
	}

	if !requestedByAdmin {
		// TODO : more access checking
		if requesterId != theElection.CreatorId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	cats, err := e.repo.GetListOfRelatedCategories(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return cats, nil
}

func (e *election) GetUserRelatedElections(ctx context.Context, userId, requesterId string, requestedByAdmin bool) (*models.Elections, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if !requestedByAdmin {
		// TODO more access checks ...
		if userId != requesterId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	elections, err := e.repo.GetUserRelatedElections(ctx, userId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return elections, nil
}

func (e election) GetCategoryRelatedElections(ctx context.Context, categoryId, requesterId string, requestedByAdmin bool) (*models.Elections, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	if !requestedByAdmin {
		// TODO : checking requester accesses on category
	}

	elections, err := e.repo.GetCategoryRelatedElections(ctx, categoryId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return elections, nil
}

func (e election) ConcurrentCalculationElectionResults(ctx context.Context, electionId, requesterId string,
	requestedByAdmin bool) (*models.ElectionResults, map[string]error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}
	if e.candidateLogic == nil {
		e.candidateLogic = NewCandidateLogic()
	}
	if e.voteLogic == nil {
		e.voteLogic = NewVoteLogic()
	}

	// returning errors map
	Errors := make(map[string]error)

	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		Errors[constants.ReadingElectionError] = err
		return nil, Errors
	}

	if !requestedByAdmin {
		if theElection.CreatorId != requesterId {
			Errors[constants.AccessError] = errors.New(constants.AccessDenied)
			return nil, Errors
		}
	}

	allCandidates, err := e.candidateLogic.GetAllElectionCandidates(ctx, electionId, requesterId, requestedByAdmin)
	if err != nil {
		Errors[constants.ReadingCandidatesError] = err
		return nil, Errors
	}

	results := make([]models.CandidateElectionResult, len(allCandidates))
	for k, v := range allCandidates {
		results[k].CandidateId = v.CandidateId
		results[k].CandidateMetaId = v.MetaId.String()
		results[k].CandidateName = v.Name

		go func() {
			positiveVotes, err := e.voteLogic.AgregateOfCandidatePositiveVotes(ctx, v.MetaId.String(), requesterId, requestedByAdmin)
			if err != nil {
				Errors[constants.PositiveVotes+v.MetaId.String()] = err
				return
			}
			results[k].PositiveVotesCount = positiveVotes.Count
		}()

		go func() {
			negativeVotes, err := e.voteLogic.AgregateOfCandidateNegativeVotes(ctx, v.MetaId.String(), requesterId, requestedByAdmin)
			if err != nil {
				Errors[constants.NegativeVotes+v.MetaId.String()] = err
				return
			}
			results[k].PositiveVotesCount = negativeVotes.Count
		}()
	}

	result := models.ElectionResults{
		Id:         uuid.New(),
		ElectionId: electionId,
		Title:      theElection.Title,
		Type:       theElection.Type,
		HasEnded:   theElection.HasEnded,
		Results:    results,
	}

	_, err = e.repo.SaveElectionResult(ctx, result)
	if err != nil {
		Errors[constants.SavingElectionResultsError] = err
		return nil, Errors
	}

	return &result, Errors
}

func (e election) CalculationElectionResults(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*models.ElectionResults, error) {
	// this part of code follows singlton design pattern
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}
	if e.candidateLogic == nil {
		e.candidateLogic = NewCandidateLogic()
	}
	if e.voteLogic == nil {
		e.voteLogic = NewVoteLogic()
	}

	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, err
	}

	if !requestedByAdmin {
		if theElection.CreatorId != requesterId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	allCandidates, err := e.candidateLogic.GetAllElectionCandidates(ctx, electionId, requesterId, requestedByAdmin)
	if err != nil {
		return nil, err
	}

	results := make([]models.CandidateElectionResult, len(allCandidates))
	for k, v := range allCandidates {
		results[k].CandidateId = v.CandidateId
		results[k].CandidateMetaId = v.MetaId.String()
		results[k].CandidateName = v.Name

		positiveVotes, err := e.voteLogic.AgregateOfCandidatePositiveVotes(ctx, v.MetaId.String(), requesterId, requestedByAdmin)
		if err != nil {
			return nil, err
		}

		negativeVotes, err := e.voteLogic.AgregateOfCandidateNegativeVotes(ctx, v.MetaId.String(), requesterId, requestedByAdmin)
		if err != nil {
			return nil, err
		}

		results[k].PositiveVotesCount = positiveVotes.Count
		results[k].NegativeVotesCount = negativeVotes.Count
	}

	result := models.ElectionResults{
		Id:         uuid.New(),
		ElectionId: electionId,
		Title:      theElection.Title,
		Type:       theElection.Type,
		HasEnded:   theElection.HasEnded,
		Results:    results,
	}

	_, err = e.repo.SaveElectionResult(ctx, result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (e election) ChangeElectionTerminationStatus(ctx context.Context, electionId, requesterId string, status, requestedByAdmin bool) error {
	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return err
	}

	if !requestedByAdmin {
		if requesterId != theElection.CreatorId {
			return errors.New(constants.AccessDenied)
		}
	}

	err = e.repo.ChangeElectionTerminationStatus(ctx, electionId, status)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}

func (e election) GetListOfStartedElections(ctx context.Context, requestedByAdmin bool) ([]models.Election, error) {
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}
	if !requestedByAdmin {
		return nil, errors.New(constants.AccessDenied)
	}

	elections, err := e.repo.GetListOfStartedElections(ctx, constants.DESCorder)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return elections, nil
}

func (e election) ElectionResultExists(ctx context.Context, electionId, requesterId string, requestedByAdmin bool) (*bool, error) {
	if e.repo == nil {
		e.repo = repository.NewElectionRepo()
	}

	theElection, err := e.ReadElectionData(ctx, electionId)
	if err != nil {
		return nil, err
	}

	if !requestedByAdmin {
		if theElection.CreatorId != requesterId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	exists, err := e.repo.ElectionResultExists(ctx, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return exists, nil
}
