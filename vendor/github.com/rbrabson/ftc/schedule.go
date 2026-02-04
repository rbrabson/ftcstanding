package ftc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EventSchedules is the list of event schedules
type EventSchedules struct {
	Schedule []*EventSchedule `json:"schedule"`
}

// HybridSchedules is the list of results for a schedule at a given event
type HybridSchedules struct {
	Schedule []*HybridSchedule `json:"schedule"`
}

// EventSchedule is the schedule at a given event
type EventSchedule struct {
	Description     string          `json:"description,omitempty"`
	Field           string          `json:"field,omitempty"`
	TournamentLevel string          `json:"tournamentLevel,omitempty"`
	StartTime       string          `json:"startTime,omitempty"`
	Series          int             `json:"series,omitempty"`
	MatchNumber     int             `json:"matchNumber,omitempty"`
	Teams           []ScheduledTeam `json:"teams"`
	ModifiedOn      string          `json:"modifiedOn,omitempty"`
}

// HybridSchedule is the result for a scheduled match at a given event
type HybridSchedule struct {
	Description              string          `json:"description"`
	TournamentLevel          string          `json:"tournamentLevel"`
	Series                   int             `json:"series"`
	MatchNumber              int             `json:"matchNumber"`
	StartTime                string          `json:"startTime"`
	ActualStartTime          string          `json:"actualStartTime"`
	PostResultTime           string          `json:"postResultTime"`
	ScoreRedFinal            int             `json:"scoreRedFinal"`
	ScoreRedFoul             int             `json:"scoreRedFoul"`
	ScoreRedAuto             int             `json:"scoreRedAuto"`
	ScoreBlueFinal           int             `json:"scoreBlueFinal"`
	ScoreBlueFoul            int             `json:"scoreBlueFoul"`
	ScoreBlueAuto            int             `json:"scoreBlueAuto"`
	ScoreBlueDriveControlled *int            `json:"scoreBlueDriveControlled,omitempty"`
	ScoreBlueEndgame         *int            `json:"scoreBlueEndgame,omitempty"`
	RedWins                  bool            `json:"redWins"`
	BlueWins                 bool            `json:"blueWins"`
	Teams                    []ScheduledTeam `json:"teams"`
}

// ScheduledTeam is the team that is scheduled at a given tournament
type ScheduledTeam struct {
	TeamNumber        int    `json:"teamNumber,omitempty"`
	DisplayTeamNumber string `json:"displayTeamNumber,omitempty"`
	Station           string `json:"station,omitempty"`
	Team              string `json:"team,omitempty"`
	TeamName          string `json:"teamName,omitempty"`
	Surrogate         bool   `json:"surrogate,omitempty"`
	NoShow            bool   `json:"noShow,omitempty"`
	DQ                *bool  `json:"dq,omitempty"`
	OnField           *bool  `json:"onField,omitempty"`
}

// GetEventSchedule gets the match schedule for a given event.
func GetEventSchedule(season, eventCode string, tournamentLevel MatchType, teamNumber ...string) ([]*EventSchedule, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/schedule")
	sb.WriteString("/")
	sb.WriteString(eventCode)
	sb.WriteString("?")
	sb.WriteString("tournamentLevel=")
	sb.WriteString(string(tournamentLevel))
	if len(teamNumber) > 0 {
		sb.WriteString("&teamNumber=")
		sb.WriteString(teamNumber[0])
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output EventSchedules
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Schedule, nil
}

// GetHybridSchedule gets the hybrid schedule information for a given event.
func GetHybridSchedule(season, eventCode string, tournamentLevel MatchType) ([]*HybridSchedule, error) {
	url := fmt.Sprintf("%s/%s/schedule/%s?tournamentLevel=%s", server, season, eventCode, string(tournamentLevel))

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output HybridSchedules
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Schedule, nil
}

func (s EventSchedule) String() string {

	body, _ := json.Marshal(s)
	return string(body)
}

func (s HybridSchedule) String() string {
	body, _ := json.Marshal(s)
	return string(body)
}
