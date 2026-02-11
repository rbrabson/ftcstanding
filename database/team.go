package database

import "fmt"

// Team represents a team that participates in competitions.
type Team struct {
	TeamID     int    `json:"team_id"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	City       string `json:"city"`
	StateProv  string `json:"state_prov"`
	Country    string `json:"country"`
	Website    string `json:"website"`
	RookieYear int    `json:"rookie_year"`
	HomeRegion string `json:"home_region"`
	RobotName  string `json:"robot_name"`
}

// TeamRanking represents the ranking information for a team based on their performance in matches at a specific event.
type TeamRanking struct {
	TeamID     int     `json:"team_id"`
	EventID    string  `json:"event_id"`
	NumMatches int     `json:"num_matches"`
	CCWM       float64 `json:"ccwm"`
	OPR        float64 `json:"opr"`
	NpOPR      float64 `json:"np_opr"`
	DPR        float64 `json:"dpr"`
	NpDPR      float64 `json:"np_dpr"`
	NpAvg      float64 `json:"np_avg"`
}

// String returns a string representation of the Team.
func (t *Team) String() string {
	return fmt.Sprintf("Team{ID: %d, Name: %q, City: %s, %s, Region: %s}",
		t.TeamID, t.Name, t.City, t.StateProv, t.HomeRegion)
}

// String returns a string representation of the TeamRanking.
func (tr *TeamRanking) String() string {
	return fmt.Sprintf("TeamRanking{TeamID: %d, EventID: %q, NumMatches: %d, CCWM: %.2f, OPR: %.2f, NpOPR: %.2f, DPR: %.2f, NpDPR: %.2f, NpAvg: %.2f}",
		tr.TeamID, tr.EventID, tr.NumMatches, tr.CCWM, tr.OPR, tr.NpOPR, tr.DPR, tr.NpDPR, tr.NpAvg)
}

// TeamFilter defines criteria for filtering teams.
type TeamFilter struct {
	TeamIDs     []int
	Countries   []string
	HomeRegions []string
	EventCodes  []string
}

// TeamRankingFilter defines criteria for filtering team rankings.
type TeamRankingFilter struct {
	TeamIDs  []int
	EventIDs []string
}
