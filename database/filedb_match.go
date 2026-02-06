package database

// GetMatch retrieves a match from the file database by its ID.
func (db *filedb) GetMatch(matchID string) *Match {
	db.mu.RLock()
	defer db.mu.RUnlock()

	match, ok := db.matches[matchID]
	if !ok {
		return nil
	}
	// Return a copy to avoid external modifications
	matchCopy := *match
	return &matchCopy
}

// GetAllMatches retrieves all matches from the file database with optional filters.
// If no filters are provided, returns all matches.
// Filters are combined with OR logic within each field.
func (db *filedb) GetAllMatches(filters ...MatchFilter) []*Match {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// If no filters, return all matches
	if len(filters) == 0 {
		matches := make([]*Match, 0, len(db.matches))
		for _, match := range db.matches {
			matchCopy := *match
			matches = append(matches, &matchCopy)
		}
		return matches
	}

	filter := filters[0]
	matches := make([]*Match, 0)

	for _, match := range db.matches {
		matchesFilter := true

		// Check EventID filter (OR within field)
		if len(filter.EventIDs) > 0 {
			found := false
			for _, eventID := range filter.EventIDs {
				if match.EventID == eventID {
					found = true
					break
				}
			}
			if !found {
				matchesFilter = false
			}
		}

		if matchesFilter {
			matchCopy := *match
			matches = append(matches, &matchCopy)
		}
	}

	return matches
}

// GetMatchesByEvent retrieves all matches for a specific event.
func (db *filedb) GetMatchesByEvent(eventID string) []*Match {
	db.mu.RLock()
	defer db.mu.RUnlock()

	matches := make([]*Match, 0)
	for _, match := range db.matches {
		if match.EventID == eventID {
			matchCopy := *match
			matches = append(matches, &matchCopy)
		}
	}
	return matches
}

// SaveMatch saves or updates a match in the file database.
func (db *filedb) SaveMatch(match *Match) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Make a copy to avoid external modifications
	matchCopy := *match
	db.matches[match.MatchID] = &matchCopy

	// Persist to disk
	return db.saveJSONFile("matches.json", db.matches)
}

// GetMatchAllianceScore retrieves the score for a specific alliance in a match.
func (db *filedb) GetMatchAllianceScore(matchID, alliance string) *MatchAllianceScore {
	db.mu.RLock()
	defer db.mu.RUnlock()

	matchScores, ok := db.matchScores[matchID]
	if !ok {
		return nil
	}

	score, ok := matchScores[alliance]
	if !ok {
		return nil
	}

	// Return a copy to avoid external modifications
	scoreCopy := *score
	return &scoreCopy
}

// SaveMatchAllianceScore saves or updates the score for a specific alliance in a match.
func (db *filedb) SaveMatchAllianceScore(score *MatchAllianceScore) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Ensure the match score map exists
	if db.matchScores[score.MatchID] == nil {
		db.matchScores[score.MatchID] = make(map[string]*MatchAllianceScore)
	}

	// Make a copy to avoid external modifications
	scoreCopy := *score
	db.matchScores[score.MatchID][score.Alliance] = &scoreCopy

	// Persist to disk
	return db.saveJSONFile("match_scores.json", db.matchScores)
}

// GetMatchTeams retrieves all teams participating in a specific match.
func (db *filedb) GetMatchTeams(matchID string) []*MatchTeam {
	db.mu.RLock()
	defer db.mu.RUnlock()

	teams, ok := db.matchTeams[matchID]
	if !ok {
		return nil
	}

	// Return copies
	result := make([]*MatchTeam, len(teams))
	for i, team := range teams {
		teamCopy := *team
		result[i] = &teamCopy
	}
	return result
}

// SaveMatchTeam saves or updates a match team in the file database.
func (db *filedb) SaveMatchTeam(team *MatchTeam) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if this team already exists for this match
	teams := db.matchTeams[team.MatchID]
	found := false
	for i, existing := range teams {
		if existing.TeamID == team.TeamID {
			// Update existing
			teamCopy := *team
			teams[i] = &teamCopy
			found = true
			break
		}
	}

	if !found {
		// Add new
		teamCopy := *team
		db.matchTeams[team.MatchID] = append(teams, &teamCopy)
	}

	// Persist to disk
	return db.saveJSONFile("match_teams.json", db.matchTeams)
}

// GetTeamsByEvent retrieves all unique team IDs that participated at a specific event.
func (db *filedb) GetTeamsByEvent(eventID string) []int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	teamIDMap := make(map[int]bool)

	// Iterate through all matches to find those belonging to this event
	for matchID, match := range db.matches {
		if match.EventID == eventID {
			// Get all teams for this match
			teams, ok := db.matchTeams[matchID]
			if ok {
				for _, team := range teams {
					teamIDMap[team.TeamID] = true
				}
			}
		}
	}

	// Convert map to sorted slice
	teamIDs := make([]int, 0, len(teamIDMap))
	for teamID := range teamIDMap {
		teamIDs = append(teamIDs, teamID)
	}
	return teamIDs
}
