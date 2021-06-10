package logic

import (
	"context"
	"errors"
	"voting-system/constants"
	"voting-system/domain/models"
	"voting-system/repository"

	"github.com/google/uuid"
)

type ContributorLogic interface {
	// SaveNewContributor adds a new contributor to system db, using received contributor data
	SaveNewContributor(ctx context.Context, requesterId string, contributorData models.Contributor) (*models.Id, error)
	// ReadContributorData gets data of a specific contributor
	ReadContributor(ctx context.Context, contributorMetaId, requesterId string, requestedByAdmin bool) (*models.Contributor, error)
	// ContributorExists checks for contributor existance in election
	ContributorExists(ctx context.Context, contributorId, electionId string) (*bool, error)
}

type contributor struct {
	repo repository.ContributorRepo
}

func NewContributorLogic() ContributorLogic {
	return new(contributor)
}

func (c contributor) SaveNewContributor(ctx context.Context, requesterId string, contributorData models.Contributor) (*models.Id, error) {
	if c.repo == nil {
		c.repo = repository.NewContributorRepo()
	}

	if err := IdValidation(requesterId); err != nil {
		return nil, err
	}

	if err := contributorData.Validate(); err != nil {
		return nil, err
	}

	// checking on existing data for participation of requester in the election
	exists, err := c.repo.IsContributorExists(ctx, requesterId, contributorData.ElectionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}
	if *exists {
		return nil, errors.New(constants.ContributorAlredyExists)
	}

	contributorData.ContributorId = requesterId
	contributorData.MetaId = uuid.New()

	id, err := c.repo.SaveNewContributor(ctx, contributorData)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return id, nil
}

func (c contributor) ReadContributor(ctx context.Context, contributorMetaId, requesterId string, requestedByAdmin bool) (*models.Contributor, error) {
	if c.repo == nil {
		c.repo = repository.NewContributorRepo()
	}

	if err := IdValidation(contributorMetaId); err != nil {
		return nil, err
	}

	exists, err := c.repo.IsContributionExists(ctx, contributorMetaId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}
	if !*exists {
		return nil, errors.New(constants.ContributorDoesNotExist)
	}

	theContributor, err := c.repo.ReadContributorData(ctx, contributorMetaId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	if !requestedByAdmin {
		if theContributor.ContributorId != requesterId {
			return nil, errors.New(constants.AccessDenied)
		}
	}

	return theContributor, nil
}

func (c contributor) ContributorExists(ctx context.Context, contributorId, electionId string) (*bool, error) {
	if c.repo == nil {
		c.repo = repository.NewContributorRepo()
	}

	exists, err := c.repo.IsContributorExists(ctx, contributorId, electionId)
	if err != nil {
		return nil, errors.New(constants.InternalServerError)
	}

	return exists, nil
}
