package database

// Award is an award that is given in a given season
type Award struct {
	AwardID     int    `json:"award_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ForPerson   bool   `json:"for_person"`
}
