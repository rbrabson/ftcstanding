package request

// Add code to request and build the database models and save them in the database.
// This should use the ftc package to do all of the processing.

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
)

// RequestAndSaveEvents requests events from the FTC API for a given season and saves them in the database.
func RequestAndSaveEvents(season string) []*database.Event {
	events := RequestEvents(season)
	for _, event := range events {
		db.SaveEvent(event)
	}
	return events
}

// RequestEvents requests events from the FTC API for a given season.
func RequestEvents(season string) []*database.Event {
	ftcEvents, err := ftc.GetEvents(season)
	if err != nil {
		slog.Error("Error requesting events:", "year", season, "error", err)
		return nil
	}
	slog.Info("Retrieved events...", "count", len(ftcEvents))
	year, _ := strconv.Atoi(season)
	events := make([]*database.Event, 0, len(ftcEvents))
	for _, ftcEvent := range ftcEvents {
		dateStart := time.Time(ftcEvent.DateStart)
		dateEnd := time.Time(ftcEvent.DateEnd)
		eventID := database.GetEventID(ftcEvent, dateStart)

		event := database.Event{
			EventID:    eventID,
			EventCode:  ftcEvent.Code,
			Year:       year,
			Name:       ftcEvent.Name,
			Type:       ftcEvent.Type,
			RegionCode: ftcEvent.RegionCode,
			Venue:      ftcEvent.Venue,
			Address:    ftcEvent.Address,
			City:       ftcEvent.City,
			StateProv:  ftcEvent.Stateprov,
			Country:    ftcEvent.Country,
			Timezone:   ftcEvent.Timezone,
			DateStart:  dateStart,
			DateEnd:    dateEnd,
		}
		if ftcEvent.DivisionCode != nil {
			event.DivisionCode = *ftcEvent.DivisionCode
		}
		if ftcEvent.LeagueCode != nil {
			event.LeagueCode = *ftcEvent.LeagueCode
		}
		events = append(events, &event)
	}
	slog.Info("Finished processing events", "count", len(events))
	return events
}

// RequestAndSaveEventAwards requests event awards from the FTC API for a given event and saves them in the database.
func RequestAndSaveEventAwards(event *database.Event) []*database.EventAward {
	eventAwards := RequestEventAwards(event)
	for _, eventAward := range eventAwards {
		db.SaveEventAward(eventAward)
	}
	return eventAwards
}

// RequestEventAwards requests event awards from the FTC API for a given event.
func RequestEventAwards(event *database.Event) []*database.EventAward {
	ftcEventAwards, err := ftc.GetEventAwards(strconv.Itoa(event.Year), event.EventCode)
	if err != nil {
		slog.Error("Error requesting event awards:", "year", event.Year, "eventCode", event.EventCode, "error", err)
		return nil
	}
	slog.Info("Retrieved event awards...", "count", len(ftcEventAwards))
	eventAwards := make([]*database.EventAward, 0, len(ftcEventAwards))
	for _, ftcEventAward := range ftcEventAwards {
		eventAward := database.EventAward{
			EventID: event.EventID,
			AwardID: ftcEventAward.AwardID,
			TeamID:  ftcEventAward.TeamNumber,
			Name:    ftcEventAward.Name,
			Series:  ftcEventAward.Series,
		}
		eventAwards = append(eventAwards, &eventAward)
	}
	slog.Info("Finished processing event awards", "count", len(eventAwards))
	return eventAwards
}

// RequestAndSaveEventRankings requests event rankings from the FTC API for a given event and saves them in the database.
func RequestAndSaveEventRankings(event *database.Event) []*database.EventRanking {
	eventRankings := RequestEventRanking(event)
	for _, eventRanking := range eventRankings {
		db.SaveEventRanking(eventRanking)
	}
	return eventRankings
}

// RequestEventRanking requests event rankings from the FTC API for a given event.
func RequestEventRanking(event *database.Event) []*database.EventRanking {
	ftcEventRankings, err := ftc.GetRankings(strconv.Itoa(event.Year), event.EventCode)
	if err != nil {
		slog.Error("Error requesting event rankings:", "year", event.Year, "eventCode", event.EventCode, "error", err)
		return nil
	}
	eventRankings := make([]*database.EventRanking, 0, len(ftcEventRankings))
	for _, ftcEventRanking := range ftcEventRankings {
		eventRanking := database.EventRanking{
			EventID:        event.EventID,
			TeamID:         ftcEventRanking.TeamNumber,
			Rank:           ftcEventRanking.Rank,
			SortOrder1:     ftcEventRanking.SortOrder1,
			SortOrder2:     ftcEventRanking.SortOrder2,
			SortOrder3:     ftcEventRanking.SortOrder3,
			SortOrder4:     ftcEventRanking.SortOrder4,
			SortOrder5:     ftcEventRanking.SortOrder5,
			SortOrder6:     ftcEventRanking.SortOrder6,
			Wins:           ftcEventRanking.Wins,
			Losses:         ftcEventRanking.Losses,
			Ties:           ftcEventRanking.Ties,
			Dq:             ftcEventRanking.DQ,
			MatchesPlayed:  ftcEventRanking.MatchesPlayed,
			MatchesCounted: ftcEventRanking.MatchesCounted,
		}
		eventRankings = append(eventRankings, &eventRanking)
	}
	slog.Info("Finished processing event rankings", "count", len(eventRankings))
	return eventRankings
}

// RequestAndSaveEventAdvancements requests event advancements from the FTC API for a given event and saves them in the database.
func RequestAndSaveEventAdvancements(event *database.Event) []*database.EventAdvancement {
	eventAdvancements := RequestEventAdvancements(event)
	for _, eventAdvancement := range eventAdvancements {
		db.SaveEventAdvancement(eventAdvancement)
	}
	return eventAdvancements
}

// RequestEventAdvancements requests event advancements from the FTC API for a given season and event.
func RequestEventAdvancements(event *database.Event) []*database.EventAdvancement {
	ftcEventAdvancements, err := ftc.GetAdvancementsTo(strconv.Itoa(event.Year), event.EventCode)
	if err != nil {
		slog.Error("Error requesting event advancements:", "year", event.Year, "eventCode", event.EventCode, "error", err)
		return nil
	}
	eventAdvancements := make([]*database.EventAdvancement, 0, len(ftcEventAdvancements.Advancement))
	for _, ftcEventAdvancement := range ftcEventAdvancements.Advancement {
		eventAdvancement := database.EventAdvancement{
			EventID: event.EventID,
			TeamID:  ftcEventAdvancement.Team,
			Status:  strings.ToLower(ftcEventAdvancement.Status),
		}
		eventAdvancements = append(eventAdvancements, &eventAdvancement)
	}
	slog.Info("Finished processing event advancements", "count", len(eventAdvancements))
	return eventAdvancements
}

func RequestTeamsInEvent(event *database.Event) []*database.EventTeam {
	// Get all matches for the event from the database
	matches := db.GetMatchesByEvent(event.EventID)
	if len(matches) == 0 {
		slog.Warn("no matches found for event", "eventID", event.EventID)
		return nil
	}

	// Collect all unique team IDs from matches
	teamIDsMap := make(map[int]bool)
	for _, match := range matches {
		matchTeams := db.GetMatchTeams(match.MatchID)
		for _, mt := range matchTeams {
			teamIDsMap[mt.TeamID] = true
		}
	}

	eventTeams := make([]*database.EventTeam, 0, len(teamIDsMap))
	// Store EventTeam entries for all unique teams
	for teamID := range teamIDsMap {
		eventTeam := &database.EventTeam{
			EventID: event.EventID,
			TeamID:  teamID,
		}
		eventTeams = append(eventTeams, eventTeam)
	}

	slog.Info("retrieved event teams for event", "eventCode", event.EventCode, "teamCount", len(teamIDsMap))
	return eventTeams
}

// RequestAndSaveTeamsInEvent retrieves all teams for an event and saves them to the database.
func RequestAndSaveTeamsInEvent(event *database.Event) []*database.EventTeam {
	eventTeams := RequestTeamsInEvent(event)

	for _, eventTeam := range eventTeams {
		if err := db.SaveEventTeam(eventTeam); err != nil {
			slog.Error("failed to save event team", "eventID", event.EventID, "teamID", eventTeam.TeamID, "error", err)
			return nil
		}
	}

	slog.Info("stored event teams from matches", "eventID", event.EventID, "eventCode", event.EventCode, "teamCount", len(eventTeams))
	return eventTeams
}
