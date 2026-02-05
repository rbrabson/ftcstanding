package database

import "fmt"

// Award is an award that is given in a given season
type Award struct {
	AwardID     int    `json:"award_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ForPerson   bool   `json:"for_person"`
}

// String returns a string representation of the Award.
func (a *Award) String() string {
	personType := "Team"
	if a.ForPerson {
		personType = "Person"
	}
	return fmt.Sprintf("Award{ID: %d, Name: %q, Type: %s}", a.AwardID, a.Name, personType)
}
