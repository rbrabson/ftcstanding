package ftc

import (
	"encoding/json"
	"strings"
)

// Type of FTC match
type MatchType string

// The types of FTC matches
const (
	QUALIFIER MatchType = "qual"
	PLAYOFF   MatchType = "playoff"
)

// Matches is the list of matches
type Matches struct {
	Matches []*Match `json:"matches"`
}

// Match is a match that takes place at a given event
type Match struct {
	ActualStartTime string       `json:"actualStartTime"`
	Description     string       `json:"description"`
	TournamentLevel string       `json:"tournamentLevel"`
	Series          int          `json:"series"`
	MatchNumber     int          `json:"matchNumber"`
	ScoreRedFinal   int          `json:"scoreRedFinal"`
	ScoreRedFoul    int          `json:"scoreRedFoul"`
	ScoreRedAuto    int          `json:"scoreRedAuto"`
	ScoreBlueFinal  int          `json:"scoreBlueFinal"`
	ScoreBlueFoul   int          `json:"scoreBlueFoul"`
	ScoreBlueAuto   int          `json:"scoreBlueAuto"`
	PostResultTime  string       `json:"postResultTime"`
	Teams           []*MatchTeam `json:"teams"`
	ModifiedOn      string       `json:"modifiedOn"`
}
type MatchTeam struct {
	TeamNumber int    `json:"teamNumber"`
	Station    string `json:"station"`
	DQ         bool   `json:"dq"`
	OnField    bool   `json:"onField"`
}

// GetMatchResults returns the results of a given event.
func GetMatchResults(season, eventCode string, tournamentLevel MatchType, teamNumber ...string) ([]*Match, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/matches/")
	sb.WriteString(eventCode)
	sb.WriteString("?tournamentLevel=")
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

	var output Matches
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Matches, nil
}

func (m Match) String() string {
	body, _ := json.Marshal(m)
	return string(body)
}

func (m MatchTeam) String() string {
	body, _ := json.Marshal(m)
	return string(body)
}
