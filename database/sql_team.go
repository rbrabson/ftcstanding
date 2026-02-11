package database

import "fmt"

// InitTeamStatements prepares all SQL statements for team operations.
func (db *sqldb) initTeamStatements() error {
	queries := map[string]string{
		"getTeam":          "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE team_id = ?",
		"getAllTeams":      "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams ORDER BY team_id",
		"getTeamsByRegion": "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE home_region = ? ORDER BY team_id",
		"saveTeam":         "INSERT INTO teams (team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), full_name = VALUES(full_name), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), website = VALUES(website), rookie_year = VALUES(rookie_year), home_region = VALUES(home_region), robot_name = VALUES(robot_name)",
		"saveTeamRanking":  "INSERT INTO team_rankings (team_id, event_id, num_matches, ccwm, opr, np_opr, dpr, np_dpr, np_avg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE num_matches = VALUES(num_matches), ccwm = VALUES(ccwm), opr = VALUES(opr), np_opr = VALUES(np_opr), dpr = VALUES(dpr), np_dpr = VALUES(np_dpr), np_avg = VALUES(np_avg)",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// GetTeam retrieves a team from a database by its ID.
func (db *sqldb) GetTeam(teamID int) *Team {
	var team Team
	stmt := db.getStatement("getTeam")
	if stmt == nil {
		return nil
	}
	err := stmt.QueryRow(teamID).Scan(
		&team.TeamID,
		&team.Name,
		&team.FullName,
		&team.City,
		&team.StateProv,
		&team.Country,
		&team.Website,
		&team.RookieYear,
		&team.HomeRegion,
		&team.RobotName,
	)
	if err != nil {
		return nil
	}
	return &team
}

// GetAllTeams retrieves all teams from the database with optional filters.
// If no filters are provided, returns all teams.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *sqldb) GetAllTeams(filters ...TeamFilter) []*Team {
	// If no filters, use the prepared statement
	if len(filters) == 0 {
		stmt := db.getStatement("getAllTeams")
		if stmt == nil {
			return nil
		}
		rows, err := stmt.Query()
		if err != nil {
			return nil
		}
		defer rows.Close()

		var teams []*Team
		for rows.Next() {
			var team Team
			err := rows.Scan(
				&team.TeamID,
				&team.Name,
				&team.FullName,
				&team.City,
				&team.StateProv,
				&team.Country,
				&team.Website,
				&team.RookieYear,
				&team.HomeRegion,
				&team.RobotName,
			)
			if err != nil {
				continue
			}
			teams = append(teams, &team)
		}
		return teams
	}

	filter := filters[0]

	// If EventCodes filter is provided, get team IDs from those events
	var eventTeamIDs []int
	if len(filter.EventCodes) > 0 {
		eventTeamIDMap := make(map[int]bool)
		for _, eventCode := range filter.EventCodes {
			// Get all events matching this code
			events := db.GetAllEvents(EventFilter{EventCodes: []string{eventCode}})
			for _, event := range events {
				// Get all teams for this event
				eventTeams := db.GetEventTeams(event.EventID)
				for _, et := range eventTeams {
					eventTeamIDMap[et.TeamID] = true
				}
			}
		}
		// Convert map to slice
		for teamID := range eventTeamIDMap {
			eventTeamIDs = append(eventTeamIDs, teamID)
		}
	}

	// Build dynamic query
	query := "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE 1=1"
	args := []interface{}{}

	// Add EventCodes filter (team must be in the events)
	if len(eventTeamIDs) > 0 {
		query += " AND team_id IN ("
		for i, id := range eventTeamIDs {
			if i > 0 {
				query += ","
			}
			query += "?"
			args = append(args, id)
		}
		query += ")"
	}

	// Add TeamID filter
	if len(filter.TeamIDs) > 0 {
		query += " AND team_id IN ("
		for i, id := range filter.TeamIDs {
			if i > 0 {
				query += ","
			}
			query += "?"
			args = append(args, id)
		}
		query += ")"
	}

	// Add Country filter
	if len(filter.Countries) > 0 {
		query += " AND country IN ("
		for i, country := range filter.Countries {
			if i > 0 {
				query += ","
			}
			query += "?"
			args = append(args, country)
		}
		query += ")"
	}

	// Add HomeRegion filter
	if len(filter.HomeRegions) > 0 {
		query += " AND home_region IN ("
		for i, region := range filter.HomeRegions {
			if i > 0 {
				query += ","
			}
			query += "?"
			args = append(args, region)
		}
		query += ")"
	}

	query += " ORDER BY team_id"

	// Execute query
	rows, err := db.sqldb.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teams []*Team
	for rows.Next() {
		var team Team
		err := rows.Scan(
			&team.TeamID,
			&team.Name,
			&team.FullName,
			&team.City,
			&team.StateProv,
			&team.Country,
			&team.Website,
			&team.RookieYear,
			&team.HomeRegion,
			&team.RobotName,
		)
		if err != nil {
			continue
		}
		teams = append(teams, &team)
	}
	return teams
}

// SaveTeam saves or updates a team in the
func (db *sqldb) SaveTeam(team *Team) error {
	stmt := db.getStatement("saveTeam")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		team.TeamID,
		team.Name,
		team.FullName,
		team.City,
		team.StateProv,
		team.Country,
		team.Website,
		team.RookieYear,
		team.HomeRegion,
		team.RobotName,
	)
	return err
}

// GetTeamsByRegion retrieves all teams in a given home region, ordered by team ID.
func (db *sqldb) GetTeamsByRegion(region string) []*Team {
	stmt := db.getStatement("getTeamsByRegion")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(region)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teams []*Team
	for rows.Next() {
		var team Team
		err := rows.Scan(
			&team.TeamID,
			&team.Name,
			&team.FullName,
			&team.City,
			&team.StateProv,
			&team.Country,
			&team.Website,
			&team.RookieYear,
			&team.HomeRegion,
			&team.RobotName,
		)
		if err != nil {
			continue
		}
		teams = append(teams, &team)
	}
	return teams
}

// GetTeamRankings retrieves team rankings with optional filters.
// Filters support filtering by TeamID and/or EventID.
// If no filters are provided, returns all team rankings.
func (db *sqldb) GetTeamRankings(filters ...TeamRankingFilter) []*TeamRanking {
	// Build dynamic query
	query := "SELECT team_id, event_id, num_matches, ccwm, opr, np_opr, dpr, np_dpr, np_avg FROM team_rankings WHERE 1=1"
	args := []interface{}{}

	if len(filters) > 0 {
		filter := filters[0]

		// Add TeamID filter
		if len(filter.TeamIDs) > 0 {
			query += " AND team_id IN ("
			for i, id := range filter.TeamIDs {
				if i > 0 {
					query += ","
				}
				query += "?"
				args = append(args, id)
			}
			query += ")"
		}

		// Add EventID filter
		if len(filter.EventIDs) > 0 {
			query += " AND event_id IN ("
			for i, id := range filter.EventIDs {
				if i > 0 {
					query += ","
				}
				query += "?"
				args = append(args, id)
			}
			query += ")"
		}
	}

	query += " ORDER BY event_id, team_id"

	// Execute query
	rows, err := db.sqldb.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var rankings []*TeamRanking
	for rows.Next() {
		var ranking TeamRanking
		err := rows.Scan(
			&ranking.TeamID,
			&ranking.EventID,
			&ranking.NumMatches,
			&ranking.CCWM,
			&ranking.OPR,
			&ranking.NpOPR,
			&ranking.DPR,
			&ranking.NpDPR,
			&ranking.NpAvg,
		)
		if err != nil {
			continue
		}
		rankings = append(rankings, &ranking)
	}
	return rankings
}

// SaveTeamRanking saves or updates a team ranking in the database.
func (db *sqldb) SaveTeamRanking(ranking *TeamRanking) error {
	stmt := db.getStatement("saveTeamRanking")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		ranking.TeamID,
		ranking.EventID,
		ranking.NumMatches,
		ranking.CCWM,
		ranking.OPR,
		ranking.NpOPR,
		ranking.DPR,
		ranking.NpDPR,
		ranking.NpAvg,
	)
	return err
}
