package dbmodel

import (
	"fmt"
)

// Award is an award that is given in a given season
type Award struct {
	AwardID     int    `json:"award_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ForPerson   bool   `json:"for_person"`
}

// InitAwardStatements prepares all SQL statements for award operations.
func InitAwardStatements() error {
	queries := map[string]string{
		"getAward":     "SELECT award_id, name, description, for_person FROM awards WHERE award_id = ?",
		"getAllAwards": "SELECT award_id, name, description, for_person FROM awards",
		"saveAward":    "INSERT INTO awards (award_id, name, description, for_person) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description), for_person = VALUES(for_person)",
	}

	for name, query := range queries {
		if err := PrepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// GetAward retrieves an award from a database by its ID.
func GetAward(awardID int) *Award {
	var award Award
	stmt := GetStatement("getAward")
	if stmt == nil {
		return nil
	}
	err := stmt.QueryRow(awardID).Scan(
		&award.AwardID,
		&award.Name,
		&award.Description,
		&award.ForPerson,
	)
	if err != nil {
		return nil
	}
	return &award
}

// GetAllAwards retrieves all awards from the
func GetAllAwards() []Award {
	stmt := GetStatement("getAllAwards")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()

	var awards []Award
	for rows.Next() {
		var award Award
		err := rows.Scan(
			&award.AwardID,
			&award.Name,
			&award.Description,
			&award.ForPerson,
		)
		if err != nil {
			continue
		}
		awards = append(awards, award)
	}
	return awards
}

// SaveAward saves or updates an award in the
func SaveAward(award *Award) error {
	stmt := GetStatement("saveAward")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(award.AwardID, award.Name, award.Description, award.ForPerson)
	return err
}
