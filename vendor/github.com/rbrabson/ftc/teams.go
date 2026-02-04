package ftc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Teams returns a list of FTC teams. The information is returned in `pages`, so multiple requests
// may be required to get all FTC teams.
type Teams struct {
	Teams          []*Team `json:"teams"`
	TeamCountTotal int     `json:"teamCountTotal"`
	TeamCountPage  int     `json:"teamCountPage"`
	PageCurrent    int     `json:"pageCurrent"`
	PageTotal      int     `json:"pageTotal"`
}

// Team is information for a given FTC team.
type Team struct {
	TeamNumber        int     `json:"teamNumber,omitempty"`
	DisplayTeamNumber string  `json:"displayTeamNumber,omitempty"`
	NameFull          string  `json:"nameFull,omitempty"`
	NameShort         string  `json:"nameShort,omitempty"`
	SchoolName        *string `json:"schoolName,omitempty"`
	City              string  `json:"city,omitempty"`
	StateProv         string  `json:"stateProv,omitempty"`
	Country           string  `json:"country,omitempty"`
	Website           *string `json:"website,omitempty"`
	RookieYear        int     `json:"rookieYear,omitempty"`
	RobotName         *string `json:"robotName,omitempty"`
	DistrictCode      *string `json:"districtCode,omitempty"`
	HomeCMP           *string `json:"homeCMP,omitempty"`
	HomeRegion        *string `json:"homeRegion,omitempty"`
}

// GetTeams returns a `page` of FTC teams.
func GetTeams(season string, teamNumber ...string) ([]*Team, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/teams")
	if len(teamNumber) > 0 {
		sb.WriteString("?")
		sb.WriteString(teamNumber[0])
	}
	url := sb.String()

	// Get the first page of teams
	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Teams
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Make the slice large enough to contain all teams
	teams := make([]*Team, 0, output.TeamCountTotal)
	teams = append(teams, output.Teams...)

	// Loop through all remaining pages, appending the teams to the list
	numPages := output.PageTotal
	for i := 2; i <= numPages; i++ {
		pageURL := fmt.Sprintf("%s?page=%d", url, i)
		body, err := getURL(pageURL)
		if err != nil {
			return nil, err
		}

		var output Teams
		err = json.Unmarshal(body, &output)
		if err != nil {
			return nil, err
		}

		teams = append(teams, output.Teams...)
	}

	return teams, nil
}

// String returns a string representation of Teams. In this case, it is a json string.
func (teams *Teams) String() string {
	body, _ := json.Marshal(teams)
	return string(body)
}
