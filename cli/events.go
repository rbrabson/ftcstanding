package cli

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

// RenderTeamsByEvent renders event details and all participating teams in a formatted table.
func RenderTeamsByEvent(eventTeams *query.EventTeams) string {
	if eventTeams == nil || eventTeams.Event == nil {
		return "No event data available\n"
	}

	var sb strings.Builder

	// Render event information header
	sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Event Information\n"))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Code: %s\n", eventTeams.Event.EventCode))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Name: %s\n", eventTeams.Event.Name))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", eventTeams.Event.Year))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Location: %s, %s, %s\n",
		eventTeams.Event.City, eventTeams.Event.StateProv, eventTeams.Event.Country))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Dates: %s to %s\n\n",
		eventTeams.Event.DateStart.Format("Jan 2, 2006"),
		eventTeams.Event.DateEnd.Format("Jan 2, 2006")))

	// Render teams table
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0 (Team Number)
				{},                                     // Inherit default (cyan) for column 1 (Team Name)
				{FG: renderer.Colors{color.FgHiRed}},   // High-intensity red for column 2 (Location)
				{},                                     // Inherit default (cyan) for remaining columns
			},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow bold footer
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)
	table.Header([]string{"Number", "Name", "Location", "Region", "Rookie Year"})

	if len(eventTeams.Teams) == 0 {
		sb.WriteString("\nNo teams found for this event.\n")
	} else {
		for _, team := range eventTeams.Teams {
			location := fmt.Sprintf("%s, %s, %s", team.City, team.StateProv, team.Country)
			table.Append([]string{
				strconv.Itoa(team.TeamID),
				team.Name,
				location,
				team.HomeRegion,
				strconv.Itoa(team.RookieYear),
			})
		}

		// Add footer with team count
		table.Footer([]string{
			fmt.Sprintf("Total Teams: %d", len(eventTeams.Teams)),
			"",
			"",
			"",
			"",
		})

		table.Render()
	}

	return sb.String()
}

// RenderTeamRankings renders event details and team rankings in a formatted table.
func RenderTeamRankings(eventRankings *query.EventTeamRankings) string {
	if eventRankings == nil || eventRankings.Event == nil {
		return "No event data available\n"
	}

	var sb strings.Builder

	// Render event information header
	sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Event Rankings\n"))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Code: %s\n", eventRankings.Event.EventCode))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Name: %s\n", eventRankings.Event.Name))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", eventRankings.Event.Year))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Location: %s, %s, %s\n",
		eventRankings.Event.City, eventRankings.Event.StateProv, eventRankings.Event.Country))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Dates: %s to %s\n\n",
		eventRankings.Event.DateStart.Format("Jan 2, 2006"),
		eventRankings.Event.DateEnd.Format("Jan 2, 2006")))

	// Render rankings table
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta, color.Bold}}, // Magenta for column 0 (Rank)
				{FG: renderer.Colors{color.FgYellow}},              // Yellow for column 1 (Team)
				{},                                                 // Default for remaining columns
			},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow bold footer
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)
	table.Header([]string{"Rank", "Team", "RS", "Match Pts", "Base Pts", "Auto Pts", "High Score", "W-L-T", "Matches"})

	if len(eventRankings.TeamRankings) == 0 {
		sb.WriteString("\nNo rankings found for this event.\n")
	} else {
		for _, tr := range eventRankings.TeamRankings {
			team := fmt.Sprintf("%d - %s", tr.Team.TeamID, tr.Team.Name)
			wlt := fmt.Sprintf("%d-%d-%d", tr.Ranking.Wins, tr.Ranking.Losses, tr.Ranking.Ties)

			table.Append([]string{
				strconv.Itoa(tr.Ranking.Rank),
				team,
				fmt.Sprintf("%.2f", tr.Ranking.SortOrder1),
				fmt.Sprintf("%.2f", tr.Ranking.SortOrder2),
				fmt.Sprintf("%.2f", tr.Ranking.SortOrder3),
				fmt.Sprintf("%.2f", tr.Ranking.SortOrder4),
				strconv.Itoa(tr.HighMatchScore),
				wlt,
				strconv.Itoa(tr.Ranking.MatchesPlayed),
			})
		}

		// Add footer with team count
		table.Footer([]string{
			fmt.Sprintf("Total Teams: %d", len(eventRankings.TeamRankings)),
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		})

		table.Render()
	}

	return sb.String()
}
