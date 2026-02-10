package database

import (
	"slices"
	"sort"
)

// GetTeam retrieves a team from the file database by its ID.
func (db *filedb) GetTeam(teamID int) *Team {
	db.teamsMu.RLock()
	defer db.teamsMu.RUnlock()

	team, ok := db.teams[teamID]
	if !ok {
		return nil
	}
	// Return a copy to avoid external modifications
	teamCopy := *team
	return &teamCopy
}

// GetAllTeams retrieves all teams from the file database with optional filters.
// If no filters are provided, returns all teams.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *filedb) GetAllTeams(filters ...TeamFilter) []*Team {
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
		return teams
	}

	filter := filters[0]

	// If EventCodes filter is provided, get team IDs from those events
	var eventTeamIDs map[int]bool
	if len(filter.EventCodes) > 0 {
		eventTeamIDs = make(map[int]bool)
		for _, eventCode := range filter.EventCodes {
			// Get all events matching this code
			events := db.GetAllEvents(EventFilter{EventCodes: []string{eventCode}})
			for _, event := range events {
				// Get all teams for this event
				eventTeams := db.GetEventTeams(event.EventID)
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
	return teams
}

// SaveTeam saves or updates a team in the file database.
func (db *filedb) SaveTeam(team *Team) error {
	db.teamsMu.Lock()
	defer db.teamsMu.Unlock()

	// Make a copy to avoid external modifications
	teamCopy := *team
	db.teams[team.TeamID] = &teamCopy

	// Persist to disk
	return db.saveJSONFile("teams.json", db.teams)
}

// GetTeamsByRegion retrieves all teams in a given home region from the file database.
func (db *filedb) GetTeamsByRegion(region string) []*Team {
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
	return teams
}
