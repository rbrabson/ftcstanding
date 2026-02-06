package database

import "fmt"

// InitTeamStatements prepares all SQL statements for team operations.
func (db *sqldb) initTeamStatements() error {
	queries := map[string]string{
		"getTeam":          "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE team_id = ?",
		"getAllTeams":      "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams",
		"getTeamsByRegion": "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE home_region = ? ORDER BY team_id",
		"saveTeam":         "INSERT INTO teams (team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), full_name = VALUES(full_name), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), website = VALUES(website), rookie_year = VALUES(rookie_year), home_region = VALUES(home_region), robot_name = VALUES(robot_name)",
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

	// Build dynamic query
	query := "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE 1=1"
	args := []interface{}{}

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
