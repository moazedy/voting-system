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
	SaveNewContributor(ctx context.Context, requesterId string, contributorData models.Contributor) error
}

type contributor struct {
	repo repository.ContributorRepo
}

func NewContributorLogic() ContributorLogic {
	return new(contributor)
}

func (c contributor) SaveNewContributor(ctx context.Context, requesterId string, contributorData models.Contributor) error {
	if c.repo == nil {
		c.repo = repository.NewContributorRepo()
	}

	if err := IdValidation(requesterId); err != nil {
		return err
	}

	if err := contributorData.Validate(); err != nil {
		return err
	}

	uid, err := uuid.Parse(requesterId)
	if err != nil {
		return errors.New(constants.InvalidContributorId)
	}
	contributorData.Id = uid

	err = c.repo.SaveNewContributor(ctx, contributorData)
	if err != nil {
		return errors.New(constants.InternalServerError)
	}

	return nil
}
