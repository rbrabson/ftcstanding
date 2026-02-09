package terminal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
)

// RenderTeams renders a list of teams in a table format.
func RenderTeams(teams []*database.Team) string {
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
			BG: renderer.Colors{color.BgBlack},             // White background
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
	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
				}},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
					tw.AlignLeft,
				}},
			},
		}),
	)
	table.Header([]string{"Team", "Country", "Region", "Rookie Year"})

	for _, team := range teams {
		table.Append([]string{
			strconv.Itoa(team.TeamID) + " - " + team.Name,
			team.Country,
			team.HomeRegion,
			strconv.Itoa(team.RookieYear),
		})
	}

	table.Footer([]string{"Total Teams: " + strconv.Itoa(len(teams)), "", "", ""})

	table.Render()
	return sb.String()
}

// formatRecord formats a Record as a W-L-T string.
func formatRecord(r query.Record) string {
	return fmt.Sprintf("%d-%d-%d", r.Wins, r.Losses, r.Ties)
}

// RenderTeamDetails renders detailed information about a team including events, records, and awards.
func RenderTeamDetails(details *query.TeamDetails) string {
	if details == nil {
		return "No team details available\n"
	}

	var sb strings.Builder

	// Team Header Information
	sb.WriteString(color.HiCyanString("═══════════════════════════════════════════════════════════════\n"))
	sb.WriteString(color.HiGreenString("Team %d - %s\n", details.TeamID, details.Name))
	sb.WriteString(color.HiCyanString("═══════════════════════════════════════════════════════════════\n"))
	if details.FullName != "" {
		sb.WriteString(color.WhiteString("Details:  %s\n", details.FullName))
	}
	sb.WriteString(color.WhiteString("Location: %s, %s, %s\n", details.City, details.StateProv, details.Country))
	sb.WriteString(color.WhiteString("Region:   %s\n", details.Region))
	sb.WriteString("\n")

	// Overall Records
	sb.WriteString(color.YellowString("Overall Record:\n"))
	sb.WriteString(color.WhiteString("  Total:         %s\n", formatRecord(details.TotalRecord)))
	sb.WriteString(color.WhiteString("  Qualification: %s\n", formatRecord(details.QualRecord)))
	sb.WriteString(color.WhiteString("  Playoff:       %s\n", formatRecord(details.PlayoffRecord)))
	sb.WriteString("\n")

	// Events Table
	if len(details.Events) > 0 {
		sb.WriteString(color.YellowString("Events:\n"))

		colorCfg := renderer.ColorizedConfig{
			Header: renderer.Tint{
				FG: renderer.Colors{color.FgGreen, color.Bold},
				BG: renderer.Colors{color.BgBlack},
			},
			Column: renderer.Tint{
				FG: renderer.Colors{color.FgCyan},
				Columns: []renderer.Tint{
					{FG: renderer.Colors{color.FgMagenta}}, // Event Code
					{FG: renderer.Colors{color.FgWhite}},   // Event Name
					{},                                     // Total Record
					{},                                     // Qual Record
					{},                                     // Playoff Record
					{FG: renderer.Colors{color.FgHiGreen}}, // Advanced
					{FG: renderer.Colors{color.FgYellow}},  // Awards
				},
			},
			Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}},
			Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}},
		}

		var tableSb strings.Builder
		table := tablewriter.NewTable(&tableSb,
			tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
			tablewriter.WithConfig(tablewriter.Config{
				Header: tw.CellConfig{
					Alignment: tw.CellAlignment{PerColumn: []tw.Align{
						tw.AlignLeft,   // Event Code
						tw.AlignLeft,   // Event Name
						tw.AlignCenter, // Total Record
						tw.AlignCenter, // Qual Record
						tw.AlignCenter, // Playoff Record
						tw.AlignCenter, // Advanced
						tw.AlignLeft,   // Awards
					}},
				},
			}),
		)

		table.Header([]string{"Event Code", "Event Name", "Total", "Qual", "Playoff", "Advanced", "Awards"})

		for _, event := range details.Events {
			advancedStr := ""
			if event.Advanced {
				advancedStr = "✓"
			}

			awardsStr := ""
			if len(event.Awards) > 0 {
				awardsStr = strings.Join(event.Awards, ", ")
			}

			table.Append([]string{
				event.EventCode,
				event.EventName,
				formatRecord(event.TotalRecord),
				formatRecord(event.QualRecord),
				formatRecord(event.PlayoffRecord),
				advancedStr,
				awardsStr,
			})
		}

		table.Render()
		sb.WriteString(tableSb.String())
	} else {
		sb.WriteString(color.YellowString("No events found for this team.\n"))
	}

	return sb.String()
}
