package database

import (
	"fmt"
	"time"
)

// InitEventStatements prepares all SQL statements for event operations.
func InitEventStatements() error {
	queries := map[string]string{
		"getEvent":             "SELECT event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end FROM events WHERE event_id = ?",
		"saveEvent":            "INSERT INTO events (event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE event_code = VALUES(event_code), year = VALUES(year), name = VALUES(name), type = VALUES(type), division_code = VALUES(division_code), region_code = VALUES(region_code), league_code = VALUES(league_code), venue = VALUES(venue), address = VALUES(address), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), timezone = VALUES(timezone), date_start = VALUES(date_start), date_end = VALUES(date_end)",
		"getEventAwards":       "SELECT event_id, team_id, award_id FROM event_awards WHERE event_id = ?",
		"saveEventAward":       "INSERT INTO event_awards (event_id, team_id, award_id) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE event_id = event_id",
		"getEventRankings":     "SELECT event_id, team_id, rank, sort_order1, sort_order2, sort_order3, sort_order4, sort_order5, sort_order6, wins, losses, ties, dq, matches_played, matches_counted FROM event_rankings WHERE event_id = ?",
		"saveEventRanking":     "INSERT INTO event_rankings (event_id, team_id, rank, sort_order1, sort_order2, sort_order3, sort_order4, sort_order5, sort_order6, wins, losses, ties, dq, matches_played, matches_counted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE rank = VALUES(rank), sort_order1 = VALUES(sort_order1), sort_order2 = VALUES(sort_order2), sort_order3 = VALUES(sort_order3), sort_order4 = VALUES(sort_order4), sort_order5 = VALUES(sort_order5), sort_order6 = VALUES(sort_order6), wins = VALUES(wins), losses = VALUES(losses), ties = VALUES(ties), dq = VALUES(dq), matches_played = VALUES(matches_played), matches_counted = VALUES(matches_counted)",
		"getEventAdvancements": "SELECT event_id, team_id FROM event_advancements WHERE event_id = ?",
		"saveEventAdvancement": "INSERT INTO event_advancements (event_id, team_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE event_id = event_id",
		"getRegionCodes":       "SELECT DISTINCT region_code FROM events WHERE region_code IS NOT NULL AND region_code != '' ORDER BY region_code",
	}

	for name, query := range queries {
		if err := PrepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// Event represents a competition event.
type Event struct {
	EventID      string    `json:"event_id"`
	EventCode    string    `json:"event_code"`
	Year         int       `json:"year"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	DivisionCode string    `json:"division_code"`
	RegionCode   string    `json:"region_code"`
	LeagueCode   string    `json:"league_code"`
	Venue        string    `json:"venue"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	StateProv    string    `json:"state_prov"`
	Country      string    `json:"country"`
	Timezone     string    `json:"timezone"`
	DateStart    time.Time `json:"date_start"`
	DateEnd      time.Time `json:"date_end"`
}

// EventAward represents an award given to a team at an event. EventID, TeamID, and AwardID together form the primary key.
type EventAward struct {
	EventID string `json:"event_id"`
	TeamID  int    `json:"team_id"`
	AwardID int    `json:"award_id"`
}

// EventRanking represents a team's ranking in an event. EventID and TeamID together form the primary key.
type EventRanking struct {
	EventID        string  `json:"event_id"`
	TeamID         int     `json:"team_id"`
	Rank           int     `json:"rank"`
	SortOrder1     float64 `json:"sort_order1"`
	SortOrder2     float64 `json:"sort_order2"`
	SortOrder3     float64 `json:"sort_order3"`
	SortOrder4     float64 `json:"sort_order4"`
	SortOrder5     float64 `json:"sort_order5"`
	SortOrder6     float64 `json:"sort_order6"`
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	Ties           int     `json:"ties"`
	Dq             int     `json:"dq"`
	MatchesPlayed  int     `json:"matches_played"`
	MatchesCounted int     `json:"matches_counted"`
}

// EventAdvancement represents a team advancing from an event. EventID and TeamID together form the primary key.
type EventAdvancement struct {
	EventID string `json:"event_id"`
	TeamID  int    `json:"team_id"`
}

// GetEventID generates an EventID from the given EventCode and DateStart.
func GetEventID(eventCode string, dateStart time.Time) string {
	return fmt.Sprintf("%s : %04d-%02d-%02d", eventCode, dateStart.Year(), int(dateStart.Month()), dateStart.Day())
}

// GetEvent retrieves an event from the database by its ID.
func GetEvent(eventID string) *Event {
	var event Event
	stmt := GetStatement("getEvent")
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
func SaveEvent(event *Event) error {
	stmt := GetStatement("saveEvent")
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
func GetEventAwards(eventID string) []EventAward {
	stmt := GetStatement("getEventAwards")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var awards []EventAward
	for rows.Next() {
		var ea EventAward
		err := rows.Scan(&ea.EventID, &ea.TeamID, &ea.AwardID)
		if err != nil {
			continue
		}
		awards = append(awards, ea)
	}
	return awards
}

// SaveEventAward saves or updates an event award in the
func SaveEventAward(ea *EventAward) error {
	stmt := GetStatement("saveEventAward")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(ea.EventID, ea.TeamID, ea.AwardID)
	return err
}

// GetEventRankings retrieves all rankings for a specific event.
func GetEventRankings(eventID string) []EventRanking {
	stmt := GetStatement("getEventRankings")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var rankings []EventRanking
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
		rankings = append(rankings, er)
	}
	return rankings
}

// SaveEventRanking saves or updates an event ranking in the
func SaveEventRanking(er *EventRanking) error {
	stmt := GetStatement("saveEventRanking")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(er.EventID, er.TeamID, er.Rank, er.SortOrder1, er.SortOrder2, er.SortOrder3, er.SortOrder4, er.SortOrder5, er.SortOrder6, er.Wins, er.Losses, er.Ties, er.Dq, er.MatchesPlayed, er.MatchesCounted)
	return err
}

// GetEventAdvancements retrieves all team advancements for a specific event.
func GetEventAdvancements(eventID string) []EventAdvancement {
	stmt := GetStatement("getEventAdvancements")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query(eventID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var advancements []EventAdvancement
	for rows.Next() {
		var ea EventAdvancement
		err := rows.Scan(&ea.EventID, &ea.TeamID)
		if err != nil {
			continue
		}
		advancements = append(advancements, ea)
	}
	return advancements
}

// SaveEventAdvancement saves or updates an event advancement in the
func SaveEventAdvancement(ea *EventAdvancement) error {
	stmt := GetStatement("saveEventAdvancement")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(ea.EventID, ea.TeamID)
	return err
}

// GetRegionCodes retrieves all unique region codes from events, sorted alphabetically.
func GetRegionCodes() []string {
	stmt := GetStatement("getRegionCodes")
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
