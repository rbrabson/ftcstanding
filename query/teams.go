package query

import (
	"strings"

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
func TeamsQuery(filter ...database.TeamFilter) []*database.Team {
	return db.GetAllTeams(filter...)
}

// TeamDetailsQuery returns detailed information about a specific team.
func TeamDetailsQuery(teamID int) *TeamDetails {
	// Get team basic information
	team := db.GetTeam(teamID)
	if team == nil {
		return nil
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
	eventIDs := db.GetEventsByTeam(teamID)

	// Process each event
	for _, eventID := range eventIDs {
		event := db.GetEvent(eventID)
		if event == nil {
			continue
		}

		eventDetail := EventDetails{
			EventCode: event.EventCode,
			EventName: event.Name,
		}

		// Get qualification ranking for this team at this event
		rankings := db.GetEventRankings(eventID)
		for _, ranking := range rankings {
			if ranking.TeamID == teamID {
				eventDetail.QualRank = ranking.Rank
				break
			}
		}

		// Get matches for this event
		matches := db.GetMatchesByEvent(eventID)

		// Calculate records by going through each match
		for _, match := range matches {
			matchTeams := db.GetMatchTeams(match.MatchID)

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
			teamScore := db.GetMatchAllianceScore(match.MatchID, teamAlliance)
			opponentAlliance := database.AllianceRed
			if teamAlliance == database.AllianceRed {
				opponentAlliance = database.AllianceBlue
			}
			opponentScore := db.GetMatchAllianceScore(match.MatchID, opponentAlliance)

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
		advancements := db.GetEventAdvancements(eventID)
		for _, adv := range advancements {
			if adv.TeamID == teamID {
				eventDetail.Advanced = true
				break
			}
		}

		// Get awards won at this event
		awards := db.GetTeamAwardsByEvent(eventID, teamID)
		eventDetail.Awards = make([]string, 0, len(awards))
		for _, award := range awards {
			eventDetail.Awards = append(eventDetail.Awards, award.Name)
		}

		details.Events = append(details.Events, eventDetail)
	}

	return details
}
