package query

import (
	"slices"
	"strings"

	"github.com/rbrabson/ftcstanding/database"
)

// TeamAward represents an award with full team details.
type TeamAward struct {
	Award *database.EventAward
	Team  *database.Team
}

// EventAwards represents an event with all team awards.
type EventAwards struct {
	Event  *database.Event
	Awards []*TeamAward
}

// AwardsByEventQuery retrieves all awards won by teams at a given event.
// It returns an EventAwards object containing the event and all awards with full team details.
func AwardsByEventQuery(eventCode string, year int) (*EventAwards, error) {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events, err := db.GetAllEvents(filter)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}

	// Find the event matching the year
	var event *database.Event
	for _, e := range events {
		if e.Year == year {
			event = e
			break
		}
	}

	if event == nil {
		return nil, nil
	}

	// Get all awards for the event
	eventAwards, err := db.GetEventAwards(event.EventID)
	if err != nil {
		return nil, err
	}
	if len(eventAwards) == 0 {
		return &EventAwards{
			Event:  event,
			Awards: []*TeamAward{},
		}, nil
	}

	// Retrieve the full team details for each award
	var teamAwards []*TeamAward
	for _, award := range eventAwards {
		team, err := db.GetTeam(award.TeamID)
		if err != nil {
			return nil, err
		}
		if team != nil {
			teamAwards = append(teamAwards, &TeamAward{
				Award: award,
				Team:  team,
			})
		}
	}

	// Sort by award ID, then series, then team ID
	slices.SortFunc(teamAwards, func(a, b *TeamAward) int {
		// Get sort priorities for each award
		priorityA := getAwardSortPriority(a.Award.Name)
		priorityB := getAwardSortPriority(b.Award.Name)

		// Sort by priority first
		if priorityA != priorityB {
			return priorityA - priorityB
		}

		// Within same type, sort by series
		if a.Award.Series != b.Award.Series {
			return a.Award.Series - b.Award.Series
		}

		// Finally, sort by team ID
		return a.Team.TeamID - b.Team.TeamID
	})

	return &EventAwards{
		Event:  event,
		Awards: teamAwards,
	}, nil
}

// getAwardSortPriority returns the sort priority for an award based on its name.
// Lower numbers come first.
func getAwardSortPriority(awardName string) int {
	switch {
	case strings.EqualFold(awardName, "inspire award"):
		return 1
	case strings.EqualFold(awardName, "inspire award 2nd place"):
		return 2
	case strings.EqualFold(awardName, "inspire award 3rd place"):
		return 3
	case strings.EqualFold(awardName, "winning alliance - captain"):
		return 4
	case strings.EqualFold(awardName, "winning alliance - 1st team selected"):
		return 5
	case strings.EqualFold(awardName, "finalist alliance - captain"):
		return 6
	case strings.EqualFold(awardName, "finalist alliance - 1st team selected"):
		return 7
	case strings.EqualFold(awardName, "think award"):
		return 8
	case strings.EqualFold(awardName, "think award 2nd place"):
		return 9
	case strings.EqualFold(awardName, "connect award"):
		return 10
	case strings.EqualFold(awardName, "connect award 2nd place"):
		return 11
	case strings.EqualFold(awardName, "sustain award"):
		return 12
	case strings.EqualFold(awardName, "sustain award 2nd place"):
		return 13
	case strings.EqualFold(awardName, "innovate award sponsored by rtx"):
		return 14
	case strings.EqualFold(awardName, "innovate award sponsored by rtx 2nd place"):
		return 15
	case strings.EqualFold(awardName, "control award"):
		return 16
	case strings.EqualFold(awardName, "control award 2nd place"):
		return 17
	case strings.EqualFold(awardName, "reach award"):
		return 18
	case strings.EqualFold(awardName, "reach award 2nd place"):
		return 19
	case strings.EqualFold(awardName, "design award"):
		return 20
	case strings.EqualFold(awardName, "design award 2nd place"):
		return 21
	case strings.EqualFold(awardName, "judges' choice award"):
		return 22
	case strings.EqualFold(awardName, "dean's list winners"):
		return 23
	default:
		return 999 // Unknown awards go to the end
	}
}
