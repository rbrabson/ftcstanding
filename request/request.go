package request

// RequestAndSaveAll requests and saves all data for a given season.
func RequestAndSaveAll(season string) {
	RequestAndSaveAwards(season)
	RequestAndSaveTeams(season)
	events := RequestAndSaveEvents(season)
	for _, event := range events {
		RequestAndSaveEventAwards(event)
		RequestAndSaveEventRankings(event)
		RequestAndSaveEventAdvancements(event)
		RequestAndSaveMatches(event)
	}
}
