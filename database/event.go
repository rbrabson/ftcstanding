package database

import (
	"fmt"
	"time"

	"github.com/rbrabson/ftc"
)

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

// GetEventID generates a unique ID for an event based on its FTC code and start date.
func GetEventID(ftcEvent *ftc.Event, dateStart time.Time) string {
	return fmt.Sprintf("%s : %d", ftcEvent.Code, dateStart.Year())
}
