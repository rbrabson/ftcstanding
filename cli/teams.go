package cli

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
)

func RenderTeams(teams []*database.Team) string {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"ID", "Name", "Country", "Region", "Rookie Year"})

	for _, team := range teams {
		tw.AppendRow(
			table.Row{team.TeamID, team.Name, team.Country, team.HomeRegion, team.RookieYear},
		)
	}

	return tw.Render()
}

// RenderTeamMatchDetails renders a list of team match details as a table.
func RenderTeamMatchDetails(details []*query.TeamMatchDetails) string {

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Event", "Match", "Alliance", "Total Points", "Teams"})

	for _, detail := range details {
		teamIDs := make([]int, 0, len(detail.Teams))
		for _, mt := range detail.Teams {
			teamIDs = append(teamIDs, mt.TeamID)
		}
		tw.AppendRow(table.Row{detail.Event.EventCode, detail.Match.MatchNumber, detail.AllianceScore.Alliance, detail.AllianceScore.TotalPoints, teamIDs})
	}

	return tw.Render()
}
