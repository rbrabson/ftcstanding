package request

// Add code to request and build the database models and save them in the database.
// This should use the ftc package to do all of the processing.

import (
	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
)

// RequestAndSaveAwards requests awards from the FTC API for a given season and saves them in the database.
func RequestAndSaveAwards(season string) []*database.Award {
	awards := RequestAwards(season)
	for _, award := range awards {
		db.SaveAward(award)
	}
	return awards
}

// RequestAwards requests awards from the FTC API for a given season.
func RequestAwards(season string) []*database.Award {
	ftcAwards, err := ftc.GetAwardListing(season)
	if err != nil {
		return nil
	}
	awards := make([]*database.Award, 0, len(ftcAwards))
	for _, ftcAward := range ftcAwards {
		award := database.Award{
			AwardID:     ftcAward.AwardID,
			Name:        ftcAward.Name,
			Description: ftcAward.Description,
			ForPerson:   ftcAward.ForPerson,
		}
		awards = append(awards, &award)
	}
	return awards
}
