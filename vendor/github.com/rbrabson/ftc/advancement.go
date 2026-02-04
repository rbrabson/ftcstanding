package ftc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AdvancementsTo is the list of teams that advanced from or to a given tournament
type AdvancementsTo struct {
	AdvancesTo  string         `json:"advancesTo"`
	Slots       int            `json:"slots"`
	Advancement []*Advancement `json:"advancement"`
}

// AdvancementsFrom is the list of teams advancing to a future tournament
type AdvancementsFrom struct {
	AdvancedFrom       string         `json:"advancedFrom"`
	AdvancedFromRegion *string        `json:"advancedFromRegion"`
	Slots              int            `json:"slots"`
	Advancement        []*Advancement `json:"advancement"`
}

// Advancement is the advancement information for a given team
type Advancement struct {
	Team        int    `json:"team"`
	DisplayTeam string `json:"displayTeam"`
	Slot        int    `json:"slot"`
	Criteria    string `json:"criteria"`
	Status      string `json:"status"`
}

// GetAdvancementFrom returns the source events from which teams advanced from to reach
// the specified event.
func GetAdvancementsFrom(season, eventCode string) ([]*AdvancementsFrom, error) {
	url := fmt.Sprintf("%s/%s/advancement/%s/source", server, season, eventCode)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output []*AdvancementsFrom
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output, nil
}

// GetAdvancementTo returns the list of teams advancing from the event and to which event
// they are advancing.
func GetAdvancementsTo(season, eventCode string, excludeSkipped ...bool) (*AdvancementsTo, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/advancement/")
	sb.WriteString(eventCode)
	if len(excludeSkipped) > 0 {
		sb.WriteString("?excludedSkipped=")
		if excludeSkipped[0] {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output AdvancementsTo
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return &output, nil
}

func (a AdvancementsTo) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}

func (a AdvancementsFrom) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}

func (a Advancement) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}
