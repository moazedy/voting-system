package couchbaseQueries

import "voting-system/constants"

const (
	SaveElectionQuery = `INSERT INTO ` + constants.ElectionsBucket + ` (KEY, VALUE) 
   VALUES ($1, $2)`

	ReadElectionQuery = ` SELECT * FROM ` + constants.ElectionsBucket + ` WHERE 
    (deleted= false OR deleted IS MISSING OR deleted IS NULL) AND
     id=$1`

	DeleteElectionQuery = `UPDATE ` + constants.ElectionsBucket + ` SET 
      deleted=true, deleted_at= CLOCK_UTC() WHERE id=$1`

	UpdateElectionQuery = ` UPDATE ` + constants.ElectionsBucket + ` SET 
       title=$1, start_time=$2, end_time=$3, has_ended=$4 type=$5, candidate_count_limit=$6,
       creator_id= $7 WHERE id=$8 `

	GetElectionContributorsCountQuery = ` SELECT count(*) FROM ` + constants.ContributorsBucket +
		` WHERE (deleted=false OR deleted IS MISSING OR deleted IS NULL ) AND 
         election_id=$1 `

	ElectionExistsQuery = ` SELECT count(*) FROM ` + constants.ElectionsBucket +
		` WHERE (deleted=false OR deleted IS MISSING OR deleted IS NULL ) AND 
         id=$1 `

	GetListOfRelatedUsersQuery = ` SELECT related_users FROM ` + constants.ElectionsBucket + ` WHERE 
           (deleted= false OR deleted IS MISSING OR deleted IS NULL) AND 
            id = $1 
         `
	GetListOfRelatedCategoriesQuery = ` SELECT related_categories FROM ` + constants.ElectionsBucket + ` WHERE 
          (deleted= false OR deleted IS MISSING OR deleted IS NULL) AND 
            id = $1 
         `

	GetUserRelatedElectionsQuery = ` SELECT e.id FROM ` + constants.ElectionsBucket + ` AS e 
          UNNEST e.related_persons AS persons
          WHERE (deleted= false OR deleted IS MISSING OR deleted IS NULL) AND
          persons.id = $1 `

	GetCategoryRelatedElectionsQuery = ` SELECT e.id FROM ` + constants.ElectionsBucket + ` AS e 
          UNNEST e.related_categories AS cats
          WHERE (deleted= false OR deleted IS MISSING OR deleted IS NULL) AND
          cats.id = $1 `
)
