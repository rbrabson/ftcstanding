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
	var events []*database.Event
	if refresh {
		RequestAndSaveAwards(season)
		RequestAndSaveTeams(season)
		events = RequestAndSaveEvents(season)
	} else {
		events = db.GetAllEvents()
	}

	for i, event := range events {
		slog.Info("Processing event", "eventNumber", i+1, "totalEvents", len(events), "event", event.EventCode)
		if event.DateEnd.After(time.Now()) {
			slog.Info("Skipping event details for future event", "event", event.EventCode, "dateEnd", event.DateEnd)
			continue
		}
		filter := database.MatchFilter{
			EventIDs: []string{event.EventID},
		}
		matches := db.GetAllMatches(filter)
		if len(matches) > 0 && event.DateEnd.Before(time.Now().Add(-24*time.Hour)) && !refresh {
			slog.Info("Skipping event details for already processed event", "event", event.EventCode, "dateEnd", event.DateEnd)
			continue
		}
		slog.Info("Processing event details for event", "event", event.EventCode, "matches", len(matches), "dateEnd", event.DateEnd)
		RequestAndSaveEventAwards(event)
		RequestAndSaveEventRankings(event)
		RequestAndSaveEventAdvancements(event)
		RequestAndSaveMatches(event)
		RequestAndSaveTeamsInEvent(event)
		slog.Info("Finished processing event details for event", "event", event.EventCode)
	}
}
