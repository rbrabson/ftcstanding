package ftc

import (
	"encoding/json"
	"fmt"
)

// AllianceSelections is the list of alliance selections.
type AllianceSelections struct {
	Selections []*AllianceSelection `json:"selections"`
	Count      int                  `json:"count"`
}

// AllianceSelection is the alliance results for a given team
type AllianceSelection struct {
	Index  int    `json:"index"`
	Team   int    `json:"team"`
	Result string `json:"result"`
}

// Alliances is the list of alliances in a given tournament
type Alliances struct {
	Alliances []*Alliance `json:"alliances"`
	Count     int         `json:"count"`
}

// Alliance is the results for one alliance in a match between two alliances
type Alliance struct {
	Number         int     `json:"number"`
	Name           string  `json:"name"`
	Captain        int     `json:"captain"`
	CaptainDisplay string  `json:"captainDisplay"`
	Round1         int     `json:"round1,omitempty"`
	Round1Display  string  `json:"round1Display,omitempty"`
	Round2         int     `json:"round2,omitempty"`
	Round2Display  string  `json:"round2Display,omitempty"`
	Round3         *int    `json:"round3,omitempty"`
	Backup         *string `json:"backup,omitempty"`
	BackupReplaced *string `json:"backupReplaced,omitempty"`
}

// GetEventAlliances returns the alliance selectsions for the playoffs for the given event.
func GetEventAlliances(season, eventCode string) ([]*Alliance, error) {
	url := fmt.Sprintf("%s/%s/alliances/%s", server, season, eventCode)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output Alliances
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Alliances, nil
}

// GetAllianceSelections returns the teams that were selected into alliances for the given event.
func GetAllianceSelections(season, eventCode string) ([]*AllianceSelection, error) {
	url := fmt.Sprintf("%s/%s/alliances/%s/selection", server, season, eventCode)

	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var output AllianceSelections
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Selections, nil
}

func (a AllianceSelection) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}

func (a Alliances) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}

func (a Alliance) String() string {
	body, _ := json.Marshal(a)
	return string(body)
}
