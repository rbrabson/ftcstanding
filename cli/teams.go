package cli

import (
	"bytes"

	"github.com/fatih/color"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rodaine/table"
)

// RenderTeams renders a list of teams as a table.
func RenderTeams(teams []*database.Team) string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "Country", "Region", "Rookie Year")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, team := range teams {
		tbl.AddRow(team.TeamID, team.Name, team.Country, team.HomeRegion, team.RookieYear)
	}

	buffer := &bytes.Buffer{}
	tbl.WithWriter(buffer).Print()
	return buffer.String()
}

// RenderTeamMatchDetails renders a list of team match details as a table.
func RenderTeamMatchDetails(details []*query.TeamMatchDetails) string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Event", "Match", "Alliance", "Total Points", "Teams")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, detail := range details {
		teamIDs := make([]int, 0, len(detail.Teams))
		for _, mt := range detail.Teams {
			teamIDs = append(teamIDs, mt.TeamID)
		}
		tbl.AddRow(detail.Event.EventCode, detail.Match.MatchNumber, detail.AllianceScore.Alliance, detail.AllianceScore.TotalPoints, teamIDs)
	}

	buffer := &bytes.Buffer{}
	tbl.WithWriter(buffer).Print()
	return buffer.String()
}
