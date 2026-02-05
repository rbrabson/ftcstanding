package database

// GetAward retrieves an award from the file database by its ID.
func (db *filedb) GetAward(awardID int) *Award {
	db.mu.RLock()
	defer db.mu.RUnlock()

	award, ok := db.awards[awardID]
	if !ok {
		return nil
	}
	// Return a copy to avoid external modifications
	awardCopy := *award
	return &awardCopy
}

// GetAllAwards retrieves all awards from the file database.
func (db *filedb) GetAllAwards() []*Award {
	db.mu.RLock()
	defer db.mu.RUnlock()

	awards := make([]*Award, 0, len(db.awards))
	for _, award := range db.awards {
		awardCopy := *award
		awards = append(awards, &awardCopy)
	}
	return awards
}

// SaveAward saves or updates an award in the file database.
func (db *filedb) SaveAward(award *Award) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Make a copy to avoid external modifications
	awardCopy := *award
	db.awards[award.AwardID] = &awardCopy

	// Persist to disk
	return db.saveJSONFile("awards.json", db.awards)
}
