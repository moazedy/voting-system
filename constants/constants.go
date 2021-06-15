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
	MaximomTitleLength         int = 50
	MaximomCandidatesCount     int = 10
	MaximomCandidateNameLength int = 50
)

const (
	PaginationDefaultPerPage int    = 25
	PaginationDefaultOrder   string = "DESC"
	PaginationMaxPerPage     int    = 100
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
