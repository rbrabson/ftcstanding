package ftc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Awards is the list of awards for a given season
type Awards struct {
	Awards []*Award `json:"awards"`
}

// Award is an award that is given in a given season
type Award struct {
	AwardID     int    `json:"awardId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ForPerson   bool   `json:"forPerson"`
}

// TeamAwards is the list of awards received by a team
type TeamAwards struct {
	Awards []*TeamAward `json:"awards"`
}

// TeamAward is an award that is received by a given team
type TeamAward struct {
	AwardID      int     `json:"awardId"`
	EventCode    string  `json:"eventCode"`
	Name         string  `json:"name"`
	Series       int     `json:"series"`
	TeamNumber   int     `json:"teamNumber"`
	SchoolName   *string `json:"schoolName,omitempty"`
	FullTeamName string  `json:"fullTeamName"`
	Person       *string `json:"person,omitempty"`
}

// GetAwardListing returns the list of awards for a given season
func GetAwardListing(season string) ([]*Award, error) {
	url := fmt.Sprintf("%s/%s/awards/list", server, season)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Awards
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Awards, nil
}

// GetEventAwards gets the list of awards given at an event
func GetEventAwards(season, eventCode string, teamNumber ...string) ([]*TeamAward, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/awards/")
	sb.WriteString(eventCode)
	if len(teamNumber) > 0 {
		sb.WriteString("?teamNumber")
		sb.WriteString(teamNumber[0])
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output TeamAwards
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Awards, nil
}

// GetTeamAwards gets the list of awards for a given team
func GetTeamAwards(season, teamNumber string, eventCode ...string) ([]*TeamAward, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString("/awards/")
	sb.WriteString(teamNumber)
	if teamNumber != "" {
		sb.WriteString("?eventCode")
		sb.WriteString(eventCode[0])
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output []*TeamAward
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output, nil
}

func (a Award) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}

func (a TeamAward) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}
