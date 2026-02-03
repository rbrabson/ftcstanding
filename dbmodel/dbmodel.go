package dbmodel

// InitStatements initializes all prepared statements for the dbmodel package.
func InitStatements() error {
	if err := InitEventStatements(); err != nil {
		return err
	}
	if err := InitAwardStatements(); err != nil {
		return err
	}
	if err := InitMatchStatements(); err != nil {
		return err
	}
	if err := InitTeamStatements(); err != nil {
		return err
	}

	return nil
}
