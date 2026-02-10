package main

import (
	"flag"
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

// teamsInRegion lists all teams in a given region.
func teamsByRegion(region string) {
	teamsFilter := database.TeamFilter{
		HomeRegions: []string{region},
	}
	teams := query.TeamsQuery(teamsFilter)
	teamsOutput := terminal.RenderTeams(teams)
	fmt.Println(teamsOutput)
}

// teamDetails displays detailed information about a specific team.
func teamDetails(teamID int) {
	details := query.TeamDetailsQuery(teamID)
	if details == nil {
		fmt.Printf("Error: Team %d not found\n", teamID)
		return
	}
	output := terminal.RenderTeamDetails(details)
	fmt.Println(output)
}

// teamsByEvent lists all teams at a specific event.
func teamsByEvent(event string, year int) {
	eventTeams := query.TeamsByEventQuery(event, year)
	eventTeamsOutput := terminal.RenderTeamsByEvent(eventTeams)
	fmt.Println(eventTeamsOutput)
}

// teamRankingsByEvent lists the team rankings at a specific event.
func teamRankingsByEvent(event string, year int) {
	rankings := query.EventTeamRankingQuery(event, year)
	teamRankingsOutput := terminal.RenderTeamRankings(rankings)
	fmt.Println(teamRankingsOutput)
}

// awardWinnersByEvent lists the award winners at a specific event.
func awardWinnersByEvent(event string, year int) {
	awardsResults := query.AwardsByEventQuery(event, year)
	awardResultsOutput := terminal.RenderAwardsByEvent(awardsResults)
	fmt.Println(awardResultsOutput)
}

// advancementReportByEvent lists the advancement report at a specific event.
func advancementReportByEvent(event string, year int) {
	advancementReport := query.AdvancementReportQuery(event, year)
	advancementReportOutput := terminal.RenderAdvancementReport(advancementReport)
	fmt.Println(advancementReportOutput)
}

// matchResultsByEvent lists the match results at a specific event.
func matchResultsByEvent(event string, year int) {
	matchResults := query.MatchesByEventQuery(event, year)
	matchResultsOutput := terminal.RenderMatchDetails(matchResults)
	fmt.Println(matchResultsOutput)
}

// matchResultsForTeamByEvent lists the match results for a specific team at a specific event.
func matchResultsForTeamByEvent(event string, team int, year int) {
	matchResults := query.MatchesByEventAndTeamQuery(event, team, year)
	matchResultsOutput := terminal.RenderMatchesByEventAndTeam(matchResults)
	fmt.Println(matchResultsOutput)
}

// regionAdvancementReport lists all teams advancing in a region.
func regionAdvancementReport(region string, year int) {
	report := query.RegionAdvancementQuery(region, year)
	output := terminal.RenderRegionAdvancementReport(report)
	fmt.Println(output)
}

// regionalPerformanceRankings displays performance rankings for teams in a region.
func regionalPerformanceRankings(region string, eventCode string, year int, sortBy string) {
	performances, err := query.RegionalTeamRankingsQuery(region, eventCode, year)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
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

	output := terminal.RenderTeamPerformance(performances, eventCode, sort, region, year)
	fmt.Println(output)
}

// printUsage prints the usage information for the CLI.
func printUsage() {
	fmt.Println("Usage: ftc <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  advancement <eventCode> [-year]     Show advancement report for an event")
	fmt.Println("  awards <eventCode> [-year]          List award winners at an event")
	fmt.Println("  event-teams <eventCode> [-year]     List teams at an event")
	fmt.Println("  matches <eventCode> [-year]         Show match results at an event")
	fmt.Println("  rankings <eventCode> [-year]        List team rankings at an event")
	fmt.Println("  region-advancement <region> [-year] Show all advancing teams in a region")
	fmt.Println("  region-rankings <region> [-year] [-sort] [-event]  Show performance rankings for teams in a region")
	fmt.Println("  team <teamID>                       Show detailed information about a team")
	fmt.Println("  team-matches <eventCode> <teamID> [-year]  Show match results for a team at an event")
	fmt.Println("  teams <region>                      List teams in a region")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -year int    Year (defaults to FTC_SEASON environment variable)")
	fmt.Println("  -sort string Sort by: opr, npopr, ccwm, dpr, npdpr, npavg, matches, team (default: ccwm)")
	fmt.Println("  -event string Event code to filter matches (optional, for region-rankings)")
	fmt.Println()
}

// run executes the CLI command.
func run() int {
	season := os.Getenv("FTC_SEASON")
	if season == "" {
		fmt.Fprintln(os.Stderr, "Error: FTC_SEASON environment variable not set")
		return 1
	}

	defaultYear, err := strconv.Atoi(season)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid FTC_SEASON value: %s\n", season)
		return 1
	}

	db, err := database.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize database: %v\n", err)
		return 1
	}
	defer db.Close()

	request.Init(db)
	query.Init(db)
	terminal.Init(db)

	if len(os.Args) < 2 {
		printUsage()
		return 1
	}

	command := os.Args[1]

	switch command {
	case "team":
		if len(os.Args) < 3 {
			fmt.Println("Error: team command requires a teamID argument")
			fmt.Println("Usage: ftc team <teamID>")
			return 1
		}
		teamID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: invalid teamID '%s', must be a number\n", os.Args[2])
			return 1
		}
		teamDetails(teamID)

	case "teams":
		if len(os.Args) < 3 {
			fmt.Println("Error: teams command requires a region argument")
			fmt.Println("Usage: ftc teams <region>")
			return 1
		}
		region := os.Args[2]
		teamsByRegion(region)

	case "event-teams":
		fs := flag.NewFlagSet("event-teams", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: event-teams command requires an eventCode argument")
			fmt.Println("Usage: ftc event-teams <eventCode> [-year <year>]")
			return 1
		}
		eventCode := fs.Arg(0)
		teamsByEvent(eventCode, *year)

	case "rankings":
		fs := flag.NewFlagSet("rankings", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: rankings command requires an eventCode argument")
			fmt.Println("Usage: ftc rankings <eventCode> [-year <year>]")
			return 1
		}
		eventCode := fs.Arg(0)
		teamRankingsByEvent(eventCode, *year)

	case "awards":
		fs := flag.NewFlagSet("awards", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: awards command requires an eventCode argument")
			fmt.Println("Usage: ftc awards <eventCode> [-year <year>]")
			return 1
		}
		eventCode := fs.Arg(0)
		awardWinnersByEvent(eventCode, *year)

	case "advancement":
		fs := flag.NewFlagSet("advancement", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: advancement command requires an eventCode argument")
			fmt.Println("Usage: ftc advancement <eventCode> [-year <year>]")
			return 1
		}
		eventCode := fs.Arg(0)
		advancementReportByEvent(eventCode, *year)

	case "matches":
		fs := flag.NewFlagSet("matches", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: matches command requires an eventCode argument")
			fmt.Println("Usage: ftc matches <eventCode> [-year <year>]")
			return 1
		}
		eventCode := fs.Arg(0)
		matchResultsByEvent(eventCode, *year)

	case "team-matches":
		fs := flag.NewFlagSet("team-matches", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 2 {
			fmt.Println("Error: team-matches command requires eventCode and teamID arguments")
			fmt.Println("Usage: ftc team-matches <eventCode> <teamID> [-year <year>]")
			return 1
		}
		eventCode := fs.Arg(0)
		teamID, err := strconv.Atoi(fs.Arg(1))
		if err != nil {
			fmt.Printf("Error: invalid teamID '%s', must be a number\n", fs.Arg(1))
			return 1
		}
		matchResultsForTeamByEvent(eventCode, teamID, *year)

	case "region-advancement":
		fs := flag.NewFlagSet("region-advancement", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: region-advancement command requires a region argument")
			fmt.Println("Usage: ftc region-advancement <region> [-year <year>]")
			return 1
		}
		region := fs.Arg(0)
		regionAdvancementReport(region, *year)

	case "region-rankings":
		fs := flag.NewFlagSet("region-rankings", flag.ExitOnError)
		year := fs.Int("year", defaultYear, "Year")
		sortBy := fs.String("sort", "ccwm", "Sort by: opr, npopr, ccwm, dpr, npdpr, npavg, matches, team")
		eventCode := fs.String("event", "", "Event code to filter matches (optional)")
		fs.Parse(os.Args[2:])

		if fs.NArg() < 1 {
			fmt.Println("Error: region-rankings command requires a region argument")
			fmt.Println("Usage: ftc region-rankings <region> [-year <year>] [-sort <criteria>] [-event <eventCode>]")
			return 1
		}
		region := fs.Arg(0)
		if len(os.Args) > 3 {
			fs.Parse(os.Args[3:])
		}
		regionalPerformanceRankings(region, *eventCode, *year, *sortBy)

	default:
		fmt.Printf("Error: unknown command '%s'\n\n", command)
		printUsage()
		return 1
	}

	return 0
}

func main() {
	godotenv.Load()
	setLogLevelFromEnv()
	os.Exit(run())
}
