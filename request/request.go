package request

import (
	"log/slog"
	"time"

	"github.com/rbrabson/ftcstanding/database"
)

var (
	db database.DB
)

// Init initializes the request package with a database connection.
func Init(database database.DB) {
	db = database
}

// RequestAndSaveAll requests and saves all data for a given season.
func RequestAndSaveAll(season string, refresh bool) {
	RequestAndSaveAwards(season)
	RequestAndSaveTeams(season)
	events := RequestAndSaveEvents(season)
	for _, event := range events {
		// TODO: keep track of the last time we requested event details, and don't re-process details unless requested to refresh all data.
		//       alternatively, see if we've previously processed it by looking one of the other DB tables and skipping if the details are
		//       there. This should dramatically speed things up.
		if event.DateStart.Before(time.Now()) {
			if !refresh {
				// Skip events that finished more than 2 days ago
				if event.DateEnd.Before(time.Now().Add(-48 * time.Hour)) {
					slog.Info("Skipping event details for already processed event", "event", event.EventCode, "dateStart", event.DateStart)
					continue
				} else {
					slog.Info("Processing event details for recent event", "event", event.EventCode, "dateStart", event.DateStart)
				}
			}
			RequestAndSaveEventAwards(event)
			RequestAndSaveEventRankings(event)
			RequestAndSaveEventAdvancements(event)
			RequestAndSaveMatches(event)
		} else {
			slog.Info("Skipping event details for future event", "event", event.EventCode, "dateStart", event.DateStart)
		}
	}
}
