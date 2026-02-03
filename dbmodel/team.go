package dbmodel

import (
	"fmt"

	"github.com/rbrabson/ftcstanding/database"
)

// Team represents a team that participates in competitions.
type Team struct {
	TeamID     int    `json:"team_id"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	City       string `json:"city"`
	StateProv  string `json:"state_prov"`
	Country    string `json:"country"`
	Website    string `json:"website"`
	RookieYear int    `json:"rookie_year"`
	HomeRegion string `json:"home_region"`
	RobotName  string `json:"robot_name"`
}

// InitTeamStatements prepares all SQL statements for team operations.
func InitTeamStatements() error {
	queries := map[string]string{
		"getTeam":     "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE team_id = ?",
		"getAllTeams": "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams",
		"saveTeam":    "INSERT INTO teams (team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), full_name = VALUES(full_name), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), website = VALUES(website), rookie_year = VALUES(rookie_year), home_region = VALUES(home_region), robot_name = VALUES(robot_name)",
	}

	for name, query := range queries {
		if err := database.PrepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// GetTeam retrieves a team from a database by its ID.
func GetTeam(teamID int) *Team {
	var team Team
	stmt := database.GetStatement("getTeam")
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
func GetAllTeams() []Team {
	stmt := database.GetStatement("getAllTeams")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teams []Team
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
		teams = append(teams, team)
	}
	return teams
}

// SaveTeam saves or updates a team in the database.
func SaveTeam(team *Team) error {
	stmt := database.GetStatement("saveTeam")
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
