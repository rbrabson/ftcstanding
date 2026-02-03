package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DATA_SOURCE_NAME")
	if dsn == "" {
		panic("DATA_SOURCE_NAME environment variable not set")
	}
	if err := database.InitDB(dsn); err != nil {
		panic(err)
	}
	defer database.CloseStatements()

	// Initialize prepared statements
	if err := database.InitStatements(); err != nil {
		panic(err)
	}

}
