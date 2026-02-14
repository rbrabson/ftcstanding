package database

// GetAward retrieves an award from the file database by its ID.
func (db *filedb) GetAward(awardID int) (*Award, error) {
	if err := db.refreshAwardsIfChanged(); err != nil {
		return nil, err
	}

	db.awardsMu.RLock()
	defer db.awardsMu.RUnlock()

	award, ok := db.awards[awardID]
	if !ok {
		return nil, nil
	}
	// Return a copy to avoid external modifications
	awardCopy := *award
	return &awardCopy, nil
}

// GetAllAwards retrieves all awards from the file database.
func (db *filedb) GetAllAwards() ([]*Award, error) {
	if err := db.refreshAwardsIfChanged(); err != nil {
		return nil, err
	}

	db.awardsMu.RLock()
	defer db.awardsMu.RUnlock()

	awards := make([]*Award, 0, len(db.awards))
	for _, award := range db.awards {
		awardCopy := *award
		awards = append(awards, &awardCopy)
	}
	return awards, nil
}

// SaveAward saves or updates an award in the file database.
func (db *filedb) SaveAward(award *Award) error {
	if err := db.refreshAwardsIfChanged(); err != nil {
		return err
	}

	db.awardsMu.Lock()
	defer db.awardsMu.Unlock()

	// Make a copy to avoid external modifications
	awardCopy := *award
	db.awards[award.AwardID] = &awardCopy

	// Persist to disk
	return db.saveJSONFile("awards.json", db.awards)
}
