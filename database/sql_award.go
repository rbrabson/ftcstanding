package database

import "fmt"

// InitAwardStatements prepares all SQL statements for award operations.
func (db *sqldb) initAwardStatements() error {
	queries := map[string]string{
		"getAward":     "SELECT award_id, name, description, for_person FROM awards WHERE award_id = ?",
		"getAllAwards": "SELECT award_id, name, description, for_person FROM awards",
		"saveAward":    "INSERT INTO awards (award_id, name, description, for_person) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description), for_person = VALUES(for_person)",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// GetAward retrieves an award from a database by its ID.
func (db *sqldb) GetAward(awardID int) (*Award, error) {
	var award Award
	stmt := db.getStatement("getAward")
	if stmt == nil {
		return nil, fmt.Errorf("prepared statement not found")
	}
	err := stmt.QueryRow(awardID).Scan(
		&award.AwardID,
		&award.Name,
		&award.Description,
		&award.ForPerson,
	)
	if err != nil {
		return nil, nil
	}
	return &award, nil
}

// GetAllAwards retrieves all awards from the
func (db *sqldb) GetAllAwards() ([]*Award, error) {
	stmt := db.getStatement("getAllAwards")
	if stmt == nil {
		return nil, fmt.Errorf("prepared statement not found")
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []*Award
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
		awards = append(awards, &award)
	}
	return awards, nil
}

// SaveAward saves or updates an award in the
func (db *sqldb) SaveAward(award *Award) error {
	stmt := db.getStatement("saveAward")
	if stmt == nil {
		return fmt.Errorf("prepared statement not found")
	}
	_, err := stmt.Exec(award.AwardID, award.Name, award.Description, award.ForPerson)
	return err
}
