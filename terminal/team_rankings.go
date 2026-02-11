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
	"github.com/rbrabson/ftcstanding/query"
)

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
// If limit is greater than 0, only the top 'limit' teams are displayed.
func RenderTeamPerformance(performances []query.TeamPerformance, eventCode string, sortBy SortBy, region string, year int, limit int) string {
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
			return performances[i].NpAVG > performances[j].NpAVG
		}
	})

	// Apply limit if specified
	if limit > 0 && limit < len(performances) {
		performances = performances[:limit]
	}

	var sb strings.Builder

	// Header
	sb.WriteString(color.HiCyanString("═══════════════════════════════════════════════════════════════\n"))
	if eventCode != "" {
		sb.WriteString(color.HiGreenString("Team Performance Rankings - %s (%d) - Event: %s\n", region, year, eventCode))
	} else {
		sb.WriteString(color.HiGreenString("Team Performance Rankings - %s (%d)\n", region, year))
	}
	sb.WriteString(color.HiYellowString("Sorted by: %s\n", sortBy))
	sb.WriteString(color.HiCyanString("═══════════════════════════════════════════════════════════════\n"))

	// Metric definitions
	sb.WriteString(color.HiWhiteString("\nMetric Definitions:\n\n"))

	sb.WriteString(color.HiYellowString("CCWM — Calculated Contribution to Winning Margin\n"))
	sb.WriteString(color.WhiteString("  Estimates how much a team affects the margin of victory or loss.\n"))
	sb.WriteString(color.WhiteString("  Positive CCWM → team usually helps alliances win by more\n"))
	sb.WriteString(color.WhiteString("  Negative CCWM → alliances with this team often lose by more\n"))
	sb.WriteString(color.HiCyanString("  👉 This blends offense, defense, and penalties into one \"do they help us win?\" number.\n\n"))

	sb.WriteString(color.HiYellowString("OPR — Offensive Power Rating\n"))
	sb.WriteString(color.WhiteString("  An estimate of how many points a team contributes per match to their alliance.\n"))
	sb.WriteString(color.WhiteString("  Calculated using math across all matches, factoring in partners and opponents.\n"))
	sb.WriteString(color.WhiteString("  Higher OPR = stronger overall scoring impact.\n"))
	sb.WriteString(color.HiCyanString("  👉 Think of it as: \"If this team plays, how many points do they add?\"\n\n"))

	sb.WriteString(color.HiYellowString("NP OPR — Non-Penalty Offensive Power Rating\n"))
	sb.WriteString(color.WhiteString("  Same idea as OPR, but penalties are removed.\n"))
	sb.WriteString(color.WhiteString("  Only counts points scored through gameplay, not points gained because opponents messed up.\n"))
	sb.WriteString(color.HiCyanString("  👉 Useful when you want to see true scoring ability, not \"we won because the other\n"))
	sb.WriteString(color.HiCyanString("     alliance kept getting penalties.\"\n\n"))

	sb.WriteString(color.HiYellowString("DPR — Defensive Power Rating\n"))
	sb.WriteString(color.WhiteString("  Estimates how many points a team allows opponents to score.\n"))
	sb.WriteString(color.WhiteString("  Lower DPR = better defense.\n"))
	sb.WriteString(color.WhiteString("  A strong defensive robot often has a noticeably low DPR even if OPR isn't huge.\n"))
	sb.WriteString(color.HiCyanString("  👉 Think of it as: \"If this team plays, how well do they keep the opponents from scoring?\"\n\n"))

	sb.WriteString(color.HiYellowString("NP DPR — Non-Penalty Defensive Power Rating\n"))
	sb.WriteString(color.WhiteString("  Same as DPR, but ignores penalty points.\n"))
	sb.WriteString(color.WhiteString("  Focuses only on how well a team limits actual scoring, not ref calls.\n"))
	sb.WriteString(color.HiCyanString("  👉 Great for identifying clean, effective defense.\n\n"))

	sb.WriteString(color.HiYellowString("NP AVG — Non-Penalty Average Score\n"))
	sb.WriteString(color.WhiteString("  The average number of non-penalty points a team's alliance scores in matches involving them.\n"))
	sb.WriteString(color.WhiteString("  Subtracts the penalties commited by the team's alliance to determine the true scoring contribution.\n"))
	sb.WriteString(color.WhiteString("  Less math-heavy than OPR, more literal.\n"))
	sb.WriteString(color.WhiteString("  Still partner-dependent, but easier to interpret.\n"))
	sb.WriteString(color.HiCyanString("  👉 Think: \"On average, when this team plays, how many real points get scored?\"\n\n"))

	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold},
			BG: renderer.Colors{color.BgBlack},
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan},
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}},   // Rank
				{FG: renderer.Colors{color.FgHiWhite}},   // Team
				{FG: renderer.Colors{color.FgHiCyan}},    // Region
				{FG: renderer.Colors{color.FgHiRed}},     // Matches
				{FG: renderer.Colors{color.FgHiMagenta}}, // CCWM
				{FG: renderer.Colors{color.FgHiGreen}},   // OPR
				{FG: renderer.Colors{color.FgHiGreen}},   // npOPR
				{FG: renderer.Colors{color.FgHiYellow}},  // DPR
				{FG: renderer.Colors{color.FgHiYellow}},  // npDPR
				{FG: renderer.Colors{color.FgHiMagenta}}, // npAVG
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}},
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}},
	}

	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,   // Rank
					tw.AlignLeft,   // Team
					tw.AlignLeft,   // Region
					tw.AlignCenter, // Matches
					tw.AlignCenter, // CCWM
					tw.AlignCenter, // OPR
					tw.AlignCenter, // npOPR
					tw.AlignCenter, // DPR
					tw.AlignCenter, // npDPR
					tw.AlignCenter, // npAVG
				}},
			},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{PerColumn: []tw.Align{
					tw.AlignLeft,  // Rank
					tw.AlignLeft,  // Team
					tw.AlignLeft,  // Region
					tw.AlignRight, // Matches
					tw.AlignRight, // CCWM
					tw.AlignRight, // OPR
					tw.AlignRight, // npOPR
					tw.AlignRight, // DPR
					tw.AlignRight, // npDPR
					tw.AlignRight, // npAVG
				}},
			},
		}),
	)

	table.Header([]string{"Rank", "Team", "Region", "Matches", "CCWM", "OPR", "npOPR", "DPR", "npDPR", "npAVG"})

	for i, perf := range performances {
		table.Append([]string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%5d - %s", perf.TeamID, perf.TeamName),
			perf.Region,
			strconv.Itoa(perf.Matches),
			fmt.Sprintf("%.2f", perf.CCWM),
			fmt.Sprintf("%.2f", perf.OPR),
			fmt.Sprintf("%.2f", perf.NpOPR),
			fmt.Sprintf("%.2f", perf.DPR),
			fmt.Sprintf("%.2f", perf.NpDPR),
			fmt.Sprintf("%.2f", perf.NpAVG),
		})
	}

	table.Render()

	return sb.String()
}
