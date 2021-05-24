package couchbaseQueries

import "voting-system/constants"

const (
	SaveNewCandidateQuery = ` INSERT INTO ` + constants.CandidatesBucket + ` (KEY,VALUE) 
   VALUES ($1, $2)`
)
