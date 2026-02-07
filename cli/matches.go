package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/rbrabson/ftcstanding/query"
)

// RenderMatchDetails renders a list of MatchDetails in a formatted table.
func RenderMatchDetails(details []*query.MatchDetails) string {
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0 (Match Type)
				{FG: renderer.Colors{color.FgYellow}},  // Yellow for column 1 (Match Number)
				{FG: renderer.Colors{color.FgRed}},     // Red for column 2 (Red Teams)
				{FG: renderer.Colors{color.FgBlue}},    // Blue for column 3 (Blue Teams)
				{FG: renderer.Colors{color.FgHiRed}},   // High-intensity red for column 4 (Red Score)
				{FG: renderer.Colors{color.FgHiBlue}},  // High-intensity blue for column 5 (Blue Score)
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	}

	var sb strings.Builder
	table := tablewriter.NewTable(&sb, tablewriter.WithRenderer(renderer.NewColorized(colorCfg)))
	table.Header([]string{"Type", "Match #", "Red Teams", "Blue Teams", "Red Score", "Blue Score"})

	for _, detail := range details {
		// Get red alliance teams
		redTeams := make([]string, 0, len(detail.RedAlliance.Teams))
		for _, team := range detail.RedAlliance.Teams {
			redTeams = append(redTeams, strconv.Itoa(team.TeamID))
		}
		redTeamsStr := strings.Join(redTeams, ", ")

		// Get blue alliance teams
		blueTeams := make([]string, 0, len(detail.BlueAlliance.Teams))
		for _, team := range detail.BlueAlliance.Teams {
			blueTeams = append(blueTeams, strconv.Itoa(team.TeamID))
		}
		blueTeamsStr := strings.Join(blueTeams, ", ")

		// Get scores
		redScore := "-"
		if detail.RedAlliance.Score != nil {
			redScore = strconv.Itoa(detail.RedAlliance.Score.TotalPoints)
		}

		blueScore := "-"
		if detail.BlueAlliance.Score != nil {
			blueScore = strconv.Itoa(detail.BlueAlliance.Score.TotalPoints)
		}

		table.Append([]string{
			detail.Match.MatchType,
			strconv.Itoa(detail.Match.MatchNumber),
			redTeamsStr,
			blueTeamsStr,
			redScore,
			blueScore,
		})
	}

	// Add footer with match count
	table.Footer([]string{
		"Total Matches",
		fmt.Sprintf("%d", len(details)),
		"",
		"",
		"",
		"",
	})

	table.Render()
	return sb.String()
}
