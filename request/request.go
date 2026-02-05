package request

import "github.com/rbrabson/ftcstanding/database"

var (
	db database.DB
)

// Init initializes the request package with a database connection.
func Init(database database.DB) {
	db = database
}

// RequestAndSaveAll requests and saves all data for a given season.
func RequestAndSaveAll(season string) {
	RequestAndSaveAwards(season)
	RequestAndSaveTeams(season)
	events := RequestAndSaveEvents(season)
	for _, event := range events {
		RequestAndSaveEventAwards(event)
		RequestAndSaveEventRankings(event)
		RequestAndSaveEventAdvancements(event)
		RequestAndSaveMatches(event)
	}
}
