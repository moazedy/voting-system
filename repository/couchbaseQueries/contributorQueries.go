package couchbaseQueries

import "voting-system/constants"

const (
	SaveContributorQuery = ` INSERT INTO ` + constants.ContributorsBucket + ` (KEY, VALUE) VALUES ($1, $2) 
	 RETURNING meta().id AS id`

	GetElectionContributorsQuery = ` SELECT * FROM ` + constants.ContributorsBucket +
		` WHERE (deleted=false OR deleted IS MISSING OR deleted IS NULL ) AND 
         election_id=$1 ORDER BY contribute_time %s OFFSET $2 LIMIT $3 `

	ReadContributorQuery = `SELECT * FROM ` + constants.ContributorsBucket + ` WHERE 
          (deleted =false OR deleted IS MISSING OR deleted IS NULL) AND 
           meta_id=$1 `

	DeleteContributorQuery = ` UPDATE ` + constants.ContributorsBucket + ` SET 
            deleted= true, deleted_at= CLOCK_UTC() WHERE meta_id= $1`

	ContributorExistanceQuery = ` SELECT count(*) FROM ` + constants.ContributorsBucket +
		` WHERE (deleted= FALSE OR deleted IS MISSING OR deleted IS NULL) AND 
              contributor_id=$1  AND election_id = $2 `

	ContributionExistanceQuery = ` SELECT count(*) FROM ` + constants.ContributorsBucket +
		` WHERE (deleted= FALSE OR deleted IS MISSING OR deleted IS NULL) AND 
              meta_id =$1 `
)
