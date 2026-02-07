package query

import (
	"log/slog"
	"slices"

	"github.com/rbrabson/ftcstanding/database"
)

// MatchDetails represents a match with both alliance scores and teams per alliance.
type MatchDetails struct {
	Event        *database.Event
	Match        *database.Match
	RedAlliance  *MatchAllianceDetails
	BlueAlliance *MatchAllianceDetails
}

// MatchAllianceDetails represents the details of an alliance in a match, including its score and participating teams.
type MatchAllianceDetails struct {
	Alliance string
	Score    *database.MatchAllianceScore
	Teams    []*database.Team
}

// MatchesByEventQuery retrieves all matches for an event, including alliance scores and all participating teams.
func MatchesByEventQuery(eventCode string, year int) []*MatchDetails {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events := db.GetAllEvents(filter)
	if len(events) == 0 {
		return nil
	}
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

	// Get all matches for the event
	matches := db.GetMatchesByEvent(event.EventID)
	if matches == nil {
		return nil
	}

	var results []*MatchDetails

	// Process each match
	for _, match := range matches {
		// Get all teams in this match
		matchTeams := db.GetMatchTeams(match.MatchID)
		if matchTeams == nil {
			slog.Debug("no teams found", "matchID", match.MatchID)
			continue
		}

		// Get alliance scores
		redScore := db.GetMatchAllianceScore(match.MatchID, database.AllianceRed)
		blueScore := db.GetMatchAllianceScore(match.MatchID, database.AllianceBlue)

		// Separate teams by alliance
		var redTeams, blueTeams []*database.Team
		for _, team := range matchTeams {
			t := db.GetTeam(team.TeamID)
			if team.Alliance == database.AllianceRed {
				redTeams = append(redTeams, t)
			} else {
				blueTeams = append(blueTeams, t)
			}
		}

		results = append(results, &MatchDetails{
			Event: event,
			Match: match,
			RedAlliance: &MatchAllianceDetails{
				Alliance: database.AllianceRed,
				Score:    redScore,
				Teams:    redTeams,
			},
			BlueAlliance: &MatchAllianceDetails{
				Alliance: database.AllianceBlue,
				Score:    blueScore,
				Teams:    blueTeams,
			},
		})
	}

	slices.SortFunc(results, func(a, b *MatchDetails) int {
		if b.Match.TournamentLevel < a.Match.TournamentLevel {
			return -1
		}
		if b.Match.TournamentLevel > a.Match.TournamentLevel {
			return 1
		}
		if a.Match.MatchNumber < b.Match.MatchNumber {
			return -1
		}
		if a.Match.MatchNumber > b.Match.MatchNumber {
			return 1
		}
		return 0
	})

	return results
}
