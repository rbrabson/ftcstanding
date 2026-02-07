package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/cli"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rbrabson/ftcstanding/request"
)

func main() {
	godotenv.Load()
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

	// matchresults := query.MatchesByEventQuery("USNCSHQ2", 2025)
	// output := cli.RenderMatchDetails(matchresults)
	// fmt.Println(output)

	eventTeams := query.TeamsByEventQuery("USNCROQ", 2025)
	output := cli.RenderTeamsByEvent(eventTeams)
	fmt.Println(output)
}
