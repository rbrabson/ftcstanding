package query

import (
	"log/slog"
	"slices"

	"github.com/rbrabson/ftcstanding/database"
)

// TeamMatchDetails represents a match with alliance scores and all participating teams.
type TeamMatchDetails struct {
	Event         *database.Event
	Match         *database.Match
	AllianceScore *database.MatchAllianceScore
	Teams         []*database.MatchTeam
}

// TeamsQuery returns a list of teams that match the given filter.
func TeamsQuery(filter ...database.TeamFilter) []*database.Team {
	return db.GetAllTeams(filter...)
}

// TeamMatchesByEventQuery retrieves all matches for a team at an event, including alliance scores and all participating teams.
func TeamMatchesByEventQuery(teamID int, eventCode string, year int) []*TeamMatchDetails {
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

	var results []*TeamMatchDetails

	// Process each match
	for _, match := range matches {
		// Get all teams in this match
		matchTeams := db.GetMatchTeams(match.MatchID)
		if matchTeams == nil {
			slog.Debug("no teams found", "matchID", match.MatchID)
			continue
		}

		// Check if the team is in this match and find their alliance
		var teamAlliance string
		teamFound := false
		for _, mt := range matchTeams {
			if mt.TeamID == teamID {
				teamAlliance = mt.Alliance
				teamFound = true
				break
			}
		}

		if !teamFound {
			slog.Debug("team not found in match", "matchID", match.MatchID, "teamID", teamID)
			continue
		}

		// Get the alliance score for this team's alliance
		allianceScore := db.GetMatchAllianceScore(match.MatchID, teamAlliance)

		results = append(results, &TeamMatchDetails{
			Event:         event,
			Match:         match,
			AllianceScore: allianceScore,
			Teams:         matchTeams,
		})
	}

	slices.SortFunc(results, func(a, b *TeamMatchDetails) int {
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
