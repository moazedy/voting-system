package couchbaseQueries

const (
	SaveVoteQuery = ` INSERT INTO ` + constants.VotesBucket + ` (KEY, VALUE) VALUES ($1, $2) `
)
