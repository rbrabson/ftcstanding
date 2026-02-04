package request

import (
	"strconv"

	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
)

// RequestAndStoreTeams retrieves the list of teams for a given season and stores them in the database.
func RequestAndStoreTeams(season int) {
	teams := RequestTeams(season)
	if teams == nil {
		return
	}
	for _, team := range teams {
		database.SaveTeam(team)
	}
}

// RequestTeams retrieves the list of teams for a given season.
func RequestTeams(season int) []*database.Team {
	ftcTeams, err := ftc.GetTeams(strconv.Itoa(season))
	if err != nil {
		return nil
	}
	teams := make([]*database.Team, 0, len(ftcTeams))
	for _, ftcTeam := range ftcTeams {
		team := database.Team{
			TeamID:     ftcTeam.TeamNumber,
			Name:       ftcTeam.NameShort,
			FullName:   ftcTeam.NameFull,
			City:       ftcTeam.City,
			StateProv:  ftcTeam.StateProv,
			Country:    ftcTeam.Country,
			RookieYear: ftcTeam.RookieYear,
		}
		if ftcTeam.Website != nil {
			team.Website = *ftcTeam.Website
		}
		if ftcTeam.HomeRegion != nil {
			team.HomeRegion = *ftcTeam.HomeRegion
		}
		if ftcTeam.RobotName != nil {
			team.RobotName = *ftcTeam.RobotName
		}
		teams = append(teams, &team)
	}
	return teams
}
