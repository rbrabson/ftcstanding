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
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
		Settings:  tw.Settings{Separators: tw.Separators{BetweenRows: tw.Off}},
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
		greenColor := color.New(color.FgGreen)
		for _, ta := range report.TeamAdvancements {
			// Format team with advancement status
			teamName := fmt.Sprintf("%5d - %s", ta.Team.TeamID, ta.Team.Name)
			var advancementNumber string
			switch {
			case ta.Status == "already_advancing":
				teamName = fmt.Sprintf("%s\n        %s", teamName, greenColor.Sprint("(already advanced)"))
				advancementNumber = "-"
			case ta.AdvancementNumber != "-":
				advancementRank++
				advancementNumber = strconv.Itoa(advancementRank)
			default:
				advancementNumber = "-"
			}

			// Color advancing teams in green
			if ta.Advances {
				table.Append([]string{
					greenColor.Sprint(fmt.Sprintf("%d", ta.Rank)),
					greenColor.Sprint(teamName),
					fmt.Sprintf("%d", ta.TotalPoints),
					fmt.Sprintf("%d", ta.JudgingPoints),
					fmt.Sprintf("%d", ta.PlayoffPoints),
					fmt.Sprintf("%d", ta.SelectionPoints),
					fmt.Sprintf("%d", ta.QualificationPoints),
					advancementNumber,
				})
			} else {
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
		}
	}

	table.Render()
	return sb.String()
}

// RenderRegionAdvancementReport renders region-wide advancement information for all advancing teams.
// It shows each team's advancing event, awards from that event, and other events they participated in.
func RenderRegionAdvancementReport(report *query.RegionAdvancementReport) string {
	if report == nil {
		return "No region data available\n"
	}

	var sb strings.Builder

	// Render header
	sb.WriteString(color.New(color.FgGreen, color.Bold).Sprint("Region Advancement Report\n"))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Region: %s\n", report.RegionCode))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Year: %d\n", report.Year))
	sb.WriteString(color.New(color.FgCyan).Sprintf("Advancing Teams: %d\n\n", len(report.TeamAdvancements)))

	if len(report.TeamAdvancements) == 0 {
		sb.WriteString("No advancing teams found for this region.\n")
		return sb.String()
	}

	// Main table configuration
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgHiMagenta, color.Bold}}, // Yellow bold for team
				{FG: renderer.Colors{color.FgGreen}},                 // Green for advancing event
				{FG: renderer.Colors{color.FgCyan}},                  // Cyan for other events
			},
		},
		Footer:    renderer.Tint{FG: renderer.Colors{color.FgYellow, color.Bold}}, // Yellow bold footer
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}},              // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}},              // White separators
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
		}),
	)
	table.Header([]string{"Team", "Advancing Event", "Other Events"})

	// Populate table rows
	for _, ta := range report.TeamAdvancements {
		// Format team
		teamName := fmt.Sprintf("%d - %s", ta.Team.TeamID, ta.Team.Name)

		// Format advancing event with awards
		advancingEvent := fmt.Sprintf("• %s - %s", ta.AdvancingEvent.EventCode, ta.AdvancingEvent.Name)
		if len(ta.AdvancingEventAwards) > 0 {
			var awardsList []string
			for _, award := range ta.AdvancingEventAwards {
				awardsList = append(awardsList, fmt.Sprintf("  ◦ %s", award.Name))
			}
			advancingEvent += "\n" + strings.Join(awardsList, "\n")
		}

		// Format other events
		var otherEventsStr string
		if len(ta.OtherEventParticipations) > 0 {
			var eventsList []string
			for _, ep := range ta.OtherEventParticipations {
				eventLine := fmt.Sprintf("• %s - %s", ep.Event.EventCode, ep.Event.Name)

				// Add awards from this event
				if len(ep.Awards) > 0 {
					var awardsList []string
					for _, award := range ep.Awards {
						awardsList = append(awardsList, fmt.Sprintf("  ◦ %s", award.Name))
					}
					eventLine += "\n" + strings.Join(awardsList, "\n")
				}
				eventsList = append(eventsList, eventLine)
			}
			otherEventsStr = strings.Join(eventsList, "\n")
		} else {
			otherEventsStr = "None"
		}

		table.Append([]string{
			teamName,
			advancingEvent,
			otherEventsStr,
		})
	}

	// Add footer with team count
	table.Footer([]string{fmt.Sprintf("Teams Advancing: %d", len(report.TeamAdvancements)), "", ""})

	table.Render()
	return sb.String()
}
