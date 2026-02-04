package ftc

import (
	"encoding/json"
	"strings"
)

// Scores is the list of match scores at a given event.
type Scores struct {
	MatchScores []*MatchScores `json:"matchScores"`
}

// MatchScores is the results of a match at a given event.
type MatchScores struct {
	MatchLevel    string           `json:"matchLevel"`
	MatchSeries   int              `json:"matchSeries"`
	MatchNumber   int              `json:"matchNumber"`
	Randomization int              `json:"randomization"`
	Alliances     []*MatchAlliance `json:"alliances"`
}

// MatchAlliance is the detailed results for a given team in a match at a given event.
// It only contains data that is reported for all seasons, and ommits season-specific data.
type MatchAlliance struct {
	Alliance            string `json:"alliance"`
	Team                int    `json:"team"`
	Robot1Auto          bool   `json:"robot1Auto"`
	Robot2Auto          bool   `json:"robot2Auto"`
	Robot1Teleop        string `json:"robot1Teleop"`
	Robot2Teleop        string `json:"robot2Teleop"`
	AutoPoints          int    `json:"autoPoints"`
	TeleopPoints        int    `json:"teleopPoints"`
	FoulPointsCommitted int    `json:"foulPointsCommitted"`
	PreFoulTotal        int    `json:"preFoulTotal"`
	TotalPoints         int    `json:"totalPoints"`
	MajorFouls          int    `json:"majorFouls"`
	MinorFouls          int    `json:"minorFouls"`
}

// GetEventScores returns the results for a given event
func GetEventScores(season, eventCode string, tournamentLevel MatchType, teamNumber ...string) ([]*MatchScores, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/scores/")
	sb.WriteString(eventCode)
	sb.WriteString("/")
	sb.WriteString(string(tournamentLevel))
	if len(teamNumber) > 0 {
		sb.WriteString("?")
		sb.WriteString(teamNumber[0])
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Scores
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.MatchScores, nil
}

func (ms MatchScores) String() string {
	body, _ := json.Marshal(ms)
	return string(body)
}

func (ma MatchAlliance) String() string {
	body, _ := json.Marshal(ma)
	return string(body)
}
