package database

import "fmt"

// GetAward retrieves an award from a database by its ID.
func (db *sqldb) GetAward(awardID int) *Award {
	var award Award
	stmt := db.getStatement("getAward")
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
func (db *sqldb) GetAllAwards() []*Award {
	stmt := db.getStatement("getAllAwards")
	if stmt == nil {
		return nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil
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
	return awards
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
