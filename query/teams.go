package query

import (
	"sort"
	"strings"
	"time"

	"github.com/rbrabson/ftcstanding/database"
)

// Record represents a win-loss-tie record.
type Record struct {
	Wins   int
	Losses int
	Ties   int
}

// EventDetails represents detailed information about a team's participation in an event.
type EventDetails struct {
	EventCode     string
	EventName     string
	DateStart     time.Time
	QualRank      int
	TotalRecord   Record
	QualRecord    Record
	PlayoffRecord Record
	Advanced      bool
	Awards        []string
}

// TeamDetails represents comprehensive information about a team.
type TeamDetails struct {
	TeamID        int
	Name          string
	FullName      string
	City          string
	StateProv     string
	Country       string
	Region        string
	RookieYear    int
	TotalRecord   Record
	QualRecord    Record
	PlayoffRecord Record
	Events        []EventDetails
}

// TeamsQuery returns a list of teams that match the given filter.
func TeamsQuery(filter ...database.TeamFilter) ([]*database.Team, error) {
	teams, err := db.GetAllTeams(filter...)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// TeamDetailsQuery returns detailed information about a specific team.
func TeamDetailsQuery(teamID int) (*TeamDetails, error) {
	// Get team basic information
	team, err := db.GetTeam(teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	// Initialize team details
	details := &TeamDetails{
		TeamID:     team.TeamID,
		Name:       team.Name,
		FullName:   team.FullName,
		City:       team.City,
		StateProv:  team.StateProv,
		Country:    team.Country,
		Region:     team.HomeRegion,
		RookieYear: team.RookieYear,
		Events:     []EventDetails{},
	}

	// Get all events for this team
	eventIDs, err := db.GetEventsByTeam(teamID)
	if err != nil {
		return nil, err
	}

	// Process each event
	for _, eventID := range eventIDs {
		event, err := db.GetEvent(eventID)
		if err != nil {
			return nil, err
		}
		if event == nil {
			continue
		}

		eventDetail := EventDetails{
			EventCode: event.EventCode,
			EventName: event.Name,
			DateStart: event.DateStart,
		}

		// Get qualification ranking for this team at this event
		rankings, err := db.GetEventRankings(eventID)
		if err != nil {
			return nil, err
		}
		for _, ranking := range rankings {
			if ranking.TeamID == teamID {
				eventDetail.QualRank = ranking.Rank
				break
			}
		}

		// Get matches for this event
		matches, err := db.GetMatchesByEvent(eventID)
		if err != nil {
			return nil, err
		}

		// Calculate records by going through each match
		for _, match := range matches {
			matchTeams, err := db.GetMatchTeams(match.MatchID)
			if err != nil {
				return nil, err
			}

			// Check if this team participated in the match
			var teamAlliance string
			found := false
			for _, mt := range matchTeams {
				if mt.TeamID == teamID && mt.OnField && !mt.Dq {
					teamAlliance = mt.Alliance
					found = true
					break
				}
			}

			if !found {
				continue
			}

			// Get alliance scores
			teamScore, err := db.GetMatchAllianceScore(match.MatchID, teamAlliance)
			if err != nil {
				return nil, err
			}
			opponentAlliance := database.AllianceRed
			if teamAlliance == database.AllianceRed {
				opponentAlliance = database.AllianceBlue
			}
			opponentScore, err := db.GetMatchAllianceScore(match.MatchID, opponentAlliance)
			if err != nil {
				return nil, err
			}

			if teamScore == nil || opponentScore == nil {
				continue
			}

			// Update records based on tournament level
			isPlayoff := strings.EqualFold(match.TournamentLevel, "playoff")

			// Determine if this team won, lost, or tied and update records
			switch {
			case teamScore.TotalPoints > opponentScore.TotalPoints:
				eventDetail.TotalRecord.Wins++
				details.TotalRecord.Wins++
				if isPlayoff {
					eventDetail.PlayoffRecord.Wins++
					details.PlayoffRecord.Wins++
				} else {
					eventDetail.QualRecord.Wins++
					details.QualRecord.Wins++
				}
			case teamScore.TotalPoints < opponentScore.TotalPoints:
				eventDetail.TotalRecord.Losses++
				details.TotalRecord.Losses++
				if isPlayoff {
					eventDetail.PlayoffRecord.Losses++
					details.PlayoffRecord.Losses++
				} else {
					eventDetail.QualRecord.Losses++
					details.QualRecord.Losses++
				}
			default:
				eventDetail.TotalRecord.Ties++
				details.TotalRecord.Ties++
				if isPlayoff {
					eventDetail.PlayoffRecord.Ties++
					details.PlayoffRecord.Ties++
				} else {
					eventDetail.QualRecord.Ties++
					details.QualRecord.Ties++
				}
			}
		}

		// Check if team advanced from this event
		advancements, err := db.GetEventAdvancements(eventID)
		if err != nil {
			return nil, err
		}
		for _, adv := range advancements {
			if adv.TeamID == teamID {
				eventDetail.Advanced = true
				break
			}
		}

		// Get awards won at this event
		awards, err := db.GetTeamAwardsByEvent(eventID, teamID)
		if err != nil {
			return nil, err
		}
		eventDetail.Awards = make([]string, 0, len(awards))
		for _, award := range awards {
			eventDetail.Awards = append(eventDetail.Awards, award.Name)
		}

		details.Events = append(details.Events, eventDetail)
	}

	// Sort events by date
	sort.Slice(details.Events, func(i, j int) bool {
		return details.Events[i].DateStart.Before(details.Events[j].DateStart)
	})

	return details, nil
}
