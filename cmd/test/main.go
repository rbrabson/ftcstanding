package main

import (
	"fmt"
	"log/slog"
	"os"
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

func main() {
	godotenv.Load()

	setLogLevelFromEnv()

	season := os.Getenv("FTC_SEASON")
	if season == "" {
		panic("FTC_SEASON environment variable not set")
	}

	db, err := database.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	request.Init(db)
	query.Init(db)
	terminal.Init(db)

	// teamsByRegion("USNC")
	// teamsByEvent("USNCSHQ2", 2025)
	teamRankingsByEvent("USNCSHQ2", 2025)
	// awardWinnersByEvent("USNCSHQ2", 2025)
	// advancementReportByEvent("USNCSHQ2", 2025)
	// matchResultsByEvent("USNCSHQ2", 2025)
	// matchResultsForTeamByEvent("USNCSHQ2", 7083, 2025)
}
