package ftc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Leagues is the data for the FTC leagues
type Leagues struct {
	Leagues     []*League `json:"leagues"`
	LeagueCount int       `json:"leagueCount"`
}

// League is the data for a given FTC league
type League struct {
	Region           string  `json:"region,omitempty"`
	Code             string  `json:"code,omitempty"`
	Name             string  `json:"name,omitempty"`
	Remote           bool    `json:"remote,omitempty"`
	ParentLeagueCode *string `json:"parentLeagueCode,omitempty"`
	ParentLeagueName *string `json:"parentLeagueName,omitempty"`
	Location         string  `json:"location,omitempty"`
}

type LeagueMembers struct {
	Members []int `json:"members"`
}

// GetLeagues returns the list of rankings for FTC leagues. Supported qparms are `regionCode` and `leagueCode`.
func GetLeagues(season string, qparms ...map[string]string) ([]*League, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/leagues")
	if len(qparms) > 0 {
		firstQparm := true
		for k, v := range qparms[0] {
			if firstQparm {
				sb.WriteString("?")
				firstQparm = false
			} else {
				sb.WriteString("&")
			}
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
		}
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Leagues
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Leagues, nil
}

// GetLeagueMembers returns the list of members in the league
func GetLeagueMembers(season, regionCode, leagueCode string) ([]int, error) {
	url := fmt.Sprintf("%s/%s/leagues/members/%s/%s", server, season, regionCode, leagueCode)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output LeagueMembers
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Members, nil
}

// String returns a string representation of League. In this case, it is a json string.
func (l League) String() string {
	body, _ := json.Marshal(l)
	return string(body)
}
