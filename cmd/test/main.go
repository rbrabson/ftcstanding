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

	// // Team Listing in NC
	// teamsFilter := database.TeamFilter{
	// 	HomeRegions: []string{"USNC"},
	// }
	// teams := query.TeamsQuery(teamsFilter)
	// teamsOutput := terminal.RenderTeams(teams)
	// fmt.Println(teamsOutput)

	// // Teams at a specific event
	// eventTeams := query.TeamsByEventQuery("USNCSHQ2", 2025)
	// eventTeamsOutput := terminal.RenderTeamsByEvent(eventTeams)
	// fmt.Println(eventTeamsOutput)

	// // Team rankings at a specific event
	// rankings := query.EventTeamRankingQuery("USNCSHQ2", 2025)
	// teamRankingsOutput := terminal.RenderTeamRankings(rankings)
	// fmt.Println(teamRankingsOutput)

	// // Award winners at a specific event
	// awardsResults := query.AwardsByEventQuery("USNCSHQ", 2025)
	// awardResultsOutput := terminal.RenderAwardsByEvent(awardsResults)
	// fmt.Println(awardResultsOutput)

	// // Advancement report for a specific event
	// advancementReport := query.AdvancementReportQuery("USNCSHQ2", 2025)
	// advancementReportOutput := terminal.RenderAdvancementReport(advancementReport)
	// fmt.Println(advancementReportOutput)

	// Match results for a specific event
	matchresults := query.MatchesByEventQuery("USNCSHQ2", 2025)
	matchResultsOutput := terminal.RenderMatchDetails(matchresults)
	fmt.Println(matchResultsOutput)

	// // Match results for a specific team at a specific event
	// matchTeamResults := query.MatchesByEventAndTeamQuery("USNCSHQ2", 24260, 2025)
	// matchTeamResultsOutput := terminal.RenderMatchesByEventAndTeam(matchTeamResults)
	// fmt.Println(matchTeamResultsOutput)
}
