package couchbaseQueries

import "voting-system/constants"

const (
	SaveContributorQuery = ` INSERT INTO ` + constants.ContributorsBucket + ` (KEY, VALUE) VALUES ($1, $2)`

	GetElectionContributorsQuery = ` SELECT * FROM ` + constants.ContributorsBucket +
		` WHERE (deleted=false OR deleted IS MISSING OR deleted IS NULL ) AND 
         election_id=$1 ORDER BY contribute_time %s OFFSET $2 LIMIT $3 `

	ReadContributorQuery = `SELECT * FROM ` + constants.ContributorsBucket + ` WHERE 
          (deleted =false OR deleted IS MISSING OR deleted IS NULL) AND 
           id=$1 `

	DeleteContributorQuery = ` UPDATE ` + constants.ContributorsBucket + ` SET 
            deleted= true, deleted_at= CLOCK_UTC() WHERE id= $1`
)
