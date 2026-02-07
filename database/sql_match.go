package database

import "fmt"

// InitMatchStatements prepares all SQL statements for match operations.
func (db *sqldb) initMatchStatements() error {
	queries := map[string]string{
		"getMatch":               "SELECT match_id, event_id, match_type, match_number, actual_start_time, description, tournament_level FROM matches WHERE match_id = ?",
		"getAllMatches":          "SELECT match_id, event_id, match_type, match_number, actual_start_time, description, tournament_level FROM matches",
		"getMatchesByEvent":      "SELECT match_id, event_id, match_type, match_number, actual_start_time, description, tournament_level FROM matches WHERE event_id = ? ORDER BY match_number",
		"saveMatch":              "INSERT INTO matches (match_id, event_id, match_type, match_number, actual_start_time, description, tournament_level) VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE event_id = VALUES(event_id), match_type = VALUES(match_type), match_number = VALUES(match_number), actual_start_time = VALUES(actual_start_time), description = VALUES(description), tournament_level = VALUES(tournament_level)",
		"getMatchAllianceScore":  "SELECT match_id, alliance, auto_points, teleop_points, foul_points_committed, pre_foul_total, total_points, major_fouls, minor_fouls FROM match_alliance_scores WHERE match_id = ? AND alliance = ?",
		"saveMatchAllianceScore": "INSERT INTO match_alliance_scores (match_id, alliance, auto_points, teleop_points, foul_points_committed, pre_foul_total, total_points, major_fouls, minor_fouls) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE auto_points = VALUES(auto_points), teleop_points = VALUES(teleop_points), foul_points_committed = VALUES(foul_points_committed), pre_foul_total = VALUES(pre_foul_total), total_points = VALUES(total_points), major_fouls = VALUES(major_fouls), minor_fouls = VALUES(minor_fouls)",
		"getMatchTeams":          "SELECT match_id, team_id, alliance, dq, on_field FROM match_teams WHERE match_id = ?",
		"getTeamsByEvent":        "SELECT DISTINCT mt.team_id FROM match_teams mt INNER JOIN matches m ON mt.match_id = m.match_id WHERE m.event_id = ? ORDER BY mt.team_id",
		"saveMatchTeam":          "INSERT INTO match_teams (match_id, team_id, alliance, dq, on_field) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE alliance = VALUES(alliance), dq = VALUES(dq), on_field = VALUES(on_field)",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// GetMatchID generates a MatchID from the given EventID and MatchNumber.
func (db *sqldb) GetMatchID(eventID string, matchNumber int) string {
	return fmt.Sprintf("%s : %d", eventID, matchNumber)
}

// GetMatch retrieves a match from the database by its ID.
func (db *sqldb) GetMatch(matchID string) *Match {
	var match Match
	stmt := db.getStatement("getMatch")
	if stmt == nil {
		return nil
	}
	err := stmt.QueryRow(matchID).Scan(
		&match.MatchID,
		&match.EventID,
		&match.MatchType,
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

// GetAllMatches retrieves all matches from the database with optional filters.
// If no filters are provided, returns all matches.
// Filters are combined with OR logic within each field.
func (db *sqldb) GetAllMatches(filters ...MatchFilter) []*Match {
	// If no filters, use the prepared statement
	if len(filters) == 0 {
		stmt := db.getStatement("getAllMatches")
		if stmt == nil {
			return nil
		}
		rows, err := stmt.Query()
		if err != nil {
			return nil
		}
		defer rows.Close()

		var matches []*Match
		for rows.Next() {
			var match Match
			err := rows.Scan(
				&match.MatchID,
				&match.EventID,
				&match.MatchType,
				&match.MatchNumber,
				&match.ActualStartTime,
				&match.Description,
				&match.TournamentLevel,
			)
			if err != nil {
				continue
			}
			matches = append(matches, &match)
		}
		return matches
	}

	filter := filters[0]

	// Build dynamic query
	query := "SELECT match_id, event_id, match_type, match_number, actual_start_time, description, tournament_level FROM matches"
	args := []interface{}{}

	if len(filter.EventIDs) > 0 {
		query += " WHERE event_id IN ("
		for i, eventID := range filter.EventIDs {
			if i > 0 {
				query += ","
			}
			query += "?"
			args = append(args, eventID)
		}
		query += ")"
	}

	query += " ORDER BY event_id, match_number"

	// Execute query
	rows, err := db.sqldb.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var matches []*Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.MatchID,
			&match.EventID,
			&match.MatchType,
			&match.MatchNumber,
			&match.ActualStartTime,
			&match.Description,
			&match.TournamentLevel,
		)
		if err != nil {
			continue
		}
		matches = append(matches, &match)
	}
	return matches
}

// GetMatchesByEvent retrieves all matches for a specific event, ordered by match number.
func (db *sqldb) GetMatchesByEvent(eventID string) []*Match {
	stmt := db.getStatement("getMatchesByEvent")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var matches []*Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.MatchID,
			&match.EventID,
			&match.MatchType,
			&match.MatchNumber,
			&match.ActualStartTime,
			&match.Description,
			&match.TournamentLevel,
		)
		if err != nil {
			continue
		}
		matches = append(matches, &match)
	}
	return matches
}

// SaveMatch saves or updates a match in the
func (db *sqldb) SaveMatch(match *Match) error {
	stmt := db.getStatement("saveMatch")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		match.MatchID,
		match.EventID,
		match.MatchType,
		match.MatchNumber,
		match.ActualStartTime,
		match.Description,
		match.TournamentLevel,
	)
	return err
}

// GetMatchAllianceScore retrieves the score for a specific alliance in a match.
func (db *sqldb) GetMatchAllianceScore(matchID, alliance string) *MatchAllianceScore {
	var score MatchAllianceScore
	stmt := db.getStatement("getMatchAllianceScore")
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
func (db *sqldb) SaveMatchAllianceScore(score *MatchAllianceScore) error {
	stmt := db.getStatement("saveMatchAllianceScore")
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
func (db *sqldb) GetMatchTeams(matchID string) []*MatchTeam {
	stmt := db.getStatement("getMatchTeams")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(matchID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teams []*MatchTeam
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
		teams = append(teams, &team)
	}
	return teams
}

// SaveMatchTeam saves or updates a match team in the
func (db *sqldb) SaveMatchTeam(team *MatchTeam) error {
	stmt := db.getStatement("saveMatchTeam")
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
func (db *sqldb) GetTeamsByEvent(eventID string) []int {
	stmt := db.getStatement("getTeamsByEvent")
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
