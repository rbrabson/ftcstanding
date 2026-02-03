package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/dbmodel"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DATA_SOURCE_NAME")
	if dsn == "" {
		panic("DATA_SOURCE_NAME environment variable not set")
	}
	if err := dbmodel.InitDB(dsn); err != nil {
		panic(err)
	}
	defer dbmodel.CloseStatements()

	// Initialize prepared statements
	if err := dbmodel.InitStatements(); err != nil {
		panic(err)
	}

}
