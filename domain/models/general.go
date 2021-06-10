package models

type Id struct {
	Id string `json:"id"`
}

type Count struct {
	Count int `json:"count"`
}

type ContributionData struct {
	ContributorId string `json:"contributor_id"`
	VoteId        string `json:"vote_id"`
	ElectionId    string `json:"election_id"`
}
