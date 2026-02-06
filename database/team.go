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

// String returns a string representation of the Team.
func (t *Team) String() string {
	return fmt.Sprintf("Team{ID: %d, Name: %q, City: %s, %s, Region: %s}",
		t.TeamID, t.Name, t.City, t.StateProv, t.HomeRegion)
}

// TeamFilter defines criteria for filtering teams.
type TeamFilter struct {
	TeamIDs     []int
	Countries   []string
	HomeRegions []string
}
