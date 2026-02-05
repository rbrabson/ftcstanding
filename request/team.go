package request

import (
	"log/slog"

	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
)

// RequestAndSaveTeams retrieves the list of teams for a given season and stores them in the database.
func RequestAndSaveTeams(season string) {
	teams := RequestTeams(season)
	if teams == nil {
		return
	}
	for _, team := range teams {
		db.SaveTeam(team)
	}
}

// RequestTeams retrieves the list of teams for a given season.
func RequestTeams(season string) []*database.Team {
	ftcTeams, err := ftc.GetTeams(season)
	if err != nil {
		slog.Error("Error requesting teams:", "error", err)
		return nil
	}
	slog.Debug("Requesting teams...", "count", len(ftcTeams))
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
	slog.Info("Finished requesting teams", "count", len(teams))
	return teams
}
