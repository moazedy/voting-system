package constants

const (
	TitleCanNotBeEmpty                  = "title_can_not_be_empty"
	TitleLengthIsMoreThanMaximom        = "title_length_is_more_than_maximom"
	EndTimeCanNotBeBeforStartTime       = "end_time_can_not_be_befor_start_time"
	StartTimeCanNotBeBeforCreationTime  = "start_time_can_not_be_befor_creation_time"
	ElectionCanNotBeEndedBeforItBegins  = "election_can_not_be_ended_befor_it_begins"
	InvalidElectionType                 = "invalid_election_type"
	ElectionCreatorIdCanNotBeEmpty      = "election_creator_id_can_not_be_empty"
	InvalidCreatorId                    = "invalid_creator_id"
	InvalidCandidatesCountLimit         = "invalid_candidates_count_limit"
	InternalServerError                 = "internal_server_error"
	InvalidElectionId                   = "invalid_election_id"
	AccessDenied                        = "access_denied"
	CandidateNameCanNotBeEmpty          = "candidate_name_can_not_be_empty"
	CandidateNameIsLongerThanExpected   = "candidate_name_is_longer_than_expected"
	InvalidCandidateType                = "invalid_candidate_type"
	ElectionIdCanNotBeEmpty             = "election_id_can_not_be_empty"
	DeletedCanNotBeTrueInTheBegining    = "deleted_can_not_be_true_in_the_begining"
	ElectionDoesNotExist                = "election_does_not_exist"
	CandidateIdCanNotBeEmpty            = "candidate_id_can_not_be_empty"
	InvalidCandidateId                  = "invalid_candidate_id"
	CandidateDoesNotExist               = "candidate_does_not_exist"
	ContributorIdCanNotBeEmpty          = "contributor_id_can_not_be_empty"
	InvalidContributorId                = "invalid_contributor_id"
	VoteDoesNotExist                    = "vote_does_not_exist"
	IdCanNotBeEmpty                     = "id_can_not_be_empty"
	InvalidId                           = "invalid_id"
	ContributorNameCanNotBeEmpty        = "contributor_name_can_not_be_empty"
	ContributorNameIsLongerThanExpected = "contributor_name_is_longer_than_expected"
	ContributorAlredyExists             = "contributor_alredy_exists"
	ContributorDoesNotExist             = "contributor_does_not_exist"
)

// error keies
const (
	ReadingElectionError       = "reading election error"
	AccessError                = "access error"
	ReadingCandidatesError     = "reading candidates error"
	SavingElectionResultsError = "saving election results error"
)
