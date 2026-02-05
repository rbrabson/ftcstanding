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

// GetAllTeams retrieves all teams from the file database.
func (db *filedb) GetAllTeams() []*Team {
	db.mu.RLock()
	defer db.mu.RUnlock()

	teams := make([]*Team, 0, len(db.teams))
	for _, team := range db.teams {
		teamCopy := *team
		teams = append(teams, &teamCopy)
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
