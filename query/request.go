package query

import "github.com/rbrabson/ftcstanding/database"

var (
	db database.DB
)

func Init(database database.DB) {
	db = database
}
