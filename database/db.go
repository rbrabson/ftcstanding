package database

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// DB defines the interface for database operations.
type DB interface {
	Close()

	GetAward(awardID int) *Award
	GetAllAwards() []*Award
	SaveAward(award *Award) error

	GetEvent(eventID string) *Event
	GetAllEvents(filters ...EventFilter) []*Event
	SaveEvent(event *Event) error
	GetEventAwards(eventID string) []*EventAward
	SaveEventAward(ea *EventAward) error
	GetTeamAwardsByEvent(eventID string, teamID int) []*EventAward
	GetAllTeamAwards(teamID int) []*EventAward
	GetEventRankings(eventID string) []*EventRanking
	SaveEventRanking(er *EventRanking) error
	GetEventAdvancements(eventID string) []*EventAdvancement
	SaveEventAdvancement(ea *EventAdvancement) error
	GetEventTeams(eventID string) []*EventTeam
	SaveEventTeam(et *EventTeam) error
	GetEventsByTeam(teamID int) []string
	GetRegionCodes() []string
	GetEventCodesByRegion(regionCode string) []string
	GetAdvancementsByRegion(regionCode string) []*EventAdvancement
	GetAllAdvancements(filters ...AdvancementFilter) []*EventAdvancement

	GetMatch(matchID string) *Match
	GetAllMatches(filters ...MatchFilter) []*Match
	GetMatchesByEvent(eventID string) []*Match
	SaveMatch(match *Match) error
	GetMatchAllianceScore(matchID, alliance string) *MatchAllianceScore
	SaveMatchAllianceScore(score *MatchAllianceScore) error
	GetMatchTeams(matchID string) []*MatchTeam
	SaveMatchTeam(team *MatchTeam) error
	GetTeamsByEvent(eventID string) []int

	GetTeam(teamID int) *Team
	GetAllTeams(filters ...TeamFilter) []*Team
	SaveTeam(team *Team) error
	GetTeamsByRegion(region string) []*Team
	GetTeamRankings(filters ...TeamRankingFilter) []*TeamRanking
	SaveTeamRanking(ranking *TeamRanking) error
}

// InitDB initializes the database connection.
// season is an optional parameter. If provided, it will be used for file-based databases.
// If not provided, the FTC_SEASON environment variable will be used.
func Init(season ...string) (DB, error) {
	godotenv.Load()
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		return nil, errors.New("DB_TYPE environment variable not set")
	}
	switch dbType {
	case "sql":
		slog.Info("Initializing SQL database")
		return initSQLDB()
	case "file":
		slog.Info("Initializing file database")
		return initFileDB(season...)
	}
	return nil, fmt.Errorf("unsupported DB_TYPE: %s", dbType)
}
