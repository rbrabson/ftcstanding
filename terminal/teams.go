package terminal

import (
	"fmt"
	"sort"
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
		parts := strings.Split(details.FullName, "&")
		for i := range parts {
			index := len(parts) - i - 1
			part := parts[index]
			part = strings.ReplaceAll(parts[index], "/", ", ")
			if i == 0 {
				sb.WriteString(color.WhiteString("Details:  %s\n", part))
			} else {
				sb.WriteString(color.WhiteString("          %s\n", part))
			}
		}
	}
	sb.WriteString(color.WhiteString("Location: %s, %s, %s\n", details.City, details.StateProv, details.Country))
	if details.RookieYear > 0 {
		sb.WriteString(color.WhiteString("Rookie:   %d\n", details.RookieYear))
	}
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
					{FG: renderer.Colors{color.FgWhite}},   // Event Name				{},                                     // Qual Rank					{},                                     // Total Record
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
						tw.AlignLeft,   // Event Name					tw.AlignCenter, // Qual Rank						tw.AlignCenter, // Total Record
						tw.AlignCenter, // Qual Record
						tw.AlignCenter, // Playoff Record
						tw.AlignCenter, // Advanced
						tw.AlignLeft,   // Awards
					}},
				},
			}),
		)

		table.Header([]string{"Event Code", "Event Name", "Rank", "Total", "Qual", "Playoff", "Advanced", "Awards"})

		for _, event := range details.Events {
			advancedStr := ""
			if event.Advanced {
				advancedStr = "✓"
			}

			awardsStr := ""
			if len(event.Awards) > 0 {
				awardsStr = strings.Join(event.Awards, ", ")
			}

			rankStr := ""
			if event.QualRank > 0 {
				rankStr = strconv.Itoa(event.QualRank)
			}

			table.Append([]string{
				event.EventCode,
				event.EventName,
				rankStr,
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

// SortBy defines the sort criteria for team performance
type SortBy string

const (
	SortByOPR     SortBy = "opr"
	SortByNpOPR   SortBy = "npopr"
	SortByCCWM    SortBy = "ccwm"
	SortByDPR     SortBy = "dpr"
	SortByNpDPR   SortBy = "npdpr"
	SortByNpAVG   SortBy = "npavg"
	SortByMatches SortBy = "matches"
	SortByTeamID  SortBy = "team"
)

// RenderTeamPerformance renders team performance metrics in a table format with sorting.
func RenderTeamPerformance(performances []query.TeamPerformance, sortBy SortBy, region string, year int) string {
	if len(performances) == 0 {
		return color.YellowString("No performance data available for region %s in year %d\n", region, year)
	}

	// Sort the performances based on the specified criteria
	sort.Slice(performances, func(i, j int) bool {
		switch sortBy {
		case SortByOPR:
			return performances[i].OPR > performances[j].OPR
		case SortByNpOPR:
			return performances[i].NpOPR > performances[j].NpOPR
		case SortByCCWM:
			return performances[i].CCWM > performances[j].CCWM
		case SortByDPR:
			return performances[i].DPR < performances[j].DPR // Lower is better for defense
		case SortByNpDPR:
			return performances[i].NpDPR < performances[j].NpDPR // Lower is better for defense
		case SortByNpAVG:
			return performances[i].NpAVG > performances[j].NpAVG
		case SortByMatches:
			return performances[i].Matches > performances[j].Matches
		case SortByTeamID:
			return performances[i].TeamID < performances[j].TeamID
		default:
			return performances[i].OPR > performances[j].OPR
		}
	})

	var sb strings.Builder

	// Header
	sb.WriteString(color.HiCyanString("═══════════════════════════════════════════════════════════════\n"))
	sb.WriteString(color.HiGreenString("Team Performance Rankings - %s (%d)\n", region, year))
	sb.WriteString(color.HiYellowString("Sorted by: %s\n", sortBy))
	sb.WriteString(color.HiCyanString("═══════════════════════════════════════════════════════════════\n\n"))

	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold},
			BG: renderer.Colors{color.BgBlack},
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}},
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}},
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignRight, // Rank
					tw.AlignLeft,  // Team
					tw.AlignRight, // Matches
					tw.AlignRight, // OPR
					tw.AlignRight, // npOPR
					tw.AlignRight, // CCWM
					tw.AlignRight, // DPR
					tw.AlignRight, // npDPR
					tw.AlignRight, // npAVG
				}},
			},
		}),
	)

	table.Header([]string{"Rank", "Team", "Matches", "OPR", "npOPR", "CCWM", "DPR", "npDPR", "npAVG"})

	for i, perf := range performances {
		table.Append([]string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%d - %s", perf.TeamID, perf.TeamName),
			strconv.Itoa(perf.Matches),
			fmt.Sprintf("%.2f", perf.OPR),
			fmt.Sprintf("%.2f", perf.NpOPR),
			fmt.Sprintf("%.2f", perf.CCWM),
			fmt.Sprintf("%.2f", perf.DPR),
			fmt.Sprintf("%.2f", perf.NpDPR),
			fmt.Sprintf("%.2f", perf.NpAVG),
		})
	}

	table.Render()

	return sb.String()
}
