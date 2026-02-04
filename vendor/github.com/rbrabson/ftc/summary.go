package ftc

import (
	"encoding/json"
	"fmt"
)

// Championships is the information about the season summary
type Championships struct {
	EventCount    int             `json:"eventCount"`
	GameName      string          `json:"gameName"`
	Kickoff       string          `json:"kickoff"`
	RookieStart   int             `json:"rookieStart"`
	TeamCount     int             `json:"teamCount"`
	Championships []*Championship `json:"fRCChampionships"`
}

// Championship is the information about a given regional championship.
type Championship struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	Location  string `json:"location"`
}

// GetSeasonSummary returns a list of the regional championships for the given season
func GetSeasonSummary(season string) (*Championships, error) {
	url := fmt.Sprintf("%s/%s", server, season)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Championships
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return &output, nil
}

func (c Championships) String() string {
	body, _ := json.Marshal(c)
	return string(body)
}

func (c Championship) String() string {
	body, _ := json.Marshal(c)
	return string(body)
}
