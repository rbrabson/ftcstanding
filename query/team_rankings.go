package query

import (
	"fmt"
	"sort"

	"github.com/rbrabson/ftcstanding/database"
)

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

// TeamRankingsQuery retrieves performance metrics for all teams in a region for a given year.
// If region is provided (non-empty), only teams from that region are included; otherwise all teams are included.
// If country is provided (non-empty), only teams from that country are included.
// If eventCode is provided (non-empty), only rankings from that event are included.
// Performance metrics are retrieved from the team_rankings database table and combined using weighted averaging
// based on the number of matches each team played in each event.
func TeamRankingsQuery(region string, country string, eventCode string, year int) ([]TeamPerformance, error) {
	// Build team filter
	var teamFilter database.TeamFilter
	if region != "" {
		teamFilter.HomeRegions = []string{region}
	}
	if country != "" {
		teamFilter.Countries = []string{country}
	}
	if eventCode != "" {
		teamFilter.EventCodes = []string{eventCode}
	}

	// Get all teams based on filters
	var teams []*database.Team
	if region == "" && country == "" && eventCode == "" {
		teams = db.GetAllTeams()
	} else {
		teams = db.GetAllTeams(teamFilter)
	}
	if len(teams) == 0 {
		if region != "" {
			return nil, fmt.Errorf("no teams found in region %s", region)
		}
		if country != "" {
			return nil, fmt.Errorf("no teams found in country %s", country)
		}
		return nil, fmt.Errorf("no teams found")
	}

	// Get team info and build a map for easy lookup
	teamMap := make(map[int]*database.Team)
	teamIDs := make([]int, 0, len(teams))
	for _, t := range teams {
		teamMap[t.TeamID] = t
		teamIDs = append(teamIDs, t.TeamID)
	}

	// Build event filter
	eventFilter := database.EventFilter{Year: year}
	if region != "" {
		eventFilter.RegionCodes = []string{region}
	}
	if eventCode != "" {
		eventFilter.EventCodes = []string{eventCode}
	} else {
		// When no specific event is specified, only include qualifiers and championships
		// (exclude scrimmages, league meets, and other non-competitive events)
		eventFilter.Types = []string{"2", "4"}
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
		if region != "" {
			return nil, fmt.Errorf("no team rankings found for teams in region %s for year %d", region, year)
		}
		return nil, fmt.Errorf("no team rankings found for year %d", year)
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

	// Sort by OPR (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].NpAVG > results[j].NpAVG
	})

	return results, nil
}

// TeamEventPerformance represents performance metrics for a team at a specific event.
type TeamEventPerformance struct {
	TeamID    int
	TeamName  string
	Region    string
	EventID   string
	EventCode string
	EventName string
	OPR       float64
	NpOPR     float64
	CCWM      float64
	DPR       float64
	NpDPR     float64
	NpAVG     float64
	Matches   int
}

// TeamEventRankingsQuery retrieves performance metrics for teams at individual events.
// Unlike TeamRankingsQuery, this does not consolidate rankings across events - each team-event
// combination is returned as a separate entry.
func TeamEventRankingsQuery(region string, country string, eventCode string, year int) ([]TeamEventPerformance, error) {
	// Build team filter
	var teamFilter database.TeamFilter
	if region != "" {
		teamFilter.HomeRegions = []string{region}
	}
	if country != "" {
		teamFilter.Countries = []string{country}
	}
	if eventCode != "" {
		teamFilter.EventCodes = []string{eventCode}
	}

	// Get all teams based on filters
	var teams []*database.Team
	if region == "" && country == "" && eventCode == "" {
		teams = db.GetAllTeams()
	} else {
		teams = db.GetAllTeams(teamFilter)
	}
	if len(teams) == 0 {
		if region != "" {
			return nil, fmt.Errorf("no teams found in region %s", region)
		}
		if country != "" {
			return nil, fmt.Errorf("no teams found in country %s", country)
		}
		return nil, fmt.Errorf("no teams found")
	}

	// Get team info and build a map for easy lookup
	teamMap := make(map[int]*database.Team)
	teamIDs := make([]int, 0, len(teams))
	for _, t := range teams {
		teamMap[t.TeamID] = t
		teamIDs = append(teamIDs, t.TeamID)
	}

	// Build event filter
	eventFilter := database.EventFilter{Year: year}
	if region != "" {
		eventFilter.RegionCodes = []string{region}
	}
	if eventCode != "" {
		eventFilter.EventCodes = []string{eventCode}
	} else {
		// When no specific event is specified, only include qualifiers and championships
		eventFilter.Types = []string{"2", "4"}
	}
	events := db.GetAllEvents(eventFilter)
	if len(events) == 0 {
		return nil, fmt.Errorf("no events found")
	}

	// Build event map for easy lookup
	eventMap := make(map[string]*database.Event)
	eventIDs := make([]string, 0, len(events))
	for _, event := range events {
		eventMap[event.EventID] = event
		eventIDs = append(eventIDs, event.EventID)
	}

	// Get all team rankings for these teams and events
	rankingFilter := database.TeamRankingFilter{
		TeamIDs:  teamIDs,
		EventIDs: eventIDs,
	}
	rankings := db.GetTeamRankings(rankingFilter)
	if len(rankings) == 0 {
		if region != "" {
			return nil, fmt.Errorf("no team rankings found for teams in region %s for year %d", region, year)
		}
		return nil, fmt.Errorf("no team rankings found for year %d", year)
	}

	// Create a result for each ranking (no consolidation)
	results := make([]TeamEventPerformance, 0, len(rankings))
	for _, ranking := range rankings {
		team := teamMap[ranking.TeamID]
		event := eventMap[ranking.EventID]

		results = append(results, TeamEventPerformance{
			TeamID:    ranking.TeamID,
			TeamName:  team.Name,
			Region:    team.HomeRegion,
			EventID:   ranking.EventID,
			EventCode: event.EventCode,
			EventName: event.Name,
			OPR:       ranking.OPR,
			NpOPR:     ranking.NpOPR,
			CCWM:      ranking.CCWM,
			DPR:       ranking.DPR,
			NpDPR:     ranking.NpDPR,
			NpAVG:     ranking.NpAvg,
			Matches:   ranking.NumMatches,
		})
	}

	// Sort by NpAVG (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].NpAVG > results[j].NpAVG
	})

	return results, nil
}
