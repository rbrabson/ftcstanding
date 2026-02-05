package database

import (
	"fmt"

	"github.com/rbrabson/ftc"
)

const (
	AllianceRed  = "red"  // Red alliance
	AllianceBlue = "blue" // Blue alliance
)

// Match represents a match in an event.
type Match struct {
	MatchID         string `json:"matchID"`
	EventID         string `json:"event_id"`
	MatchNumber     int    `json:"matchNumber"`
	ActualStartTime string `json:"actualStartTime"`
	Description     string `json:"description"`
	TournamentLevel string `json:"tournamentLevel"`
}

// MatchAllianceScore represents the score of an alliance in a match. MatchID and Alliance form a composite primary key.
type MatchAllianceScore struct {
	MatchID             string `json:"match_id"`
	Alliance            string `json:"alliance"`
	AutoPoints          int    `json:"auto_points"`
	TeleopPoints        int    `json:"teleop_points"`
	FoulPointsCommitted int    `json:"foul_points_committed"`
	PreFoulTotal        int    `json:"pre_foul_total"`
	TotalPoints         int    `json:"total_points"`
	MajorFouls          int    `json:"major_fouls"`
	MinorFouls          int    `json:"minor_fouls"`
}

// MatchTeam represents an alliance team member participating in a match. MatchID and TeamID form a composite primary key.
type MatchTeam struct {
	MatchID  string `json:"match_id"`
	TeamID   int    `json:"team_id"`
	Alliance string `json:"alliance"`
	Dq       bool   `json:"dq"`
	OnField  bool   `json:"on_field"`
}

// GetMatchID generates a unique ID for a match based on its event ID and match number.
func GetMatchID(event *Event, ftcMatch *ftc.Match) string {
	return fmt.Sprintf("%s : %d", event.EventID, ftcMatch.MatchNumber)
}
