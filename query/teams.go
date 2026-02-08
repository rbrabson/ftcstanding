package query

import (
	"github.com/rbrabson/ftcstanding/database"
)

// TeamsQuery returns a list of teams that match the given filter.
func TeamsQuery(filter ...database.TeamFilter) []*database.Team {
	return db.GetAllTeams(filter...)
}
