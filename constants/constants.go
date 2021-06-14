package constants

// couchbase bucket names
const (
	VotesBucket           = "votes"
	CandidatesBucket      = "candidates"
	ContributorsBucket    = "contributors"
	ElectionsBucket       = "elections"
	ElectionResultsBucket = "election_results"
)

const (
	MaximomTitleLength         = 50
	MaximomCandidatesCount     = 10
	MaximomCandidateNameLength = 50
)

const (
	MximomContributorNameLength = 50
)

const (
	DESCorder = "desc"
)

const (
	PositiveVotes = " positive votes "
	NegativeVotes = " negative votes "
)

const (
	ElectionManagerWorkerWorkPeriod = 5 // in seconds
)
