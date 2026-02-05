package main

import (
	"github.com/rbrabson/ftcstanding/database"
)

func main() {
	db, err := database.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
