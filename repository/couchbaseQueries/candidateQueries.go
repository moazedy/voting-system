package couchbaseQueries

import "voting-system/constants"

const (
	SaveNewCandidateQuery = ` INSERT INTO ` + constants.CandidatesBucket + ` (KEY,VALUE) 
   VALUES ($1, $2) RETURNING meta().id AS id`

	ReadCandidateDataQuery = ` SELECT * FROM ` + constants.CandidatesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND 
     id= $1 `

	GetElectionCandidatesQuery = ` SELECT * FROM ` + constants.CandidatesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND 
      election_id=$1 ORDER BY created_at %s OFFSET $2 LIMIT $3 `

	DeleteCandidateQuery = ` UPDATE ` + constants.CandidatesBucket + ` SET
       deleted = true, deleted_at = ClOCK_UTC()  WHERE id=$1 `

	UpdateCandidateQuery = ` UPDATE ` + constants.CandidatesBucket + ` SET 
   name=$1, type=$2, descriptions=$3, election_id=$4  WHERE 
   id=$5 `

	IsCandidateExistsQuery = ` SELECT count(*) FROM ` + constants.CandidatesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND 
		id = $1 
   `
)
