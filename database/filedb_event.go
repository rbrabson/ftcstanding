package database

import "sort"

// GetEvent retrieves an event from the file database by its ID.
func (db *filedb) GetEvent(eventID string) *Event {
	db.mu.RLock()
	defer db.mu.RUnlock()

	event, ok := db.events[eventID]
	if !ok {
		return nil
	}
	// Return a copy to avoid external modifications
	eventCopy := *event
	return &eventCopy
}

// GetAllEvents retrieves all events from the file database with optional filters.
// If no filters are provided, returns all events.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *filedb) GetAllEvents(filters ...EventFilter) []*Event {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// If no filters, return all events
	if len(filters) == 0 {
		events := make([]*Event, 0, len(db.events))
		for _, event := range db.events {
			eventCopy := *event
			events = append(events, &eventCopy)
		}
		return events
	}

	filter := filters[0]
	events := make([]*Event, 0)

	for _, event := range db.events {
		// Apply filters with AND logic between different filter types
		matchesFilter := true

		// Check EventCode filter (OR within field)
		if len(filter.EventCodes) > 0 {
			found := false
			for _, code := range filter.EventCodes {
				if event.EventCode == code {
					found = true
					break
				}
			}
			if !found {
				matchesFilter = false
			}
		}

		// Check RegionCode filter (OR within field)
		if matchesFilter && len(filter.RegionCodes) > 0 {
			found := false
			for _, code := range filter.RegionCodes {
				if event.RegionCode == code {
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
				if event.Country == country {
					found = true
					break
				}
			}
			if !found {
				matchesFilter = false
			}
		}

		if matchesFilter {
			eventCopy := *event
			events = append(events, &eventCopy)
		}
	}

	return events
}

// SaveEvent saves or updates an event in the file database.
func (db *filedb) SaveEvent(event *Event) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Make a copy to avoid external modifications
	eventCopy := *event
	db.events[event.EventID] = &eventCopy

	// Persist to disk
	return db.saveJSONFile("events.json", db.events)
}

// GetEventAwards retrieves all awards given at a specific event.
func (db *filedb) GetEventAwards(eventID string) []*EventAward {
	db.mu.RLock()
	defer db.mu.RUnlock()

	awards, ok := db.eventAwards[eventID]
	if !ok {
		return nil
	}

	// Return copies
	result := make([]*EventAward, len(awards))
	for i, award := range awards {
		awardCopy := *award
		result[i] = &awardCopy
	}
	return result
}

// SaveEventAward saves or updates an event award in the file database.
func (db *filedb) SaveEventAward(ea *EventAward) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if this award already exists for this event/team
	awards := db.eventAwards[ea.EventID]
	found := false
	for i, existing := range awards {
		if existing.TeamID == ea.TeamID && existing.AwardID == ea.AwardID {
			// Update existing
			eaCopy := *ea
			awards[i] = &eaCopy
			found = true
			break
		}
	}

	if !found {
		// Add new
		eaCopy := *ea
		db.eventAwards[ea.EventID] = append(awards, &eaCopy)
	}

	// Persist to disk
	return db.saveJSONFile("event_awards.json", db.eventAwards)
}

// GetTeamAwardsByEvent retrieves all awards for a specific team at a specific event.
func (db *filedb) GetTeamAwardsByEvent(eventID string, teamID int) []*EventAward {
	db.mu.RLock()
	defer db.mu.RUnlock()

	awards, ok := db.eventAwards[eventID]
	if !ok {
		return nil
	}

	result := make([]*EventAward, 0)
	for _, award := range awards {
		if award.TeamID == teamID {
			awardCopy := *award
			result = append(result, &awardCopy)
		}
	}
	return result
}

// GetAllTeamAwards retrieves all awards for a specific team across all events.
func (db *filedb) GetAllTeamAwards(teamID int) []*EventAward {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make([]*EventAward, 0)
	for _, awards := range db.eventAwards {
		for _, award := range awards {
			if award.TeamID == teamID {
				awardCopy := *award
				result = append(result, &awardCopy)
			}
		}
	}
	return result
}

// GetEventRankings retrieves all rankings for a specific event.
func (db *filedb) GetEventRankings(eventID string) []*EventRanking {
	db.mu.RLock()
	defer db.mu.RUnlock()

	rankings, ok := db.eventRankings[eventID]
	if !ok {
		return nil
	}

	// Return copies
	result := make([]*EventRanking, len(rankings))
	for i, ranking := range rankings {
		rankingCopy := *ranking
		result[i] = &rankingCopy
	}
	return result
}

// SaveEventRanking saves or updates an event ranking in the file database.
func (db *filedb) SaveEventRanking(er *EventRanking) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if this ranking already exists for this event/team
	rankings := db.eventRankings[er.EventID]
	found := false
	for i, existing := range rankings {
		if existing.TeamID == er.TeamID {
			// Update existing
			erCopy := *er
			rankings[i] = &erCopy
			found = true
			break
		}
	}

	if !found {
		// Add new
		erCopy := *er
		db.eventRankings[er.EventID] = append(rankings, &erCopy)
	}

	// Persist to disk
	return db.saveJSONFile("event_rankings.json", db.eventRankings)
}

// GetEventAdvancements retrieves all team advancements for a specific event.
func (db *filedb) GetEventAdvancements(eventID string) []*EventAdvancement {
	db.mu.RLock()
	defer db.mu.RUnlock()

	advancements, ok := db.eventAdvancements[eventID]
	if !ok {
		return nil
	}

	// Return copies
	result := make([]*EventAdvancement, len(advancements))
	for i, advancement := range advancements {
		advancementCopy := *advancement
		result[i] = &advancementCopy
	}
	return result
}

// SaveEventAdvancement saves or updates an event advancement in the file database.
func (db *filedb) SaveEventAdvancement(ea *EventAdvancement) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if this advancement already exists for this event/team
	advancements := db.eventAdvancements[ea.EventID]
	found := false
	for i, existing := range advancements {
		if existing.TeamID == ea.TeamID {
			// Update existing
			eaCopy := *ea
			advancements[i] = &eaCopy
			found = true
			break
		}
	}

	if !found {
		// Add new
		eaCopy := *ea
		db.eventAdvancements[ea.EventID] = append(advancements, &eaCopy)
	}

	// Persist to disk
	return db.saveJSONFile("event_advancements.json", db.eventAdvancements)
}

// GetRegionCodes retrieves all unique region codes from events.
func (db *filedb) GetRegionCodes() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	regionMap := make(map[string]bool)
	for _, event := range db.events {
		if event.RegionCode != "" {
			regionMap[event.RegionCode] = true
		}
	}

	regions := make([]string, 0, len(regionMap))
	for region := range regionMap {
		regions = append(regions, region)
	}
	sort.Strings(regions)
	return regions
}

// GetEventCodesByRegion retrieves all unique event codes for a given region.
func (db *filedb) GetEventCodesByRegion(regionCode string) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	eventCodeMap := make(map[string]bool)
	for _, event := range db.events {
		if event.RegionCode == regionCode {
			eventCodeMap[event.EventCode] = true
		}
	}

	eventCodes := make([]string, 0, len(eventCodeMap))
	for code := range eventCodeMap {
		eventCodes = append(eventCodes, code)
	}
	sort.Strings(eventCodes)
	return eventCodes
}

// GetAdvancementsByRegion retrieves all event advancements for events in a given region.
func (db *filedb) GetAdvancementsByRegion(regionCode string) []*EventAdvancement {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make([]*EventAdvancement, 0)
	for eventID, advancements := range db.eventAdvancements {
		event, ok := db.events[eventID]
		if ok && event.RegionCode == regionCode {
			for _, advancement := range advancements {
				advancementCopy := *advancement
				result = append(result, &advancementCopy)
			}
		}
	}
	return result
}

// GetAllAdvancements retrieves all event advancements from all events.
func (db *filedb) GetAllAdvancements() []*EventAdvancement {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make([]*EventAdvancement, 0)
	for _, advancements := range db.eventAdvancements {
		for _, advancement := range advancements {
			advancementCopy := *advancement
			result = append(result, &advancementCopy)
		}
	}
	return result
}
