package repository

import (
	"context"
	"log"
	"voting-system/domain/models"
	"voting-system/repository/couchbaseQueries"

	"github.com/couchbase/gocb/v2"
)

//  ElectionRepo is interface of  Election entity in repository layer. other layers of system can interface with it using this
// interface here
type ElectionRepo interface {
	// SaveNewElection saves a new election in db
	SaveNewElection(ctx context.Context, newElection models.Election) (*models.Id, error)
	// ReadElectionData reads some election's data in db using given Id
	ReadElectionData(ctx context.Context, electionId string) (*models.Election, error)
	// DeleteElection deletes given election
	DeleteElection(ctx context.Context, electionId string) error
	// UpdateElection updates some election using received new election data
	UpdateElection(ctx context.Context, electionData models.Election) error
	// GetElectionContributorsCount gets count of given election's contributors
	GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error)
	// ElectionExists checks on given election id existance
	ElectionExists(ctx context.Context, electionId string) (*bool, error)
	// GetListOfRelatedUsers gets list Ids of users being added as related users to the election
	GetListOfRelatedUsers(ctx context.Context, electionId string) ([]models.RelatedPerson, error)
	// GetListOfRelatedCategories gets list of category Ids being added to the election as related categories
	GetListOfRelatedCategories(ctx context.Context, electionId string) ([]models.RelatedCategory, error)
	// GetUserRelatedElections gets list of election Ids that are related to given userId
	GetUserRelatedElections(ctx context.Context, userId string) (*models.Elections, error)
	// GetCategoryRelatedElections gets list of elction Ids that are related to given category id
	GetCategoryRelatedElections(ctx context.Context, categoryId string) (*models.Elections, error)
	// SaveElectionResults saves election results in db
	SaveElectionResult(ctx context.Context, electionResult models.ElectionResults) (*models.Id, error)
}

// election is a struct that represents election entity in repository layer and its the way we can access to repository methods of
// election in this layer
type election struct {
}

// NewEelectionRepo is constractor fucntion for ElectionRepo
func NewElectionRepo() ElectionRepo {
	return new(election)
}

func (e *election) SaveNewElection(ctx context.Context, newElection models.Election) (*models.Id, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.SaveElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newElection.Id, newElection},
	})
	if err != nil {
		log.Println(" error in saving new election, error :", err.Error())
		return nil, err
	}

	var id models.Id
	err = result.One(&id)
	if err != nil {
		log.Println("error in reading new election id value, error : ", err.Error())
		return nil, err
	}
	return &id, nil
}

func (e *election) ReadElectionData(ctx context.Context, electionId string) (*models.Election, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ReadElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in reading election data, error :", err.Error())
		return nil, err
	}
	var elec models.Election
	err = result.One(&elec)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &elec, nil
		}
		log.Println("error in reading election item, error :", err.Error())
		return nil, err
	}

	return &elec, nil
}

func (e *election) DeleteElection(ctx context.Context, electionId string) error {
	_, err := DBS.Couch.Query(couchbaseQueries.DeleteElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println(" error in deleting election, error :", err.Error())
		return err
	}
	return nil
}

func (e *election) UpdateElection(ctx context.Context, electionData models.Election) error {
	_, err := DBS.Couch.Query(couchbaseQueries.UpdateElectionQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{
			electionData.Title,
			electionData.StartTime,
			electionData.EndTime,
			electionData.HasEnded,
			electionData.Type,
			electionData.CandidatesCountLimit,
			electionData.CreatorId,
			electionData.Id,
		},
	})
	if err != nil {
		log.Println(" error in updating election, error :", err.Error())
		return err
	}
	return nil
}

func (e *election) GetElectionContributorsCount(ctx context.Context, electionId string) (*models.ContributorsCount, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetElectionContributorsCountQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var count models.ContributorsCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &count, nil
		}

		log.Println("error in reading contributors count, error :", err.Error())
		return nil, err
	}

	return &count, nil
}

func (e *election) ElectionExists(ctx context.Context, electionId string) (*bool, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.ElectionExistsQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var exists bool
	var count models.ElectionsCount
	err = result.One(&count)
	if err != nil {
		if err == gocb.ErrNoResult {
			return &exists, nil
		}
		log.Println("error in reading elections count, error :", err.Error())
		return nil, err
	}

	if count.Count > 0 {
		exists = true
		return &exists, nil
	}

	return &exists, nil

}

func (e *election) GetListOfRelatedUsers(ctx context.Context, electionId string) ([]models.RelatedPerson, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetListOfRelatedUsersQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var users []models.RelatedPerson
	for result.Next() {
		var user models.RelatedPerson
		err := result.Row(&user)
		if err != nil {
			if err == gocb.ErrNoResult {
				return users, nil
			}
			log.Println("error in reading election related users, error :", err.Error())
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (e *election) GetListOfRelatedCategories(ctx context.Context, electionId string) ([]models.RelatedCategory, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetListOfRelatedCategoriesQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var cats []models.RelatedCategory
	for result.Next() {
		var cat models.RelatedCategory
		err := result.Row(&cat)
		if err != nil {
			if err == gocb.ErrNoResult {
				return cats, nil
			}
			log.Println("error in reading election related categories, error :", err.Error())
			return nil, err
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

func (e *election) GetUserRelatedElections(ctx context.Context, userId string) (*models.Elections, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetUserRelatedElectionsQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{userId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var elections models.Elections
	for result.Next() {
		var election models.ElectionId
		err := result.Row(&election)
		if err != nil {
			if err == gocb.ErrNoResult {
				return &elections, nil
			}
			log.Println("error in reading election id item, error : ", err.Error())
			return nil, err
		}

		elections.Elections = append(elections.Elections, election)
	}

	return &elections, nil
}

func (e election) GetCategoryRelatedElections(ctx context.Context, categoryId string) (*models.Elections, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.GetCategoryRelatedElectionsQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{categoryId},
	})
	if err != nil {
		log.Println("error in query execution, error :", err.Error())
		return nil, err
	}

	var elections models.Elections
	for result.Next() {
		var election models.ElectionId
		err := result.Row(&election)
		if err != nil {
			if err == gocb.ErrNoResult {
				return &elections, nil
			}
			log.Println("error in reading election id item, error : ", err.Error())
			return nil, err
		}

		elections.Elections = append(elections.Elections, election)
	}

	return &elections, nil

}

func (e election) SaveElectionResult(ctx context.Context, electionResult models.ElectionResults) (*models.Id, error) {
	result, err := DBS.Couch.Query(couchbaseQueries.SaveElectionResultsQuery, &gocb.QueryOptions{
		PositionalParameters: []interface{}{electionResult.Id, electionResult},
	})
	if err != nil {
		log.Println(" error in saving new election result, error :", err.Error())
		return nil, err
	}

	var id models.Id
	err = result.One(&id)
	if err != nil {
		log.Println("error in reading new election result id value, error : ", err.Error())
		return nil, err
	}
	return &id, nil
}
