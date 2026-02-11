package query

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/lambda"
	"github.com/rbrabson/ftcstanding/performance"
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

// RegionalTeamRankingsQuery calculates performance metrics for all teams in a region for a given year.
// If eventCode is provided (non-empty), only matches from that event are included.
// Performance metrics are calculated on a per-event basis and then combined using weighted averaging
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
	for _, t := range teams {
		teamMap[t.TeamID] = t
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

	// Structure to hold per-event performance data for each team
	type eventPerformance struct {
		OPR     float64
		NpOPR   float64
		CCWM    float64
		DPR     float64
		NpDPR   float64
		NpAVG   float64
		Matches int
	}

	// Map to accumulate weighted metrics: teamID -> []eventPerformance
	teamEventData := make(map[int][]eventPerformance)

	// Process each event separately
	for _, event := range events {
		dbMatches := db.GetMatchesByEvent(event.EventID)

		var matches []performance.Match
		teamSet := make(map[int]any)

		for _, dbMatch := range dbMatches {
			// Get alliance scores
			redScore := db.GetMatchAllianceScore(dbMatch.MatchID, database.AllianceRed)
			blueScore := db.GetMatchAllianceScore(dbMatch.MatchID, database.AllianceBlue)

			if redScore == nil || blueScore == nil {
				continue
			}

			// Get teams in the match
			matchTeams := db.GetMatchTeams(dbMatch.MatchID)

			var redTeams []int
			var blueTeams []int

			for _, mt := range matchTeams {
				if !mt.OnField || mt.Dq {
					continue
				}

				if mt.Alliance == database.AllianceRed {
					redTeams = append(redTeams, mt.TeamID)
				} else {
					blueTeams = append(blueTeams, mt.TeamID)
				}

				teamSet[mt.TeamID] = struct{}{}
			}

			// Only include matches with teams on both alliances
			if len(redTeams) == 0 || len(blueTeams) == 0 {
				continue
			}

			matches = append(matches, performance.Match{
				RedTeams:      redTeams,
				BlueTeams:     blueTeams,
				RedScore:      float64(redScore.TotalPoints),
				BlueScore:     float64(blueScore.TotalPoints),
				RedPenalties:  float64(redScore.FoulPointsCommitted),
				BluePenalties: float64(blueScore.FoulPointsCommitted),
			})
		}

		// Skip events with no matches
		if len(matches) == 0 {
			continue
		}

		// Convert teamSet to sorted slice
		eventTeams := slices.Collect(maps.Keys(teamSet))
		sort.Ints(eventTeams)

		// Calculate lambda for this event
		lambdaValue := lambda.GetLambda(matches)

		slog.Info("processing event", "event", event.EventCode, "region", region, "season", year, "matches", len(matches), "teams", len(eventTeams), "lambda", lambdaValue)

		// Calculate performance metrics for this event
		calculator := performance.Calculator{
			Matches: matches,
			Teams:   eventTeams,
			Lambda:  lambdaValue,
		}

		opr := calculator.CalculateOPR()
		npopr := calculator.CalculateNpOPR()
		ccwm := calculator.CalculateCCWM()
		dpr := calculator.CalculateDPR()
		npdpr := calculator.CalculateNpDPR()

		// Store per-event results for each team
		for _, teamID := range eventTeams {
			// Count matches for this team in this event
			matchCount := 0
			for _, m := range matches {
				if slices.Contains(m.RedTeams, teamID) || slices.Contains(m.BlueTeams, teamID) {
					matchCount++
				}
			}

			npavg := calculator.CalculateNpAVG(matches, teamID)

			teamEventData[teamID] = append(teamEventData[teamID], eventPerformance{
				OPR:     opr[teamID],
				NpOPR:   npopr[teamID],
				CCWM:    ccwm[teamID],
				DPR:     dpr[teamID],
				NpDPR:   npdpr[teamID],
				NpAVG:   npavg,
				Matches: matchCount,
			})
		}
	}

	if len(teamEventData) == 0 {
		return nil, fmt.Errorf("no matches found for teams in region %s for year %d", region, year)
	}

	// Combine per-event results using weighted averaging
	results := make([]TeamPerformance, 0, len(teamEventData))
	for teamID, eventPerfs := range teamEventData {
		// Calculate weighted averages
		var totalMatches int
		var weightedOPR, weightedNpOPR, weightedCCWM float64
		var weightedDPR, weightedNpDPR, weightedNpAVG float64

		for _, ep := range eventPerfs {
			weight := float64(ep.Matches)
			totalMatches += ep.Matches

			weightedOPR += ep.OPR * weight
			weightedNpOPR += ep.NpOPR * weight
			weightedCCWM += ep.CCWM * weight
			weightedDPR += ep.DPR * weight
			weightedNpDPR += ep.NpDPR * weight
			weightedNpAVG += ep.NpAVG * weight
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

	// Sort by OPR (descending) and assign ranks
	sort.Slice(results, func(i, j int) bool {
		return results[i].OPR > results[j].OPR
	})

	return results, nil
}
