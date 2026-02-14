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
func TeamsByEventQuery(eventCode string, year int) (*EventTeams, error) {
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

	// Find the event matching the year
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

	// Get all event teams for the event
	eventTeams, err := db.GetEventTeams(event.EventID)
	if err != nil {
		return nil, err
	}
	if len(eventTeams) == 0 {
		return nil, nil
	}

	// Retrieve the full team details
	var teams []*database.Team
	for _, et := range eventTeams {
		team, err := db.GetTeam(et.TeamID)
		if err != nil {
			return nil, err
		}
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
	}, nil
}

// TeamRanking represents a team with its ranking information.
type TeamRanking struct {
	Team           *database.Team
	Ranking        *database.EventRanking
	HighMatchScore int // Highest total points scored in any match
}

// EventTeamRankings represents an event with all team rankings.
type EventTeamRankings struct {
	Event        *database.Event
	TeamRankings []*TeamRanking
}

// EventTeamRankingQuery retrieves an event and all teams with their rankings, sorted by rank.
func EventTeamRankingQuery(eventCode string, year int) (*EventTeamRankings, error) {
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

	// Find the event matching the year
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

	// Get all event rankings for the event
	eventRankings, err := db.GetEventRankings(event.EventID)
	if err != nil {
		return nil, err
	}
	if len(eventRankings) == 0 {
		return nil, nil
	}

	// Get all matches for the event to calculate high scores
	matches, err := db.GetMatchesByEvent(event.EventID)
	if err != nil {
		return nil, err
	}

	// Calculate high score for each team
	teamHighScores := make(map[int]int)
	for _, match := range matches {
		matchTeams, err := db.GetMatchTeams(match.MatchID)
		if err != nil {
			return nil, err
		}
		for _, mt := range matchTeams {
			// Get the alliance score for this team's alliance
			allianceScore, err := db.GetMatchAllianceScore(match.MatchID, mt.Alliance)
			if err != nil {
				return nil, err
			}
			var opposingAllianceScore *database.MatchAllianceScore
			if mt.Alliance == "red" {
				opposingAllianceScore, err = db.GetMatchAllianceScore(match.MatchID, "blue")
			} else {
				opposingAllianceScore, err = db.GetMatchAllianceScore(match.MatchID, "red")
			}
			if err != nil {
				return nil, err
			}
			if allianceScore != nil {
				var totalPoints int
				if opposingAllianceScore != nil {
					totalPoints = allianceScore.TotalPoints - opposingAllianceScore.FoulPointsCommitted
				} else {
					totalPoints = allianceScore.TotalPoints
				}
				if totalPoints > teamHighScores[mt.TeamID] {
					teamHighScores[mt.TeamID] = totalPoints
				}
			}
		}
	}

	// Retrieve the full team details and combine with rankings
	var teamRankings []*TeamRanking
	for _, ranking := range eventRankings {
		team, err := db.GetTeam(ranking.TeamID)
		if err != nil {
			return nil, err
		}
		if team != nil {
			teamRankings = append(teamRankings, &TeamRanking{
				Team:           team,
				Ranking:        ranking,
				HighMatchScore: teamHighScores[ranking.TeamID],
			})
		}
	}

	// Sort by rank
	slices.SortFunc(teamRankings, func(a, b *TeamRanking) int {
		return a.Ranking.Rank - b.Ranking.Rank
	})

	return &EventTeamRankings{
		Event:        event,
		TeamRankings: teamRankings,
	}, nil
}
