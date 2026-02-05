package database

import "fmt"

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

// GetAllTeams retrieves all teams from the database.
func (db *sqldb) GetAllTeams() []*Team {
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
