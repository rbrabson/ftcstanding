package cli

import (
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
)

func RenderTeams(teams []*database.Team) string {
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
			BG: renderer.Colors{color.BgHiWhite},           // White background
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0
				{},                                     // Inherit default (cyan) for column 1
				{FG: renderer.Colors{color.FgHiRed}},   // High-intensity red for column 2
				{},                                     // Inherit default (cyan) for remaining columns
			},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow bold footer
			Columns: []renderer.Tint{
				{},                                      // Inherit default
				{FG: renderer.Colors{color.FgHiYellow}}, // High-intensity yellow for column 1
				{},                                      // Inherit default
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	}

	var sb strings.Builder
	table := tablewriter.NewTable(&sb, tablewriter.WithRenderer(renderer.NewColorized(colorCfg)))
	table.Header([]string{"ID", "Name", "Country", "Region", "Rookie Year"})

	for _, team := range teams {
		table.Append([]string{
			strconv.Itoa(team.TeamID),
			team.Name,
			team.Country,
			team.HomeRegion,
			strconv.Itoa(team.RookieYear),
		})
	}

	table.Render()
	return sb.String()
}

func RenderTeamMatchDetails(details []*query.TeamMatchDetails) string {
	var sb strings.Builder
	table := tablewriter.NewTable(&sb)
	table.Header([]string{"Event", "Match", "Alliance", "Total Points", "Teams"})

	for _, detail := range details {
		for _, mt := range detail.Teams {
			table.Append([]string{
				detail.Event.EventCode,
				strconv.Itoa(detail.Match.MatchNumber),
				detail.AllianceScore.Alliance,
				strconv.Itoa(detail.AllianceScore.TotalPoints),
				strconv.Itoa(mt.TeamID),
			})
		}
	}

	table.Render()
	return sb.String()
}
