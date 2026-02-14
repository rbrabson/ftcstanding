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

	GetAward(awardID int) (*Award, error)
	GetAllAwards() ([]*Award, error)
	SaveAward(award *Award) error

	GetEvent(eventID string) (*Event, error)
	GetAllEvents(filters ...EventFilter) ([]*Event, error)
	SaveEvent(event *Event) error
	GetEventAwards(eventID string) ([]*EventAward, error)
	SaveEventAward(ea *EventAward) error
	GetTeamAwardsByEvent(eventID string, teamID int) ([]*EventAward, error)
	GetAllTeamAwards(teamID int) ([]*EventAward, error)
	GetEventRankings(eventID string) ([]*EventRanking, error)
	SaveEventRanking(er *EventRanking) error
	GetEventAdvancements(eventID string) ([]*EventAdvancement, error)
	SaveEventAdvancement(ea *EventAdvancement) error
	GetEventTeams(eventID string) ([]*EventTeam, error)
	SaveEventTeam(et *EventTeam) error
	GetEventsByTeam(teamID int) ([]string, error)
	GetRegionCodes() ([]string, error)
	GetEventCodesByRegion(regionCode string) ([]string, error)
	GetAdvancementsByRegion(regionCode string) ([]*EventAdvancement, error)
	GetAllAdvancements(filters ...AdvancementFilter) ([]*EventAdvancement, error)

	GetMatch(matchID string) (*Match, error)
	GetAllMatches(filters ...MatchFilter) ([]*Match, error)
	GetMatchesByEvent(eventID string) ([]*Match, error)
	SaveMatch(match *Match) error
	GetMatchAllianceScore(matchID, alliance string) (*MatchAllianceScore, error)
	SaveMatchAllianceScore(score *MatchAllianceScore) error
	GetMatchTeams(matchID string) ([]*MatchTeam, error)
	SaveMatchTeam(team *MatchTeam) error
	GetTeamsByEvent(eventID string) ([]int, error)

	GetTeam(teamID int) (*Team, error)
	GetAllTeams(filters ...TeamFilter) ([]*Team, error)
	SaveTeam(team *Team) error
	GetTeamsByRegion(region string) ([]*Team, error)
	GetTeamRankings(filters ...TeamRankingFilter) ([]*TeamRanking, error)
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
