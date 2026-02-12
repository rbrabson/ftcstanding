package database

import (
	"slices"
	"sort"
)

// GetEvent retrieves an event from the file database by its ID.
func (db *filedb) GetEvent(eventID string) *Event {
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()

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
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()

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
			if !slices.Contains(filter.EventCodes, event.EventCode) {
				matchesFilter = false
			}
		}

		// Check RegionCode filter (OR within field)
		if matchesFilter && len(filter.RegionCodes) > 0 {
			if !slices.Contains(filter.RegionCodes, event.RegionCode) {
				matchesFilter = false
			}
		}

		// Check Country filter (OR within field)
		if matchesFilter && len(filter.Countries) > 0 {
			if !slices.Contains(filter.Countries, event.Country) {
				matchesFilter = false
			}
		}

		// Check Year filter
		if matchesFilter && filter.Year > 0 {
			if event.Year != filter.Year {
				matchesFilter = false
			}
		}

		// Check Type filter (OR within field)
		if matchesFilter && len(filter.Types) > 0 {
			if !slices.Contains(filter.Types, event.Type) {
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
	db.eventsMu.Lock()
	defer db.eventsMu.Unlock()

	// Make a copy to avoid external modifications
	eventCopy := *event
	db.events[event.EventID] = &eventCopy

	// Persist to disk
	return db.saveJSONFile("events.json", db.events)
}

// GetEventAwards retrieves all awards given at a specific event.
func (db *filedb) GetEventAwards(eventID string) []*EventAward {
	db.eventAwardsMu.RLock()
	defer db.eventAwardsMu.RUnlock()

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
	db.eventAwardsMu.Lock()
	defer db.eventAwardsMu.Unlock()

	// Check if this award already exists for this event/team/award combination
	awards := db.eventAwards[ea.EventID]
	found := false
	for i, existing := range awards {
		if existing.TeamID == ea.TeamID && existing.AwardID == ea.AwardID && existing.Series == ea.Series {
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
	db.eventAwardsMu.RLock()
	defer db.eventAwardsMu.RUnlock()

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
	db.eventAwardsMu.RLock()
	defer db.eventAwardsMu.RUnlock()

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
	db.eventRankingsMu.RLock()
	defer db.eventRankingsMu.RUnlock()

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
	db.eventRankingsMu.Lock()
	defer db.eventRankingsMu.Unlock()

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
	db.eventAdvancementsMu.RLock()
	defer db.eventAdvancementsMu.RUnlock()

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
	db.eventAdvancementsMu.Lock()
	defer db.eventAdvancementsMu.Unlock()

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
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()

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
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()

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
	// Need to lock both events and eventAdvancements since we read from both
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()
	db.eventAdvancementsMu.RLock()
	defer db.eventAdvancementsMu.RUnlock()

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

// GetAllAdvancements retrieves all event advancements from all events with optional filters.
// If no filters are provided, returns all advancements.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *filedb) GetAllAdvancements(filters ...AdvancementFilter) []*EventAdvancement {
	// Lock eventAdvancements for all cases
	db.eventAdvancementsMu.RLock()
	defer db.eventAdvancementsMu.RUnlock()

	// If no filters, return all advancements
	if len(filters) == 0 {
		result := make([]*EventAdvancement, 0)
		for _, advancements := range db.eventAdvancements {
			for _, advancement := range advancements {
				advancementCopy := *advancement
				result = append(result, &advancementCopy)
			}
		}
		return result
	}

	// Need to also lock events when filtering
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()

	filter := filters[0]
	result := make([]*EventAdvancement, 0)

	for eventID, advancements := range db.eventAdvancements {
		event, ok := db.events[eventID]
		if !ok {
			continue
		}

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

		if matchesFilter {
			for _, advancement := range advancements {
				advancementCopy := *advancement
				result = append(result, &advancementCopy)
			}
		}
	}

	return result
}

// GetEventTeams retrieves all teams for a specific event.
func (db *filedb) GetEventTeams(eventID string) []*EventTeam {
	db.eventTeamsMu.RLock()
	defer db.eventTeamsMu.RUnlock()

	teams, ok := db.eventTeams[eventID]
	if !ok {
		return nil
	}

	// Return copies
	result := make([]*EventTeam, len(teams))
	for i, team := range teams {
		teamCopy := *team
		result[i] = &teamCopy
	}
	return result
}

// SaveEventTeam saves or updates an event team in the file database.
func (db *filedb) SaveEventTeam(et *EventTeam) error {
	db.eventTeamsMu.Lock()
	defer db.eventTeamsMu.Unlock()

	// Check if this team already exists for this event
	teams := db.eventTeams[et.EventID]
	found := false
	for i, existing := range teams {
		if existing.TeamID == et.TeamID {
			// Update existing
			etCopy := *et
			teams[i] = &etCopy
			found = true
			break
		}
	}

	if !found {
		// Add new
		etCopy := *et
		db.eventTeams[et.EventID] = append(teams, &etCopy)
	}

	// Persist to disk
	return db.saveJSONFile("event_teams.json", db.eventTeams)
}

// GetEventsByTeam retrieves all event IDs that a team has or will participate in.
func (db *filedb) GetEventsByTeam(teamID int) []string {
	db.eventTeamsMu.RLock()
	defer db.eventTeamsMu.RUnlock()

	var eventIDs []string
	for eventID, teams := range db.eventTeams {
		for _, team := range teams {
			if team.TeamID == teamID {
				eventIDs = append(eventIDs, eventID)
				break
			}
		}
	}
	return eventIDs
}
