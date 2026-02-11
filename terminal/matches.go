package terminal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/rbrabson/ftcstanding/query"
)

// RenderMatchDetails renders a list of MatchDetails in a formatted table.
func RenderMatchDetails(details []*query.MatchDetails) string {
	var sb strings.Builder

	// Render event information header
	if len(details) > 0 {
		event := details[0].Event
		sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Event Information\n"))
		sb.WriteString(color.New(color.FgCyan).Sprintf("Code: %s\n", event.EventCode))
		sb.WriteString(color.New(color.FgCyan).Sprintf("Name: %s\n", event.Name))
		sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", event.Year))
		sb.WriteString(color.New(color.FgCyan).Sprintf("Location: %s, %s, %s\n",
			event.City, event.StateProv, event.Country))
		sb.WriteString(color.New(color.FgCyan).Sprintf("Dates: %s to %s\n\n",
			event.DateStart.Format("Jan 2, 2006"),
			event.DateEnd.Format("Jan 2, 2006")))
	}

	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0 (Match Type)
				{FG: renderer.Colors{color.FgYellow}},  // Yellow for column 1 (Match Number)
				{FG: renderer.Colors{color.FgHiRed}},   // Red for column 2 (Red Team 1)
				{FG: renderer.Colors{color.FgHiRed}},   // Red for column 3 (Red Team 2)
				{FG: renderer.Colors{color.FgHiBlue}},  // Blue for column 4 (Blue Team 1)
				{FG: renderer.Colors{color.FgHiBlue}},  // Blue for column 5 (Blue Team 2)
				{},                                     // Default for column 6 (Scores - colors applied inline)
				{},                                     // Default for column 7 (Winning Alliance)
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
		Settings:  tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Merging: tw.CellMerging{Mode: tw.MergeHorizontal},
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,   // Type - left aligned
					tw.AlignCenter, // Match #
					tw.AlignCenter, // Red Alliance
					tw.AlignCenter, // Red Alliance
					tw.AlignCenter, // Blue Alliance
					tw.AlignCenter, // Blue Alliance
					tw.AlignCenter, // Scores
					tw.AlignCenter, // Winner
				}},
			},
			Row: tw.CellConfig{
				Merging: tw.CellMerging{Mode: tw.MergeHierarchical},
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,
					tw.AlignRight,
					tw.AlignCenter,
					tw.AlignCenter,
					tw.AlignCenter,
					tw.AlignCenter,
					tw.AlignCenter,
					tw.AlignCenter,
				}},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,  // Total Matches label
					tw.AlignRight, // Match count number
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
				}},
			},
		}),
	)

	table.Header([]string{"Type", "Match #", "Red Alliance", "Red Alliance", "Blue Alliance", "Blue Alliance", "Scores", "Winner"})

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

		// Combine scores with color coding (red first, then blue)
		blueScoreColored := color.New(color.FgHiBlue, color.Bold).Sprint(blueScore)
		redScoreColored := color.New(color.FgHiRed, color.Bold).Sprint(redScore)
		combinedScores := fmt.Sprintf("%s\n%s", redScoreColored, blueScoreColored)

		var winner string
		switch {
		case redPoints > bluePoints:
			winner = color.New(color.FgRed, color.Bold).Sprint("Red")
		case bluePoints > redPoints:
			winner = color.New(color.FgBlue, color.Bold).Sprint("Blue")
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
			combinedScores,
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

// RenderMatchesByEventAndTeam renders a list of TeamMatchResult in a formatted table.
func RenderMatchesByEventAndTeam(results []*query.TeamMatchResult) string {
	var sb strings.Builder

	if len(results) == 0 {
		return "No matches found for this team at this event.\n"
	}

	event := results[0].Event
	team := results[0].Team
	sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Event Information\n"))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Team: %d - %s\n", team.TeamID, team.Name))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Code: %s\n", event.EventCode))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Name: %s\n", event.Name))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Location: %s, %s, %s\n",
		event.City, event.StateProv, event.Country))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", event.Year))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Dates: %s to %s\n\n",
		event.DateStart.Format("Jan 2, 2006"),
		event.DateEnd.Format("Jan 2, 2006")))

	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0 (Match Type)
				{FG: renderer.Colors{color.FgYellow}},  // Yellow for column 1 (Match Number)
				{},                                     // Default for column 2 (Team Alliance 1 - colors applied inline)
				{},                                     // Default for column 3 (Team Alliance 2 - colors applied inline)
				{},                                     // Default for column 4 (Opponent Alliance 1 - colors applied inline)
				{},                                     // Default for column 5 (Opponent Alliance 2 - colors applied inline)
				{},                                     // Default for column 6 (Scores - colors applied inline)
				{},                                     // Default for column 7 (Result)
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
		Settings:  tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Merging: tw.CellMerging{Mode: tw.MergeHorizontal},
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,   // Type
					tw.AlignCenter, // Match #
					tw.AlignCenter, // Team Alliance
					tw.AlignCenter, // Team Alliance
					tw.AlignCenter, // Opponent Alliance
					tw.AlignCenter, // Opponent Alliance
					tw.AlignCenter, // Scores
					tw.AlignCenter, // Result
				}},
			},
			Row: tw.CellConfig{
				Merging: tw.CellMerging{Mode: tw.MergeHierarchical},
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,   // Type
					tw.AlignRight,  // Match #
					tw.AlignCenter, // Team Alliance
					tw.AlignCenter, // Team Alliance
					tw.AlignCenter, // Opponent Alliance
					tw.AlignCenter, // Opponent Alliance
					tw.AlignCenter, // Scores
					tw.AlignCenter, // Result
				}},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)

	table.Header([]string{"Type", "Match #", "Team Alliance", "Team Alliance", "Opponent Alliance", "Opponent Alliance", "Scores", "Result"})

	for _, result := range results {
		// Get team alliance teams with coloring based on alliance
		teamAllianceColor := color.FgHiRed
		if result.TeamAlliance.Alliance == "blue" {
			teamAllianceColor = color.FgHiBlue
		}
		teamTeams := make([]string, 0, len(result.TeamAlliance.Teams))
		for _, team := range result.TeamAlliance.Teams {
			teamStr := fmt.Sprintf("%s\n%s", color.New(teamAllianceColor).Sprintf("%d", team.TeamID), color.New(teamAllianceColor).Sprintf("%s", team.Name))
			teamTeams = append(teamTeams, teamStr)
		}

		// Get opponent alliance teams with coloring based on alliance
		opponentAllianceColor := color.FgHiRed
		if result.OpponentAlliance.Alliance == "blue" {
			opponentAllianceColor = color.FgHiBlue
		}
		opponentTeams := make([]string, 0, len(result.OpponentAlliance.Teams))
		for _, team := range result.OpponentAlliance.Teams {
			teamStr := fmt.Sprintf("%s\n%s", color.New(opponentAllianceColor).Sprintf("%d", team.TeamID), color.New(opponentAllianceColor).Sprintf("%s", team.Name))
			opponentTeams = append(opponentTeams, teamStr)
		}

		// Get scores
		teamScore := "-"
		if result.TeamAlliance.Score != nil {
			teamScore = strconv.Itoa(result.TeamAlliance.Score.TotalPoints)
		}

		opponentScore := "-"
		if result.OpponentAlliance.Score != nil {
			opponentScore = strconv.Itoa(result.OpponentAlliance.Score.TotalPoints)
		}

		// Combine scores with color coding (team first, then opponent)
		teamScoreColored := color.New(teamAllianceColor, color.Bold).Sprint(teamScore)
		opponentScoreColored := color.New(opponentAllianceColor, color.Bold).Sprint(opponentScore)
		combinedScores := fmt.Sprintf("%s\n%s", teamScoreColored, opponentScoreColored)

		// Color the result based on outcome
		var resultColored string
		switch result.Result {
		case "Won":
			resultColored = color.New(color.FgGreen, color.Bold).Sprint("Won")
		case "Lost":
			resultColored = color.New(color.FgRed, color.Bold).Sprint("Lost")
		default:
			resultColored = color.New(color.FgYellow, color.Bold).Sprint("Tied")
		}

		table.Append([]string{
			result.Match.MatchType,
			strconv.Itoa(result.Match.MatchNumber),
			teamTeams[0],
			teamTeams[1],
			opponentTeams[0],
			opponentTeams[1],
			combinedScores,
			resultColored,
		})
	}

	// Add footer with match count
	table.Footer([]string{
		"Total Matches",
		fmt.Sprintf("%d", len(results)),
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
