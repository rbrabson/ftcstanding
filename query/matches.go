package query

import (
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

// TeamMatchResult represents a match from a specific team's perspective with match outcome.
type TeamMatchResult struct {
	Event            *database.Event
	Match            *database.Match
	Team             *database.Team
	TeamAlliance     *MatchAllianceDetails
	OpponentAlliance *MatchAllianceDetails
	Result           string // "Won", "Lost", or "Tied"
}

// MatchesByEventQuery retrieves all matches for an event, including alliance scores and all participating teams.
func MatchesByEventQuery(eventCode string, year int) ([]*MatchDetails, error) {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events, err := db.GetAllEvents(filter)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	var event *database.Event
	for _, e := range events {
		if e.Year == year {
			event = e
			break
		}
	}

	if event == nil {
		return nil, nil
	}

	// Get all matches for the event
	matches, err := db.GetMatchesByEvent(event.EventID)
	if err != nil {
		return nil, err
	}
	if matches == nil {
		return nil, nil
	}

	var results []*MatchDetails

	// Process each match
	for _, match := range matches {
		// Get all teams in this match
		matchTeams, err := db.GetMatchTeams(match.MatchID)
		if err != nil {
			return nil, err
		}
		if matchTeams == nil {
			continue
		}

		// Get alliance scores
		redScore, err := db.GetMatchAllianceScore(match.MatchID, database.AllianceRed)
		if err != nil {
			return nil, err
		}
		blueScore, err := db.GetMatchAllianceScore(match.MatchID, database.AllianceBlue)
		if err != nil {
			return nil, err
		}

		// Separate teams by alliance
		var redTeams, blueTeams []*database.Team
		for _, team := range matchTeams {
			t, err := db.GetTeam(team.TeamID)
			if err != nil {
				return nil, err
			}
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

	return results, nil
}

// MatchesByEventAndTeamQuery retrieves all matches for a specific team at an event.
// It shows the match from the team's perspective with their result (Won/Lost/Tied).
func MatchesByEventAndTeamQuery(eventCode string, teamID int, year int) ([]*TeamMatchResult, error) {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
		Year:       year,
	}
	events, err := db.GetAllEvents(filter)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	event := events[0]

	matches, err := db.GetMatchesByEvent(event.EventID)
	if err != nil {
		return nil, err
	}
	if matches == nil {
		return nil, nil
	}

	// Get the team object
	team, err := db.GetTeam(teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	var results []*TeamMatchResult
	for _, match := range matches {
		// Get all teams in this match
		matchTeams, err := db.GetMatchTeams(match.MatchID)
		if err != nil {
			return nil, err
		}
		if len(matchTeams) == 0 {
			continue
		}

		// Check if the specified team is in this match
		var teamAlliance string
		teamFound := false
		for _, mt := range matchTeams {
			if mt.TeamID == teamID {
				teamAlliance = mt.Alliance
				teamFound = true
				break
			}
		}

		// Skip this match if the team didn't participate
		if !teamFound {
			continue
		}

		// Get alliance scores
		redScore, err := db.GetMatchAllianceScore(match.MatchID, database.AllianceRed)
		if err != nil {
			return nil, err
		}
		blueScore, err := db.GetMatchAllianceScore(match.MatchID, database.AllianceBlue)
		if err != nil {
			return nil, err
		}

		// Separate teams by alliance
		var redTeams, blueTeams []*database.Team
		for _, mt := range matchTeams {
			t, err := db.GetTeam(mt.TeamID)
			if err != nil {
				return nil, err
			}
			if mt.Alliance == database.AllianceRed {
				redTeams = append(redTeams, t)
			} else {
				blueTeams = append(blueTeams, t)
			}
		}

		// Determine team's alliance and opponent alliance
		var teamAllianceDetails, opponentAllianceDetails *MatchAllianceDetails
		var teamScore, opponentScore *database.MatchAllianceScore

		if teamAlliance == database.AllianceRed {
			teamAllianceDetails = &MatchAllianceDetails{
				Alliance: database.AllianceRed,
				Score:    redScore,
				Teams:    redTeams,
			}
			opponentAllianceDetails = &MatchAllianceDetails{
				Alliance: database.AllianceBlue,
				Score:    blueScore,
				Teams:    blueTeams,
			}
			teamScore = redScore
			opponentScore = blueScore
		} else {
			teamAllianceDetails = &MatchAllianceDetails{
				Alliance: database.AllianceBlue,
				Score:    blueScore,
				Teams:    blueTeams,
			}
			opponentAllianceDetails = &MatchAllianceDetails{
				Alliance: database.AllianceRed,
				Score:    redScore,
				Teams:    redTeams,
			}
			teamScore = blueScore
			opponentScore = redScore
		}

		// Determine the result
		result := "Tied"
		if teamScore != nil && opponentScore != nil {
			if teamScore.TotalPoints > opponentScore.TotalPoints {
				result = "Won"
			} else if teamScore.TotalPoints < opponentScore.TotalPoints {
				result = "Lost"
			}
		}

		results = append(results, &TeamMatchResult{
			Event:            event,
			Match:            match,
			Team:             team,
			TeamAlliance:     teamAllianceDetails,
			OpponentAlliance: opponentAllianceDetails,
			Result:           result,
		})
	}

	// Sort by tournament level and match number
	slices.SortFunc(results, func(a, b *TeamMatchResult) int {
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

	return results, nil
}

// GetEventTeamsQuery retrieves all EventTeam entries for a given event.
func GetEventTeamsQuery(eventCode string, year int) ([]*database.EventTeam, error) {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events, err := db.GetAllEvents(filter)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}
	var event *database.Event
	for _, e := range events {
		if e.Year == year {
			event = e
			break
		}
	}

	if event == nil {
		return nil, nil
	}

	teams, err := db.GetEventTeams(event.EventID)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// SaveEventTeam saves an EventTeam entry to the database.
func SaveEventTeam(eventID string, teamID int) error {
	eventTeam := &database.EventTeam{
		EventID: eventID,
		TeamID:  teamID,
	}
	return db.SaveEventTeam(eventTeam)
}
