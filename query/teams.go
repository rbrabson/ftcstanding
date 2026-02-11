package query

import (
	"fmt"
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
			DateStart: event.DateStart,
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

	// Sort events by date
	sort.Slice(details.Events, func(i, j int) bool {
		return details.Events[i].DateStart.Before(details.Events[j].DateStart)
	})

	return details
}

// TeamPerformance represents performance metrics for a team across all their matches in a season.
type TeamPerformance struct {
	TeamID   int
	TeamName string
	Region   string
	OPR      float64
	NpOPR    float64
	CCWM     float64
	DPR      float64
	NpDPR    float64
	NpAVG    float64
	Matches  int
}

// RegionalTeamRankingsQuery retrieves performance metrics for all teams in a region for a given year.
// If eventCode is provided (non-empty), only rankings from that event are included.
// Performance metrics are retrieved from the team_rankings database table and combined using weighted averaging
// based on the number of matches each team played in each event.
func RegionalTeamRankingsQuery(region string, eventCode string, year int) ([]TeamPerformance, error) {
	// Get all teams in the region
	var teams []*database.Team
	if eventCode == "" {
		teams = db.GetAllTeams(database.TeamFilter{HomeRegions: []string{region}})
	} else {
		teams = db.GetAllTeams(database.TeamFilter{HomeRegions: []string{region}, EventCodes: []string{eventCode}})
	}
	if len(teams) == 0 {
		return nil, fmt.Errorf("no teams found in region %s", region)
	}

	// Get team info and build a map for easy lookup
	teamMap := make(map[int]*database.Team)
	teamIDs := make([]int, 0, len(teams))
	for _, t := range teams {
		teamMap[t.TeamID] = t
		teamIDs = append(teamIDs, t.TeamID)
	}

	// Get all events for the region, year, and optionally the event
	eventFilter := database.EventFilter{RegionCodes: []string{region}, Year: year}
	if eventCode != "" {
		eventFilter.EventCodes = []string{eventCode}
	}
	events := db.GetAllEvents(eventFilter)
	if len(events) == 0 {
		return nil, fmt.Errorf("no events found")
	}

	// Collect event IDs
	eventIDs := make([]string, 0, len(events))
	for _, event := range events {
		eventIDs = append(eventIDs, event.EventID)
	}

	// Get all team rankings for these teams and events from the database
	rankingFilter := database.TeamRankingFilter{
		TeamIDs:  teamIDs,
		EventIDs: eventIDs,
	}
	rankings := db.GetTeamRankings(rankingFilter)
	if len(rankings) == 0 {
		return nil, fmt.Errorf("no team rankings found for teams in region %s for year %d", region, year)
	}

	// Group rankings by team
	teamRankings := make(map[int][]*database.TeamRanking)
	for _, ranking := range rankings {
		teamRankings[ranking.TeamID] = append(teamRankings[ranking.TeamID], ranking)
	}

	// Combine per-event rankings using weighted averaging
	results := make([]TeamPerformance, 0, len(teamRankings))
	for teamID, eventRankings := range teamRankings {
		// Calculate weighted averages
		var totalMatches int
		var weightedOPR, weightedNpOPR, weightedCCWM float64
		var weightedDPR, weightedNpDPR, weightedNpAVG float64

		for _, ranking := range eventRankings {
			weight := float64(ranking.NumMatches)
			totalMatches += ranking.NumMatches

			weightedOPR += ranking.OPR * weight
			weightedNpOPR += ranking.NpOPR * weight
			weightedCCWM += ranking.CCWM * weight
			weightedDPR += ranking.DPR * weight
			weightedNpDPR += ranking.NpDPR * weight
			weightedNpAVG += ranking.NpAvg * weight
		}

		// Normalize by total matches
		totalWeight := float64(totalMatches)
		if totalWeight > 0 {
			weightedOPR /= totalWeight
			weightedNpOPR /= totalWeight
			weightedCCWM /= totalWeight
			weightedDPR /= totalWeight
			weightedNpDPR /= totalWeight
			weightedNpAVG /= totalWeight
		}

		team := teamMap[teamID]
		results = append(results, TeamPerformance{
			TeamID:   teamID,
			TeamName: team.Name,
			Region:   team.HomeRegion,
			OPR:      weightedOPR,
			NpOPR:    weightedNpOPR,
			CCWM:     weightedCCWM,
			DPR:      weightedDPR,
			NpDPR:    weightedNpDPR,
			NpAVG:    weightedNpAVG,
			Matches:  totalMatches,
		})
	}

	// Sort by CCWM (descending) and assign ranks
	sort.Slice(results, func(i, j int) bool {
		return results[i].OPR > results[j].CCWM
	})

	return results, nil
}
