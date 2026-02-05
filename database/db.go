package database

// InitDB initializes the database connection.
func Init() (*sqldb, error) {
	return InitSQLDB()
}

// DB defines the interface for database operations.
type DB interface {
	Close()

	GetAward(awardID int) *Award
	GetAllAwards() []Award
	SaveAward(award *Award) error

	GetEvent(eventID string) *Event
	SaveEvent(event *Event) error
	GetEventAwards(eventID string) []*EventAward
	SaveEventAward(ea *EventAward) error
	GetTeamAwardsByEvent(eventID string, teamID int) []*EventAward
	GetAllTeamAwards(teamID int) []*EventAward
	GetEventRankings(eventID string) []*EventRanking
	SaveEventRanking(er *EventRanking) error
	GetEventAdvancements(eventID string) []*EventAdvancement
	SaveEventAdvancement(ea *EventAdvancement) error
	GetRegionCodes() []string
	GetEventCodesByRegion(regionCode string) []string
	GetAdvancementsByRegion(regionCode string) []*EventAdvancement
	GetAllAdvancements() []*EventAdvancement

	GetMatch(matchID string) *Match
	GetAllMatches() []*Match
	GetMatchesByEvent(eventID string) []*Match
	SaveMatch(match *Match) error
	GetMatchAllianceScore(matchID, alliance string) *MatchAllianceScore
	SaveMatchAllianceScore(score *MatchAllianceScore) error
	GetMatchTeams(matchID string) []*MatchTeam
	SaveMatchTeam(team *MatchTeam) error
	GetTeamsByEvent(eventID string) []int

	GetTeam(teamID int) *Team
	GetAllTeams() []*Team
	SaveTeam(team *Team) error
	GetTeamsByRegion(region string) []*Team
}
