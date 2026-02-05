package request

import "github.com/rbrabson/ftcstanding/database"

var (
	db database.DB
)

// Init initializes the request package with a database connection.
func Init(database database.DB) {
	db = database
}
