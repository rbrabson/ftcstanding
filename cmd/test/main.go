package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/cli"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rbrabson/ftcstanding/request"
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
	cli.Init(db)

	// filter := database.TeamFilter{
	// 	HomeRegions: []string{"USNC"},
	// }
	// teams := query.TeamsQuery(filter)
	// output := cli.RenderTeams(teams)
	// fmt.Println(output)

	// matchResults := query.TeamMatchesByEventQuery(23532, "USNCROQ", 2025)
	// output := cli.RenderTeamMatchDetails(matchResults)
	// fmt.Println(output)

	// matchresults := query.MatchesByEventQuery("USNCSHQ", 2025)
	// output := cli.RenderMatchDetails(matchresults)
	// fmt.Println(output)

	// eventTeams := query.TeamsByEventQuery("USNCROQ", 2025)
	// output := cli.RenderTeamsByEvent(eventTeams)
	// fmt.Println(output)

	// rankings := query.EventTeamRankingQuery("USNCSHQ", 2025)
	// output := cli.RenderTeamRankings(rankings)
	// fmt.Println(output)

	// awardsResults := query.AwardsByEventQuery("USNCSHQ", 2025)
	// output := cli.RenderAwardsByEvent(awardsResults)
	// fmt.Println(output)

	advancementReport := query.AdvancementReportQuery("USNCSHQ", 2025)
	output := cli.RenderAdvancementReport(advancementReport)
	fmt.Println(output)
}
