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
		if event.DateStart.After(time.Now()) {
			slog.Info("Skipping event details for future event", "event", event.EventCode, "dateStart", event.DateStart)
			continue
		}
		if event.DateEnd.Before(time.Now().Add(-48*time.Hour)) && !refresh {
			slog.Info("Skipping event details for already processed event", "event", event.EventCode, "dateStart", event.DateStart)
			continue
		}
		slog.Info("Processing event details for event", "event", event.EventCode, "dateEnd", event.DateEnd, "timeSince", time.Since(event.DateEnd))
		RequestAndSaveEventAwards(event)
		RequestAndSaveEventRankings(event)
		RequestAndSaveEventAdvancements(event)
		RequestAndSaveMatches(event)

	}
}
