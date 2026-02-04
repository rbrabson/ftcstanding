package ftc

import (
	"encoding/json"
	"fmt"
)

// Rankings is the list of rankings for a given event.
type Rankings struct {
	Rankings []*Ranking `json:"rankings"`
}

// Ranking is the ranking for a given event.
type Ranking struct {
	Rank              int     `json:"rank"`
	TeamNumber        int     `json:"teamNumber"`
	DisplayTeamNumber string  `json:"displayTeamNumber"`
	TeamName          *string `json:"teamName,omitempty"`
	SortOrder1        float64 `json:"sortOrder1"`
	SortOrder2        float64 `json:"sortOrder2"`
	SortOrder3        float64 `json:"sortOrder3"`
	SortOrder4        float64 `json:"sortOrder4"`
	SortOrder5        float64 `json:"sortOrder5"`
	SortOrder6        float64 `json:"sortOrder6"`
	Wins              int     `json:"wins"`
	Losses            int     `json:"losses"`
	Ties              int     `json:"ties"`
	QualAverage       int     `json:"qualAverage"`
	DQ                int     `json:"dq"`
	MatchesPlayed     int     `json:"matchesPlayed"`
	MatchesCounted    int     `json:"matchesCounted"`
}

// GetRankings returns the team rankings in a given league.
func GetRankings(season, eventCode string) ([]*Ranking, error) {
	url := fmt.Sprintf("%s/%s/rankings/%s", server, season, eventCode)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Rankings
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Rankings, nil
}

// String returns a string representation of LeagueRanking. In this case, it is a json string.
func (r Ranking) String() string {
	body, _ := json.Marshal(r)
	return string(body)
}
