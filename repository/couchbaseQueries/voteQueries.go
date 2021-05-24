package couchbaseQueries

import "voting-system/constants"

const (
	SaveVoteQuery = ` INSERT INTO ` + constants.VotesBucket + ` (KEY, VALUE) VALUES ($1, $2) `

	ReadVoteQuery = ` SELECT * FROM ` + constants.VotesBucket + ` WHERE id=$1 `

	DeleteVoteQuery = ` UPDATE ` + constants.VotesBucket + ` SET deleted = true AND 
   deleted_at= CLOCK_UTC() WHERE id = $1 `

	GetCandidatePositiveVotesQuery = ` SELECT * FROM ` + constants.VotesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND vote_value=true 
    AND candidate_id = $1 
     ORDER BY %s OFFSET $2 LIMIT $3`

	GetCandidateNegativeVotesQuery = ` SELECT * FROM ` + constants.VotesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND (vote_value = false OR vote_value IS MISSING OR vote_value IS NULL )
     AND candidate_id = $1 
		  ORDER BY %s OFFSET $2 LIMIT $3 `

	GetCandidateVotesQuery = `SELECT * FROM ` + constants.VotesBucket + ` WHERE 
		(deleted = false OR deleted IS MISSING OR deleted IS NULL) AND   
		 candidate_id=$1 ORDER BY vote_time %s OFFSET $2 LIMIT $3`

	GetCandidatePositiveVotesCount = ` SELECT count(*) FROM ` + constants.VotesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND vote_value=true 
    AND candidate_id = $1 `

	GetCandidateNegativeVotesCount = ` SELECT count(*) FROM ` + constants.VotesBucket + ` WHERE 
    (deleted = false OR deleted IS MISSING OR deleted IS NULL) AND 
		 (vote_value = false OR vote_value IS MISSING OR vote_value IS NULL )
     AND candidate_id = $1 `

	UpdateVoteQuery = ` UPDATE ` + constants.VotesBucket + ` SET 
		  candidate_id=$1, contributor_id=$2, vote_value=$3, private_voting=$4,
			election_id=$5 WHERE id= $6`
)
