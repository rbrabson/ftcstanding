package ftc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Events is information about FTC events.
type Events struct {
	Events     []*Event `json:"events,omitempty"`
	EventCount int      `json:"eventCount,omitempty"`
}

// Event is information about a given FTC event.
type Event struct {
	EventID       string  `json:"eventId,omitempty"`
	Code          string  `json:"code,omitempty"`
	DivisionCode  *string `json:"divisionCode,omitempty"`
	Name          string  `json:"name,omitempty"`
	Remote        bool    `json:"remote,omitempty"`
	Hybrid        bool    `json:"hybrid,omitempty"`
	FieldCount    int     `json:"fieldCount,omitempty"`
	Published     bool    `json:"published,omitempty"`
	Type          string  `json:"type,omitempty"`
	TypeName      string  `json:"typeName,omitempty"`
	RegionCode    string  `json:"regionCode,omitempty"`
	LeagueCode    *string `json:"leagueCode,omitempty"`
	DistrictCode  string  `json:"districtCode,omitempty"`
	Venue         string  `json:"venue,omitempty"`
	Address       string  `json:"address,omitempty"`
	City          string  `json:"city,omitempty"`
	Stateprov     string  `json:"stateprov,omitempty"`
	Country       string  `json:"country,omitempty"`
	Website       string  `json:"website,omitempty"`
	LiveStreamURL string  `json:"liveStreamUrl,omitempty"`
	Coordinates   *string `json:"coordinates,omitempty"`
	Webcasts      *string `json:"webcasts,omitempty"`
	Timezone      string  `json:"timezone,omitempty"`
	DateStart     Time    `json:"dateStart,omitempty"`
	DateEnd       Time    `json:"dateEnd,omitempty"`
}

// GetEvents returns the list of events for a given season. Supported qparms are `eventCode` and `teamNumber`,
// which are mutually exclusive.
func GetEvents(season string, qparms ...map[string]string) ([]*Event, error) {
	sb := strings.Builder{}
	sb.WriteString(server)
	sb.WriteString("/")
	sb.WriteString(season)
	sb.WriteString("/events")
	if len(qparms) > 0 {
		for k, v := range qparms[0] {
			sb.WriteString("?")
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
			break
		}
	}
	url := sb.String()

	body, err := getURL(url)
	if err != nil {
		fmt.Println(url)
		return nil, err
	}

	var output Events
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	// Return the output
	return output.Events, nil
}

// String returns a string representation of an Event. In this case, it is a json string.
func (e Event) String() string {
	body, _ := json.Marshal(e)
	return string(body)
}
