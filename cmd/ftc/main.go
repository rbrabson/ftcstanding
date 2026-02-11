package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rbrabson/ftcstanding/request"
	"github.com/rbrabson/ftcstanding/terminal"
	"github.com/spf13/cobra"
)

var (
	defaultYear int
)

// setLogLevelFromEnv sets the log level from the LOG_LEVEL environment variable.
func setLogLevelFromEnv() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")

	var logLevel slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
	return logLevel
}

// initializeApp sets up database and initializes subsystems
func initializeApp() error {
	season := os.Getenv("FTC_SEASON")
	if season == "" {
		return fmt.Errorf("FTC_SEASON environment variable not set")
	}

	var err error
	defaultYear, err = strconv.Atoi(season)
	if err != nil {
		return fmt.Errorf("invalid FTC_SEASON value: %s", season)
	}

	db, err := database.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	request.Init(db)
	query.Init(db)
	terminal.Init(db)

	return nil
}

// rootCmd is the base command for the CLI application.
var rootCmd = &cobra.Command{
	Use:   "ftc",
	Short: "FTC Standing - A CLI tool for FTC competition data",
	Long:  `A command-line interface for querying and displaying FTC (FIRST Tech Challenge) competition data including teams, events, matches, rankings, and advancement information.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeApp()
	},
}

// teamCmd enders the advancement report for a specific event, showing which teams advanced and their points breakdown.
var teamCmd = &cobra.Command{
	Use:   "team [teamID]",
	Short: "Show detailed information about a team",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid teamID '%s', must be a number", args[0])
		}
		details := query.TeamDetailsQuery(teamID)
		if details == nil {
			return fmt.Errorf("team %d not found", teamID)
		}
		output := terminal.RenderTeamDetails(details)
		fmt.Println(output)
		return nil
	},
}

// teamsCmd lists all teams in a specified region, showing their team ID, name, and home region.
var teamsCmd = &cobra.Command{
	Use:   "teams [region]",
	Short: "List teams in a region",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		region := args[0]
		teamsFilter := database.TeamFilter{
			HomeRegions: []string{region},
		}
		teams := query.TeamsQuery(teamsFilter)
		teamsOutput := terminal.RenderTeams(teams)
		fmt.Println(teamsOutput)
		return nil
	},
}

// eventTeamsCmd lists all teams that participated in a specific event, showing their team ID, name, and home region.
var eventTeamsCmd = &cobra.Command{
	Use:   "event-teams [eventCode]",
	Short: "List teams at an event",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventCode := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		eventTeams := query.TeamsByEventQuery(eventCode, year)
		eventTeamsOutput := terminal.RenderTeamsByEvent(eventTeams)
		fmt.Println(eventTeamsOutput)
		return nil
	},
}

// rankingsCmd renders the team rankings at a specific event, showing each team's rank, name, points breakdown,
// and advancement status.
var rankingsCmd = &cobra.Command{
	Use:   "rankings [eventCode]",
	Short: "List team rankings at an event",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventCode := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		rankings := query.EventTeamRankingQuery(eventCode, year)
		teamRankingsOutput := terminal.RenderTeamRankings(rankings)
		fmt.Println(teamRankingsOutput)
		return nil
	},
}

// advancementCmd renders the advancement report for a specific event, showing which teams advanced
// and their points breakdown.
var awardsCmd = &cobra.Command{
	Use:   "awards [eventCode]",
	Short: "List award winners at an event",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventCode := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		awardsResults := query.AwardsByEventQuery(eventCode, year)
		awardResultsOutput := terminal.RenderAwardsByEvent(awardsResults)
		fmt.Println(awardResultsOutput)
		return nil
	},
}

// advancementCmd renders the advancement report for a specific event, showing which teams advanced
// and their points breakdown.
var advancementCmd = &cobra.Command{
	Use:   "advancement [eventCode]",
	Short: "Show advancement report for an event",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventCode := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		advancementReport := query.AdvancementReportQuery(eventCode, year)
		advancementReportOutput := terminal.RenderAdvancementReport(advancementReport)
		fmt.Println(advancementReportOutput)
		return nil
	},
}

// matchesCmd renders the match results for a specific event, showing each match's teams, scores,
// and outcomes.
var matchesCmd = &cobra.Command{
	Use:   "matches [eventCode]",
	Short: "Show match results at an event",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventCode := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		matchResults := query.MatchesByEventQuery(eventCode, year)
		matchResultsOutput := terminal.RenderMatchDetails(matchResults)
		fmt.Println(matchResultsOutput)
		return nil
	},
}

// teamMatchesCmd renders the match results for a specific team at a specific event, showing each match's teams,
// scores, and outcomes.
var teamMatchesCmd = &cobra.Command{
	Use:   "team-matches [eventCode] [teamID]",
	Short: "Show match results for a team at an event",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventCode := args[0]
		teamID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid teamID '%s', must be a number", args[1])
		}
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		matchResults := query.MatchesByEventAndTeamQuery(eventCode, teamID, year)
		matchResultsOutput := terminal.RenderMatchesByEventAndTeam(matchResults)
		fmt.Println(matchResultsOutput)
		return nil
	},
}

// renderAdvancementReport renders the advancement report for a specific event, showing which teams advanced
// and their points breakdown.
var regionAdvancementCmd = &cobra.Command{
	Use:   "region-advancement [region]",
	Short: "Show all advancing teams in a region",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		region := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		report := query.RegionAdvancementQuery(region, year)
		output := terminal.RenderRegionAdvancementReport(report)
		fmt.Println(output)
		return nil
	},
}

// eventAdvancementCmd renders region-wide advancement information for all advancing teams. It shows
// each team's advancing event, awards from that event, and other events they participated in.
var eventAdvancementCmd = &cobra.Command{
	Use:   "event-advancement [region]",
	Short: "Show qualified teams organized by qualifying events",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		region := args[0]
		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		summary := query.EventAdvancementSummaryQuery(region, year)
		output := terminal.RenderEventAdvancementSummary(summary)
		fmt.Println(output)
		return nil
	},
}

// teamRankingsCmd shows performance rankings for teams.
var teamRankingsCmd = &cobra.Command{
	Use:   "team-rankings [region]",
	Short: "Show performance rankings for teams",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		region := ""
		if len(args) > 0 {
			region = args[0]
		}
		regionFlag, _ := cmd.Flags().GetString("region")
		if regionFlag != "" {
			region = regionFlag
		}

		year, _ := cmd.Flags().GetInt("year")
		if year == 0 {
			year = defaultYear
		}
		sortBy, _ := cmd.Flags().GetString("sort")
		eventCode, _ := cmd.Flags().GetString("event")
		country, _ := cmd.Flags().GetString("country")
		limit, _ := cmd.Flags().GetInt("limit")

		performances, err := query.TeamRankingsQuery(region, country, eventCode, year)
		if err != nil {
			return err
		}

		// Convert sortBy string to SortBy type
		var sort terminal.SortBy
		switch strings.ToLower(sortBy) {
		case "opr":
			sort = terminal.SortByOPR
		case "npopr":
			sort = terminal.SortByNpOPR
		case "ccwm":
			sort = terminal.SortByCCWM
		case "dpr":
			sort = terminal.SortByDPR
		case "npdpr":
			sort = terminal.SortByNpDPR
		case "npavg":
			sort = terminal.SortByNpAVG
		case "matches":
			sort = terminal.SortByMatches
		case "team":
			sort = terminal.SortByTeamID
		default:
			sort = terminal.SortByOPR
		}

		output := terminal.RenderTeamPerformance(performances, eventCode, sort, region, year, limit)
		fmt.Println(output)
		return nil
	},
}

// init initializes the CLI commands and flags, and adds them to the root command.
func init() {
	// Add year flag to all commands that need it
	eventTeamsCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	rankingsCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	awardsCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	advancementCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	matchesCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	teamMatchesCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	regionAdvancementCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	eventAdvancementCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")
	teamRankingsCmd.Flags().IntP("year", "y", 0, "Year (defaults to FTC_SEASON environment variable)")

	// Add team-rankings specific flags
	teamRankingsCmd.Flags().StringP("sort", "s", "npavg", "Sort by: opr, npopr, ccwm, dpr, npdpr, npavg, matches, team")
	teamRankingsCmd.Flags().StringP("event", "e", "", "Event code to filter matches")
	teamRankingsCmd.Flags().StringP("region", "r", "", "Region code to filter teams")
	teamRankingsCmd.Flags().StringP("country", "c", "", "Country to filter teams")
	teamRankingsCmd.Flags().IntP("limit", "l", 0, "Limit number of teams displayed (0 = no limit)")

	// Add all commands to root
	rootCmd.AddCommand(
		teamCmd,
		teamsCmd,
		eventTeamsCmd,
		rankingsCmd,
		awardsCmd,
		advancementCmd,
		matchesCmd,
		teamMatchesCmd,
		regionAdvancementCmd,
		eventAdvancementCmd,
		teamRankingsCmd,
	)
}

// main is the entry point for the CLI application.
func main() {
	godotenv.Load()
	setLogLevelFromEnv()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
