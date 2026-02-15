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

	awards, err := db.GetAllAwards()
	if err != nil {
		slog.Warn("failed to load awards from db", "error", err)
	}
	if refresh || len(awards) == 0 {
		awards = RequestAndSaveAwards(season)
	}
	teams, err := db.GetAllTeams()
	if err != nil {
		slog.Warn("failed to load teams from db", "error", err)
	}
	if refresh || len(teams) == 0 {
		teams = RequestAndSaveTeams(season)
	}

	events, err := db.GetAllEvents()
	if err != nil {
		slog.Warn("failed to load events from db", "error", err)
	}
	if refresh || len(events) == 0 {
		events = RequestAndSaveEvents(season)
	}

	for i, event := range events {
		slog.Info("Processing event", "eventNumber", i+1, "totalEvents", len(events), "event", event.EventCode)
		if event.DateEnd.After(time.Now()) {
			slog.Info("Skipping event details for future event", "event", event.EventCode, "dateEnd", event.DateEnd)
			continue
		}
		advancementFilter := database.AdvancementFilter{
			EventCodes: []string{event.EventCode},
		}
		advancements, err := db.GetAllAdvancements(advancementFilter)
		if err != nil {
			slog.Warn("failed to load advancements", "event", event.EventCode, "error", err)
		}
		if !refresh && len(advancements) > 0 && event.DateEnd.Before(time.Now().Add(-48*time.Hour)) {
			slog.Info("Skipping event details for already processed event", "event", event.EventCode, "advancements", len(advancements), "dateEnd", event.DateEnd)
			continue
		}
		filter := database.MatchFilter{
			EventIDs: []string{event.EventID},
		}
		matches, err := db.GetAllMatches(filter)
		if err != nil {
			slog.Warn("failed to load matches", "event", event.EventCode, "error", err)
		}
		if !refresh && len(matches) > 0 && event.DateEnd.Before(time.Now().Add(-24*6*time.Hour)) {
			slog.Info("Skipping event details for already processed event with advancements", "event", event.EventCode, "matches", len(matches), "dateEnd", event.DateEnd)
			continue
		}
		slog.Info("Processing event details for event", "event", event.EventCode, "matches", len(matches), "advancements", len(advancements), "dateEnd", event.DateEnd)
		RequestAndSaveEventAwards(event)
		RequestAndSaveEventRankings(event)
		RequestAndSaveEventAdvancements(event)
		RequestAndSaveMatches(event)
		RequestAndSaveTeamsInEvent(event)
		RequestAndSaveTeamRankings(event)
		slog.Info("Finished processing event details for event", "event", event.EventCode)
	}
}
