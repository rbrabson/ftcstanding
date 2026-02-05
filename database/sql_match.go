package database

import "fmt"

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
func (db *sqldb) GetAllMatches() []*Match {
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
