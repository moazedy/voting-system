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
       title=$1, start_time=$2, end_time=$3, has_ended=$4 type=$5, candidate_count_limit=$6 
			  WHERE id=$7 `

	GetElectionContributorsCountQuery = ` SELECT count(*) FROM ` + constants.ContributorsBucket +
		` WHERE (deleted=false OR deleted IS MISSING OR deleted IS NULL ) AND 
         election_id=$1 `
)
