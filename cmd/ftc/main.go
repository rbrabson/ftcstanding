package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/request"
)

func main() {
	godotenv.Load()
	FTC_SEASON := os.Getenv("FTC_SEASON")
	if FTC_SEASON == "" {
		panic("FTC_SEASON environment variable not set")
	}
	// season, _ := strconv.Atoi(FTC_SEASON)

	db, err := database.Init()
	if err != nil {
		panic(err)
	}
	request.Init(db)
	// request.RequestAndStoreTeams(season)
	team := db.GetTeam(7083)
	fmt.Println(team)
	defer db.Close()
}
