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

// RenderAdvancementReport renders event details and all team advancement information in a formatted table.
func RenderAdvancementReport(report *query.AdvancementReport) string {
	if report == nil || report.Event == nil {
		return "No event data available\n"
	}

	var sb strings.Builder

	// Render event information header
	sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Event Advancement Report\n"))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Code: %s\n", report.Event.EventCode))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Name: %s\n", report.Event.Name))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", report.Event.Year))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Location: %s, %s, %s\n\n",
		report.Event.City, report.Event.StateProv, report.Event.Country))

	// Render advancement table
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta, color.Bold}}, // Magenta bold for rank
				{FG: renderer.Colors{color.FgYellow}},              // Yellow for team
				{FG: renderer.Colors{color.FgCyan, color.Bold}},    // Cyan bold for total
				{}, // Default for remaining columns
			},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow bold footer
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
		Settings:  tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)
	table.Header([]string{"Rank", "Team", "Total Pts", "Judging", "Playoff", "Selection", "Qualification", "Adv #"})

	if len(report.TeamAdvancements) == 0 {
		sb.WriteString("\nNo teams found for this event.\n")
	} else {
		var advancementRank int
		for _, ta := range report.TeamAdvancements {
			// Format team with advancement status
			teamName := fmt.Sprintf("%5d - %s", ta.Team.TeamID, ta.Team.Name)
			var advancementNumber string
			switch {
			case ta.Status == "already_advancing":
				teamName = fmt.Sprintf("%s\n        (already advanced)", teamName)
				advancementNumber = "-"
			case ta.AdvancementNumber != "-":
				advancementRank++
				advancementNumber = strconv.Itoa(advancementRank)
			default:
				advancementNumber = "-"
			}

			table.Append([]string{
				fmt.Sprintf("%d", ta.Rank),
				teamName,
				fmt.Sprintf("%d", ta.TotalPoints),
				fmt.Sprintf("%d", ta.JudgingPoints),
				fmt.Sprintf("%d", ta.PlayoffPoints),
				fmt.Sprintf("%d", ta.SelectionPoints),
				fmt.Sprintf("%d", ta.QualificationPoints),
				advancementNumber,
			})
		}

		// Add footer with team count
		table.Footer([]string{
			fmt.Sprintf("Total Teams: %d", len(report.TeamAdvancements)),
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		})
	}

	table.Render()
	return sb.String()
}
