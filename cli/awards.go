package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/rbrabson/ftcstanding/query"
)

// RenderAwardsByEvent renders event details and all awards won by teams in a formatted table.
func RenderAwardsByEvent(eventAwards *query.EventAwards) string {
	if eventAwards == nil || eventAwards.Event == nil {
		return "No event data available\n"
	}

	var sb strings.Builder

	// Render event information header
	sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Event Awards\n"))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Code: %s\n", eventAwards.Event.EventCode))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Name: %s\n", eventAwards.Event.Name))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", eventAwards.Event.Year))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Location: %s, %s, %s\n\n",
		eventAwards.Event.City, eventAwards.Event.StateProv, eventAwards.Event.Country))

	// Render awards table
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgYellow}},  // Yellow for column 0 (Award Name)
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 1 (Winner)
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
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
			Widths: tw.CellWidth{
				PerColumn: tw.Mapper[int, int]{
					1: 80, // Winner column max width
				},
			},
		}),
	)
	table.Header([]string{"Award Name", "Winner"})

	if len(eventAwards.Awards) == 0 {
		sb.WriteString("\nNo awards found for this event.\n")
	} else {
		for _, teamAward := range eventAwards.Awards {
			// Format award name, qualified by series if needed
			awardName := teamAward.Award.Name

			// Format winner as "TeamID - Team Name" with full name on line below
			// winner := fmt.Sprintf("%d - %s\n%s", teamAward.Team.TeamID, teamAward.Team.Name, teamAward.Team.FullName)
			winner := fmt.Sprintf("%d - %s", teamAward.Team.TeamID, teamAward.Team.Name)

			table.Append([]string{
				awardName,
				winner,
			})
		}

		// Add footer with award count
		table.Footer([]string{
			fmt.Sprintf("Total: %d", len(eventAwards.Awards)),
			"",
		})
	}

	table.Render()
	return sb.String()
}
