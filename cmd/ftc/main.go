package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
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
	request.Init(db)
	request.RequestAndSaveAwards(season)
	request.RequestAndSaveTeams(season)
	events := request.RequestAndSaveEvents(season)
	for _, event := range events {
		request.RequestAndSaveEventAwards(event)
		request.RequestAndSaveEventRankings(event)
		request.RequestAndSaveEventAdvancements(event)

		request.RequestAndSaveMatches(event)

	}
	defer db.Close()
}
