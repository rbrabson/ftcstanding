package database

import (
	"slices"
	"sort"
)

// GetTeam retrieves a team from the file database by its ID.
func (db *filedb) GetTeam(teamID int) (*Team, error) {
	if err := db.refreshTeamsIfChanged(); err != nil {
		return nil, err
	}

	db.teamsMu.RLock()
	defer db.teamsMu.RUnlock()

	team, ok := db.teams[teamID]
	if !ok {
		return nil, nil
	}
	// Return a copy to avoid external modifications
	teamCopy := *team
	return &teamCopy, nil
}

// GetAllTeams retrieves all teams from the file database with optional filters.
// If no filters are provided, returns all teams.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *filedb) GetAllTeams(filters ...TeamFilter) ([]*Team, error) {
	if err := db.refreshTeamsIfChanged(); err != nil {
		return nil, err
	}

	db.teamsMu.RLock()
	defer db.teamsMu.RUnlock()

	// If no filters, return all teams
	if len(filters) == 0 {
		teams := make([]*Team, 0, len(db.teams))
		for _, team := range db.teams {
			teamCopy := *team
			teams = append(teams, &teamCopy)
		}
		// Sort by TeamID
		sort.Slice(teams, func(i, j int) bool {
			return teams[i].TeamID < teams[j].TeamID
		})
		return teams, nil
	}

	filter := filters[0]

	// If EventCodes filter is provided, get team IDs from those events
	var eventTeamIDs map[int]bool
	if len(filter.EventCodes) > 0 {
		eventTeamIDs = make(map[int]bool)
		for _, eventCode := range filter.EventCodes {
			// Get all events matching this code
			events, err := db.GetAllEvents(EventFilter{EventCodes: []string{eventCode}})
			if err != nil {
				return nil, err
			}
			for _, event := range events {
				// Get all teams for this event
				eventTeams, err := db.GetEventTeams(event.EventID)
				if err != nil {
					return nil, err
				}
				for _, et := range eventTeams {
					eventTeamIDs[et.TeamID] = true
				}
			}
		}
	}

	teams := make([]*Team, 0)

	for _, team := range db.teams {
		// Apply filters with AND logic between different filter types
		matchesFilter := true

		// Check EventCodes filter (team must be in at least one of the events)
		if matchesFilter && len(filter.EventCodes) > 0 {
			if !eventTeamIDs[team.TeamID] {
				matchesFilter = false
			}
		}

		// Check TeamID filter (OR within field)
		if matchesFilter && len(filter.TeamIDs) > 0 {
			if !slices.Contains(filter.TeamIDs, team.TeamID) {
				matchesFilter = false
			}
		}

		// Check Country filter (OR within field)
		if matchesFilter && len(filter.Countries) > 0 {
			if !slices.Contains(filter.Countries, team.Country) {
				matchesFilter = false
			}
		}

		// Check HomeRegion filter (OR within field)
		if matchesFilter && len(filter.HomeRegions) > 0 {
			if !slices.Contains(filter.HomeRegions, team.HomeRegion) {
				matchesFilter = false
			}
		}

		if matchesFilter {
			teamCopy := *team
			teams = append(teams, &teamCopy)
		}
	}

	// Sort by TeamID
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].TeamID < teams[j].TeamID
	})
	return teams, nil
}

// SaveTeam saves or updates a team in the file database.
func (db *filedb) SaveTeam(team *Team) error {
	if err := db.refreshTeamsIfChanged(); err != nil {
		return err
	}

	db.teamsMu.Lock()
	defer db.teamsMu.Unlock()

	// Make a copy to avoid external modifications
	teamCopy := *team
	db.teams[team.TeamID] = &teamCopy

	// Persist to disk
	return db.saveJSONFile("teams.json", db.teams)
}

// GetTeamsByRegion retrieves all teams in a given home region from the file database.
func (db *filedb) GetTeamsByRegion(region string) ([]*Team, error) {
	if err := db.refreshTeamsIfChanged(); err != nil {
		return nil, err
	}

	db.teamsMu.RLock()
	defer db.teamsMu.RUnlock()

	teams := make([]*Team, 0)
	for _, team := range db.teams {
		if team.HomeRegion == region {
			teamCopy := *team
			teams = append(teams, &teamCopy)
		}
	}
	// Sort by TeamID
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].TeamID < teams[j].TeamID
	})
	return teams, nil
}

// GetTeamRankings retrieves team rankings with optional filters.
// Filters support filtering by TeamID and/or EventID.
// If no filters are provided, returns all team rankings.
func (db *filedb) GetTeamRankings(filters ...TeamRankingFilter) ([]*TeamRanking, error) {
	if err := db.refreshTeamRankingsIfChanged(); err != nil {
		return nil, err
	}

	db.teamRankingsMu.RLock()
	defer db.teamRankingsMu.RUnlock()

	var rankings []*TeamRanking

	// Helper function to check if a ranking matches the filter
	matchesFilter := func(ranking *TeamRanking, filter TeamRankingFilter) bool {
		// Check TeamID filter
		if len(filter.TeamIDs) > 0 {
			match := slices.Contains(filter.TeamIDs, ranking.TeamID)
			if !match {
				return false
			}
		}

		// Check EventID filter
		if len(filter.EventIDs) > 0 {
			match := slices.Contains(filter.EventIDs, ranking.EventID)
			if !match {
				return false
			}
		}

		return true
	}

	// If no filters, return all rankings
	if len(filters) == 0 {
		for _, eventRankings := range db.teamRankings {
			for _, ranking := range eventRankings {
				rankingCopy := *ranking
				rankings = append(rankings, &rankingCopy)
			}
		}
	} else {
		filter := filters[0]
		for _, eventRankings := range db.teamRankings {
			for _, ranking := range eventRankings {
				if matchesFilter(ranking, filter) {
					rankingCopy := *ranking
					rankings = append(rankings, &rankingCopy)
				}
			}
		}
	}

	// Sort by EventID then TeamID
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].EventID != rankings[j].EventID {
			return rankings[i].EventID < rankings[j].EventID
		}
		return rankings[i].TeamID < rankings[j].TeamID
	})

	return rankings, nil
}

// SaveTeamRanking saves or updates a team ranking in the file database.
func (db *filedb) SaveTeamRanking(ranking *TeamRanking) error {
	if err := db.refreshTeamRankingsIfChanged(); err != nil {
		return err
	}

	db.teamRankingsMu.Lock()
	defer db.teamRankingsMu.Unlock()

	// Initialize the map for this event if it doesn't exist
	if db.teamRankings[ranking.EventID] == nil {
		db.teamRankings[ranking.EventID] = make(map[int]*TeamRanking)
	}

	// Make a copy and save it
	rankingCopy := *ranking
	db.teamRankings[ranking.EventID][ranking.TeamID] = &rankingCopy

	// Persist to disk
	return db.saveJSONFile("team_rankings.json", db.teamRankings)
}
