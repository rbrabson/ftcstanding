package database

// GetTeam retrieves a team from the file database by its ID.
func (db *filedb) GetTeam(teamID int) *Team {
	db.mu.RLock()
	defer db.mu.RUnlock()

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
	db.mu.RLock()
	defer db.mu.RUnlock()

	// If no filters, return all teams
	if len(filters) == 0 {
		teams := make([]*Team, 0, len(db.teams))
		for _, team := range db.teams {
			teamCopy := *team
			teams = append(teams, &teamCopy)
		}
		return teams
	}

	filter := filters[0]
	teams := make([]*Team, 0)

	for _, team := range db.teams {
		// Apply filters with AND logic between different filter types
		matchesFilter := true

		// Check TeamID filter (OR within field)
		if len(filter.TeamIDs) > 0 {
			found := false
			for _, id := range filter.TeamIDs {
				if team.TeamID == id {
					found = true
					break
				}
			}
			if !found {
				matchesFilter = false
			}
		}

		// Check Country filter (OR within field)
		if matchesFilter && len(filter.Countries) > 0 {
			found := false
			for _, country := range filter.Countries {
				if team.Country == country {
					found = true
					break
				}
			}
			if !found {
				matchesFilter = false
			}
		}

		// Check HomeRegion filter (OR within field)
		if matchesFilter && len(filter.HomeRegions) > 0 {
			found := false
			for _, region := range filter.HomeRegions {
				if team.HomeRegion == region {
					found = true
					break
				}
			}
			if !found {
				matchesFilter = false
			}
		}

		if matchesFilter {
			teamCopy := *team
			teams = append(teams, &teamCopy)
		}
	}

	return teams
}

// SaveTeam saves or updates a team in the file database.
func (db *filedb) SaveTeam(team *Team) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Make a copy to avoid external modifications
	teamCopy := *team
	db.teams[team.TeamID] = &teamCopy

	// Persist to disk
	return db.saveJSONFile("teams.json", db.teams)
}

// GetTeamsByRegion retrieves all teams in a given home region from the file database.
func (db *filedb) GetTeamsByRegion(region string) []*Team {
	db.mu.RLock()
	defer db.mu.RUnlock()

	teams := make([]*Team, 0)
	for _, team := range db.teams {
		if team.HomeRegion == region {
			teamCopy := *team
			teams = append(teams, &teamCopy)
		}
	}
	return teams
}
