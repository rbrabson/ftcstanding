package query

import (
	"slices"

	"github.com/rbrabson/ftcstanding/database"
)

// EventTeams represents an event with all participating teams.
type EventTeams struct {
	Event *database.Event
	Teams []*database.Team
}

// TeamsByEventQuery retrieves all teams that have or will participate in an event.
// It returns an EventTeams object containing the event and its participating teams.
func TeamsByEventQuery(eventCode string, year int) *EventTeams {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events := db.GetAllEvents(filter)
	if len(events) == 0 {
		return nil
	}

	// Find the event matching the year
	var event *database.Event
	for _, e := range events {
		if e.Year == year {
			event = e
			break
		}
	}

	if event == nil {
		return nil
	}

	// Get all event teams for the event
	eventTeams := db.GetEventTeams(event.EventID)
	if len(eventTeams) == 0 {
		return nil
	}

	// Retrieve the full team details
	var teams []*database.Team
	for _, et := range eventTeams {
		team := db.GetTeam(et.TeamID)
		if team != nil {
			teams = append(teams, team)
		}
	}

	slices.SortFunc(teams, func(a, b *database.Team) int {
		return a.TeamID - b.TeamID
	})

	return &EventTeams{
		Event: event,
		Teams: teams,
	}
}
