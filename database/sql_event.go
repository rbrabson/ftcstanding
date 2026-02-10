package database

import (
	"fmt"
	"time"
)

// InitEventStatements prepares all SQL statements for event operations.
func (db *sqldb) initEventStatements() error {
	queries := map[string]string{
		"getEvent":                "SELECT event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end FROM events WHERE event_id = ?",
		"saveEvent":               "INSERT INTO events (event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE event_code = VALUES(event_code), year = VALUES(year), name = VALUES(name), type = VALUES(type), division_code = VALUES(division_code), region_code = VALUES(region_code), league_code = VALUES(league_code), venue = VALUES(venue), address = VALUES(address), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), timezone = VALUES(timezone), date_start = VALUES(date_start), date_end = VALUES(date_end)",
		"getEventAwards":          "SELECT event_id, team_id, award_id, name, series FROM event_awards WHERE event_id = ?",
		"saveEventAward":          "INSERT INTO event_awards (event_id, team_id, award_id, name, series) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), series = VALUES(series)",
		"getTeamAwardsByEvent":    "SELECT event_id, team_id, award_id, name, series FROM event_awards WHERE event_id = ? AND team_id = ?",
		"getAllTeamAwards":        "SELECT event_id, team_id, award_id, name, series FROM event_awards WHERE team_id = ? ORDER BY event_id",
		"getEventRankings":        "SELECT event_id, team_id, rank, sort_order1, sort_order2, sort_order3, sort_order4, sort_order5, sort_order6, wins, losses, ties, dq, matches_played, matches_counted FROM event_rankings WHERE event_id = ?",
		"saveEventRanking":        "INSERT INTO event_rankings (event_id, team_id, rank, sort_order1, sort_order2, sort_order3, sort_order4, sort_order5, sort_order6, wins, losses, ties, dq, matches_played, matches_counted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE rank = VALUES(rank), sort_order1 = VALUES(sort_order1), sort_order2 = VALUES(sort_order2), sort_order3 = VALUES(sort_order3), sort_order4 = VALUES(sort_order4), sort_order5 = VALUES(sort_order5), sort_order6 = VALUES(sort_order6), wins = VALUES(wins), losses = VALUES(losses), ties = VALUES(ties), dq = VALUES(dq), matches_played = VALUES(matches_played), matches_counted = VALUES(matches_counted)",
		"getEventAdvancements":    "SELECT event_id, team_id, status FROM event_advancements WHERE event_id = ?",
		"saveEventAdvancement":    "INSERT INTO event_advancements (event_id, team_id, status) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE status = VALUES(status)",
		"getEventTeams":           "SELECT event_id, team_id FROM event_teams WHERE event_id = ?",
		"saveEventTeam":           "INSERT INTO event_teams (event_id, team_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE event_id = event_id",
		"getEventsByTeam":         "SELECT DISTINCT event_id FROM event_teams WHERE team_id = ? ORDER BY event_id",
		"getAllAdvancements":      "SELECT event_id, team_id, status FROM event_advancements ORDER BY event_id, team_id",
		"getRegionCodes":          "SELECT DISTINCT region_code FROM events WHERE region_code IS NOT NULL AND region_code != '' ORDER BY region_code",
		"getEventCodesByRegion":   "SELECT DISTINCT event_code FROM events WHERE region_code = ? ORDER BY event_code",
		"getAdvancementsByRegion": "SELECT ea.event_id, ea.team_id, ea.status FROM event_advancements ea INNER JOIN events e ON ea.event_id = e.event_id WHERE e.region_code = ? ORDER BY ea.event_id, ea.team_id",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// GetEventID generates an EventID from the given EventCode and DateStart.
func (db *sqldb) GetEventID(eventCode string, dateStart time.Time) string {
	return fmt.Sprintf("%s : %04d-%02d-%02d", eventCode, dateStart.Year(), int(dateStart.Month()), dateStart.Day())
}

// GetEvent retrieves an event from the database by its ID.
func (db *sqldb) GetEvent(eventID string) *Event {
	var event Event
	stmt := db.getStatement("getEvent")
	if stmt == nil {
		return nil
	}
	err := stmt.QueryRow(eventID).Scan(
		&event.EventID,
		&event.EventCode,
		&event.Year,
		&event.Name,
		&event.Type,
		&event.DivisionCode,
		&event.RegionCode,
		&event.LeagueCode,
		&event.Venue,
		&event.Address,
		&event.City,
		&event.StateProv,
		&event.Country,
		&event.Timezone,
		&event.DateStart,
		&event.DateEnd,
	)
	if err != nil {
		return nil
	}
	return &event
}

// GetAllEvents retrieves all events from the database with optional filters.
// If no filters are provided, returns all events.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *sqldb) GetAllEvents(filters ...EventFilter) []*Event {
	// Build dynamic query
	query := "SELECT event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end FROM events"
	args := []interface{}{}

	if len(filters) > 0 {
		filter := filters[0]
		query += " WHERE 1=1"

		// Add EventCode filter
		if len(filter.EventCodes) > 0 {
			query += " AND event_code IN ("
			for i, code := range filter.EventCodes {
				if i > 0 {
					query += ","
				}
				query += "?"
				args = append(args, code)
			}
			query += ")"
		}

		// Add RegionCode filter
		if len(filter.RegionCodes) > 0 {
			query += " AND region_code IN ("
			for i, code := range filter.RegionCodes {
				if i > 0 {
					query += ","
				}
				query += "?"
				args = append(args, code)
			}
			query += ")"
		}

		// Add Country filter
		if len(filter.Countries) > 0 {
			query += " AND country IN ("
			for i, country := range filter.Countries {
				if i > 0 {
					query += ","
				}
				query += "?"
				args = append(args, country)
			}
			query += ")"
		}

		// Add Year filter
		if filter.Year > 0 {
			query += " AND year = ?"
			args = append(args, filter.Year)
		}
	}

	query += " ORDER BY date_start, event_code"

	// Execute query
	rows, err := db.sqldb.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.EventID,
			&event.EventCode,
			&event.Year,
			&event.Name,
			&event.Type,
			&event.DivisionCode,
			&event.RegionCode,
			&event.LeagueCode,
			&event.Venue,
			&event.Address,
			&event.City,
			&event.StateProv,
			&event.Country,
			&event.Timezone,
			&event.DateStart,
			&event.DateEnd,
		)
		if err != nil {
			continue
		}
		events = append(events, &event)
	}
	return events
}

// SaveEvent saves or updates an event in the
func (db *sqldb) SaveEvent(event *Event) error {
	stmt := db.getStatement("saveEvent")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(
		event.EventID,
		event.EventCode,
		event.Year,
		event.Name,
		event.Type,
		event.DivisionCode,
		event.RegionCode,
		event.LeagueCode,
		event.Venue,
		event.Address,
		event.City,
		event.StateProv,
		event.Country,
		event.Timezone,
		event.DateStart,
		event.DateEnd,
	)
	return err
}

// GetEventAwards retrieves all awards given at a specific event.
func (db *sqldb) GetEventAwards(eventID string) []*EventAward {
	stmt := db.getStatement("getEventAwards")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var awards []*EventAward
	for rows.Next() {
		var ea EventAward
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID, &ea.Name, &ea.Series)
		if err != nil {
			continue
		}
		awards = append(awards, &ea)
	}
	return awards
}

// SaveEventAward saves or updates an event award in the
func (db *sqldb) SaveEventAward(ea *EventAward) error {
	stmt := db.getStatement("saveEventAward")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(ea.EventID, ea.TeamID, ea.AwardID, ea.Name, ea.Series)
	return err
}

// GetTeamAwardsByEvent retrieves all awards for a specific team at a specific event.
func (db *sqldb) GetTeamAwardsByEvent(eventID string, teamID int) []*EventAward {
	stmt := db.getStatement("getTeamAwardsByEvent")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID, teamID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var awards []*EventAward
	for rows.Next() {
		var ea EventAward
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID, &ea.Name, &ea.Series)
		if err != nil {
			continue
		}
		awards = append(awards, &ea)
	}
	return awards
}

// GetAllTeamAwards retrieves all awards for a specific team across all events, ordered by event ID.
func (db *sqldb) GetAllTeamAwards(teamID int) []*EventAward {
	stmt := db.getStatement("getAllTeamAwards")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(teamID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var awards []*EventAward
	for rows.Next() {
		var ea EventAward
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID, &ea.Name, &ea.Series)
		if err != nil {
			continue
		}
		awards = append(awards, &ea)
	}
	return awards
}

// GetEventRankings retrieves all rankings for a specific event.
func (db *sqldb) GetEventRankings(eventID string) []*EventRanking {
	stmt := db.getStatement("getEventRankings")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var rankings []*EventRanking
	for rows.Next() {
		var er EventRanking
		err := rows.Scan(
			&er.EventID,
			&er.TeamID,
			&er.Rank,
			&er.SortOrder1,
			&er.SortOrder2,
			&er.SortOrder3,
			&er.SortOrder4,
			&er.SortOrder5,
			&er.SortOrder6,
			&er.Wins,
			&er.Losses,
			&er.Ties,
			&er.Dq,
			&er.MatchesPlayed,
			&er.MatchesCounted,
		)
		if err != nil {
			continue
		}
		rankings = append(rankings, &er)
	}
	return rankings
}

// SaveEventRanking saves or updates an event ranking in the
func (db *sqldb) SaveEventRanking(er *EventRanking) error {
	stmt := db.getStatement("saveEventRanking")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(er.EventID, er.TeamID, er.Rank, er.SortOrder1, er.SortOrder2, er.SortOrder3, er.SortOrder4, er.SortOrder5, er.SortOrder6, er.Wins, er.Losses, er.Ties, er.Dq, er.MatchesPlayed, er.MatchesCounted)
	return err
}

// GetEventAdvancements retrieves all team advancements for a specific event.
func (db *sqldb) GetEventAdvancements(eventID string) []*EventAdvancement {
	stmt := db.getStatement("getEventAdvancements")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var advancements []*EventAdvancement
	for rows.Next() {
		var ea EventAdvancement
		err := rows.Scan(&ea.EventID, &ea.TeamID)
		if err != nil {
			continue
		}
		advancements = append(advancements, &ea)
	}
	return advancements
}

// SaveEventAdvancement saves or updates an event advancement in the
func (db *sqldb) SaveEventAdvancement(ea *EventAdvancement) error {
	stmt := db.getStatement("saveEventAdvancement")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(ea.EventID, ea.TeamID, ea.Status)
	return err
}

// GetEventTeams retrieves all teams for a specific event.
func (db *sqldb) GetEventTeams(eventID string) []*EventTeam {
	stmt := db.getStatement("getEventTeams")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var teams []*EventTeam
	for rows.Next() {
		var et EventTeam
		err := rows.Scan(&et.EventID, &et.TeamID)
		if err != nil {
			continue
		}
		teams = append(teams, &et)
	}
	return teams
}

// SaveEventTeam saves or updates an event team in the database.
func (db *sqldb) SaveEventTeam(et *EventTeam) error {
	stmt := db.getStatement("saveEventTeam")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(et.EventID, et.TeamID)
	return err
}

// GetEventsByTeam retrieves all event IDs that a team has or will participate in, sorted alphabetically.
func (db *sqldb) GetEventsByTeam(teamID int) []string {
	stmt := db.getStatement("getEventsByTeam")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(teamID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var eventIDs []string
	for rows.Next() {
		var eventID string
		err := rows.Scan(&eventID)
		if err != nil {
			continue
		}
		eventIDs = append(eventIDs, eventID)
	}
	return eventIDs
}

// GetRegionCodes retrieves all unique region codes from events, sorted alphabetically.
func (db *sqldb) GetRegionCodes() []string {
	stmt := db.getStatement("getRegionCodes")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()

	var regionCodes []string
	for rows.Next() {
		var regionCode string
		err := rows.Scan(&regionCode)
		if err != nil {
			continue
		}
		regionCodes = append(regionCodes, regionCode)
	}
	return regionCodes
}

// GetEventCodesByRegion retrieves all unique event codes for a given region, sorted alphabetically.
func (db *sqldb) GetEventCodesByRegion(regionCode string) []string {
	stmt := db.getStatement("getEventCodesByRegion")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(regionCode)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var eventCodes []string
	for rows.Next() {
		var eventCode string
		err := rows.Scan(&eventCode)
		if err != nil {
			continue
		}
		eventCodes = append(eventCodes, eventCode)
	}
	return eventCodes
}

// GetAdvancementsByRegion retrieves all event advancements for events in a given region, ordered by event ID and team ID.
func (db *sqldb) GetAdvancementsByRegion(regionCode string) []*EventAdvancement {
	stmt := db.getStatement("getAdvancementsByRegion")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(regionCode)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var advancements []*EventAdvancement
	for rows.Next() {
		var ea EventAdvancement
		err := rows.Scan(&ea.EventID, &ea.TeamID)
		if err != nil {
			continue
		}
		advancements = append(advancements, &ea)
	}
	return advancements
}

// GetAllAdvancements retrieves all event advancements from all events with optional filters.
// If no filters are provided, returns all advancements.
// Filters are combined with OR logic within each field and AND logic between fields.
func (db *sqldb) GetAllAdvancements(filters ...AdvancementFilter) []*EventAdvancement {
	// Build dynamic query
	query := "SELECT ea.event_id, ea.team_id, ea.status FROM event_advancements ea"
	args := []interface{}{}

	if len(filters) > 0 {
		filter := filters[0]
		// Need to join with events table for filtering
		if len(filter.Countries) > 0 || len(filter.RegionCodes) > 0 || len(filter.EventCodes) > 0 {
			query += " INNER JOIN events e ON ea.event_id = e.event_id WHERE 1=1"

			// Add EventCode filter
			if len(filter.EventCodes) > 0 {
				query += " AND e.event_code IN ("
				for i, code := range filter.EventCodes {
					if i > 0 {
						query += ","
					}
					query += "?"
					args = append(args, code)
				}
				query += ")"
			}

			// Add Country filter
			if len(filter.Countries) > 0 {
				query += " AND e.country IN ("
				for i, country := range filter.Countries {
					if i > 0 {
						query += ","
					}
					query += "?"
					args = append(args, country)
				}
				query += ")"
			}

			// Add RegionCode filter
			if len(filter.RegionCodes) > 0 {
				query += " AND e.region_code IN ("
				for i, code := range filter.RegionCodes {
					if i > 0 {
						query += ","
					}
					query += "?"
					args = append(args, code)
				}
				query += ")"
			}
		}
	}

	query += " ORDER BY ea.event_id, ea.team_id"

	// Execute query
	rows, err := db.sqldb.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var advancements []*EventAdvancement
	for rows.Next() {
		var ea EventAdvancement
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.Status)
		if err != nil {
			continue
		}
		advancements = append(advancements, &ea)
	}
	return advancements
}
