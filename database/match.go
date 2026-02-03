package database

import (
	"fmt"
)

const (
	AllianceRed  = "red"  // Red alliance
	AllianceBlue = "blue" // Blue alliance
)

// InitMatchStatements prepares all SQL statements for match operations.
func InitMatchStatements() error {
	queries := map[string]string{
		"getMatch":               "SELECT match_id, event_id, match_number, actual_start_time, description, tournament_level FROM matches WHERE match_id = ?",
		"getAllMatches":          "SELECT match_id, event_id, match_number, actual_start_time, description, tournament_level FROM matches",
		"getMatchesByEvent":      "SELECT match_id, event_id, match_number, actual_start_time, description, tournament_level FROM matches WHERE event_id = ? ORDER BY match_number",
		"saveMatch":              "INSERT INTO matches (match_id, event_id, match_number, actual_start_time, description, tournament_level) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE event_id = VALUES(event_id), match_number = VALUES(match_number), actual_start_time = VALUES(actual_start_time), description = VALUES(description), tournament_level = VALUES(tournament_level)",
		"getMatchAllianceScore":  "SELECT match_id, alliance, auto_points, teleop_points, foul_points_committed, pre_foul_total, total_points, major_fouls, minor_fouls FROM match_alliance_scores WHERE match_id = ? AND alliance = ?",
		"saveMatchAllianceScore": "INSERT INTO match_alliance_scores (match_id, alliance, auto_points, teleop_points, foul_points_committed, pre_foul_total, total_points, major_fouls, minor_fouls) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE auto_points = VALUES(auto_points), teleop_points = VALUES(teleop_points), foul_points_committed = VALUES(foul_points_committed), pre_foul_total = VALUES(pre_foul_total), total_points = VALUES(total_points), major_fouls = VALUES(major_fouls), minor_fouls = VALUES(minor_fouls)",
		"getMatchTeams":          "SELECT match_id, team_id, alliance, dq, on_field FROM match_teams WHERE match_id = ?",
		"getTeamsByEvent":        "SELECT DISTINCT mt.team_id FROM match_teams mt INNER JOIN matches m ON mt.match_id = m.match_id WHERE m.event_id = ? ORDER BY mt.team_id",
		"saveMatchTeam":          "INSERT INTO match_teams (match_id, team_id, alliance, dq, on_field) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE alliance = VALUES(alliance), dq = VALUES(dq), on_field = VALUES(on_field)",
	}

	for name, query := range queries {
		if err := PrepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

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

// GetMatchID generates a MatchID from the given EventID and MatchNumber.
func GetMatchID(eventID string, matchNumber int) string {
	return fmt.Sprintf("%s : %d", eventID, matchNumber)
}

// GetMatch retrieves a match from the database by its ID.
func GetMatch(matchID string) *Match {
	var match Match
	stmt := GetStatement("getMatch")
	if stmt == nil {
		return nil
	}
	err := stmt.QueryRow(matchID).Scan(
		&match.MatchID,
		&match.EventID,
		&match.MatchNumber,
		&match.ActualStartTime,
		&match.Description,
		&match.TournamentLevel,
	)
	if err != nil {
		return nil
	}
	return &match
}

// GetAllMatches retrieves all matches from the
func GetAllMatches() []Match {
	stmt := GetStatement("getAllMatches")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.MatchID,
			&match.EventID,
			&match.MatchNumber,
			&match.ActualStartTime,
			&match.Description,
			&match.TournamentLevel,
		)
		if err != nil {
			continue
		}
		matches = append(matches, match)
	}
	return matches
}

// GetMatchesByEvent retrieves all matches for a specific event, ordered by match number.
func GetMatchesByEvent(eventID string) []Match {
	stmt := GetStatement("getMatchesByEvent")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.MatchID,
			&match.EventID,
			&match.MatchNumber,
			&match.ActualStartTime,
			&match.Description,
			&match.TournamentLevel,
		)
		if err != nil {
			continue
		}
		matches = append(matches, match)
	}
	return matches
}

// SaveMatch saves or updates a match in the
func SaveMatch(match *Match) error {
	stmt := GetStatement("saveMatch")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		match.MatchID,
		match.EventID,
		match.MatchNumber,
		match.ActualStartTime,
		match.Description,
		match.TournamentLevel,
	)
	return err
}

// GetMatchAllianceScore retrieves the score for a specific alliance in a match.
func GetMatchAllianceScore(matchID, alliance string) *MatchAllianceScore {
	var score MatchAllianceScore
	stmt := GetStatement("getMatchAllianceScore")
	if stmt == nil {
		return nil
	}
	err := stmt.QueryRow(matchID, alliance).Scan(
		&score.MatchID,
		&score.Alliance,
		&score.AutoPoints,
		&score.TeleopPoints,
		&score.FoulPointsCommitted,
		&score.PreFoulTotal,
		&score.TotalPoints,
		&score.MajorFouls,
		&score.MinorFouls,
	)
	if err != nil {
		return nil
	}
	return &score
}

// SaveMatchAllianceScore saves or updates the score for a specific alliance in a match.
func SaveMatchAllianceScore(score *MatchAllianceScore) error {
	stmt := GetStatement("saveMatchAllianceScore")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		score.MatchID,
		score.Alliance,
		score.AutoPoints,
		score.TeleopPoints,
		score.FoulPointsCommitted,
		score.PreFoulTotal,
		score.TotalPoints,
		score.MajorFouls,
		score.MinorFouls,
	)
	return err
}

// GetMatchTeams retrieves all teams participating in a specific match.
func GetMatchTeams(matchID string) []MatchTeam {
	stmt := GetStatement("getMatchTeams")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(matchID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teams []MatchTeam
	for rows.Next() {
		var team MatchTeam
		if err := rows.Scan(
			&team.MatchID,
			&team.TeamID,
			&team.Alliance,
			&team.Dq,
			&team.OnField,
		); err != nil {
			return nil
		}
		teams = append(teams, team)
	}
	return teams
}

// SaveMatchTeam saves or updates a match team in the
func SaveMatchTeam(team *MatchTeam) error {
	stmt := GetStatement("saveMatchTeam")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		team.MatchID,
		team.TeamID,
		team.Alliance,
		team.Dq,
		team.OnField,
	)
	return err
}

// GetTeamsByEvent retrieves all unique team IDs that participated at a specific event, ordered by team ID.
func GetTeamsByEvent(eventID string) []int {
	stmt := GetStatement("getTeamsByEvent")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teamIDs []int
	for rows.Next() {
		var teamID int
		err := rows.Scan(&teamID)
		if err != nil {
			continue
		}
		teamIDs = append(teamIDs, teamID)
	}
	return teamIDs
}
