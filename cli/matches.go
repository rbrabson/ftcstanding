package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/rbrabson/ftcstanding/query"
)

// RenderMatchDetails renders a list of MatchDetails in a formatted table.
func RenderMatchDetails(details []*query.MatchDetails) string {
	// colorCfg := renderer.ColorizedConfig{
	// 	Header: renderer.Tint{
	// 		FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
	// 	},
	// 	Column: renderer.Tint{
	// 		FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
	// 		Columns: []renderer.Tint{
	// 			{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0 (Match Type)
	// 			{FG: renderer.Colors{color.FgYellow}},  // Yellow for column 1 (Match Number)
	// 			{FG: renderer.Colors{color.FgRed}},     // Red for column 2 (Red Teams)
	// 			{FG: renderer.Colors{color.FgBlue}},    // Blue for column 3 (Blue Teams)
	// 			{FG: renderer.Colors{color.FgHiRed}},   // High-intensity red for column 4 (Red Score)
	// 			{FG: renderer.Colors{color.FgHiBlue}},  // High-intensity blue for column 5 (Blue Score)
	// 			{FG: renderer.Colors{color.FgHiCyan}},  // High-intensity cyan for column 7 (Winning Alliance)
	// 		},
	// 	},
	// 	Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
	// 	Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	// }

	// var sb strings.Builder
	// table := tablewriter.NewTable(&sb,
	// 	tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
	// )

	// TODO: trying some stuff out....
	var sb strings.Builder
	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Merging:   tw.CellMerging{Mode: tw.MergeHorizontal},
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{
				Merging:   tw.CellMerging{Mode: tw.MergeHierarchical},
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
		}),
	)

	table.Header([]string{"Type", "Match #", "Red Teams", "Red Teams", "Blue Teams", "Blue Teams", "Red Score", "Blue Score", "Winner"})

	for _, detail := range details {
		// Get red alliance teams
		redTeams := make([]string, 0, len(detail.RedAlliance.Teams))
		for _, team := range detail.RedAlliance.Teams {
			teamStr := fmt.Sprintf("%d\n%s", team.TeamID, team.Name)
			redTeams = append(redTeams, teamStr)
		}

		// Get blue alliance teams
		blueTeams := make([]string, 0, len(detail.BlueAlliance.Teams))
		for _, team := range detail.BlueAlliance.Teams {
			teamStr := fmt.Sprintf("%d\n%s", team.TeamID, team.Name)
			blueTeams = append(blueTeams, teamStr)
		}

		// Get scores
		var redPoints, bluePoints int
		redScore := "-"
		if detail.RedAlliance.Score != nil {
			redScore = strconv.Itoa(detail.RedAlliance.Score.TotalPoints)
			redPoints = detail.RedAlliance.Score.TotalPoints
		}

		blueScore := "-"
		if detail.BlueAlliance.Score != nil {
			blueScore = strconv.Itoa(detail.BlueAlliance.Score.TotalPoints)
			bluePoints = detail.BlueAlliance.Score.TotalPoints
		}

		var winner string
		switch {
		case redPoints > bluePoints:
			winner = "Red"
		case bluePoints > redPoints:
			winner = "Blue"
		default:
			winner = "Tie"
		}

		table.Append([]string{
			detail.Match.MatchType,
			strconv.Itoa(detail.Match.MatchNumber),
			redTeams[0],
			redTeams[1],
			blueTeams[0],
			blueTeams[1],
			redScore,
			blueScore,
			winner,
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
		"",
		"",
	})

	table.Render()
	return sb.String()
}
