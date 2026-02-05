package database

import (
	"fmt"
	"time"
)

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
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID)
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
	_, err := stmt.Exec(ea.EventID, ea.TeamID, ea.AwardID)
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
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID)
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
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID)
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
	_, err := stmt.Exec(ea.EventID, ea.TeamID)
	return err
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

// GetAllAdvancements retrieves all event advancements from all events, ordered by event ID and team ID.
func (db *sqldb) GetAllAdvancements() []*EventAdvancement {
	stmt := db.getStatement("getAllAdvancements")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query()
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
