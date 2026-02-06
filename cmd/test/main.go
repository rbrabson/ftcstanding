package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/cli"
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
	defer db.Close()

	request.Init(db)

	output := cli.RenderTable(db)
	fmt.Println(output)
}
