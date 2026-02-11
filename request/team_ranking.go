package request

import (
	"log/slog"
	"maps"
	"slices"
	"sort"

	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/performance"
)

// RequestAndSaveTeamRankings calculates and saves team performance rankings for an event.
// It retrieves match data from the database, calculates performance metrics (OPR, NpOPR, CCWM, DPR, NpDPR, NpAvg),
// and stores the results as TeamRanking records in the database.
func RequestAndSaveTeamRankings(event *database.Event) error {
	// Get all matches for this event from the database
	dbMatches := db.GetMatchesByEvent(event.EventID)
	if len(dbMatches) == 0 {
		slog.Info("No matches found for event", "event", event.EventCode)
		return nil
	}

	var matches []performance.Match
	teamSet := make(map[int]any)

	// Convert database matches to performance.Match format
	for _, dbMatch := range dbMatches {
		// Get alliance scores
		redScore := db.GetMatchAllianceScore(dbMatch.MatchID, database.AllianceRed)
		blueScore := db.GetMatchAllianceScore(dbMatch.MatchID, database.AllianceBlue)

		if redScore == nil || blueScore == nil {
			continue
		}

		// Get teams in the match
		matchTeams := db.GetMatchTeams(dbMatch.MatchID)

		var redTeams []int
		var blueTeams []int

		for _, mt := range matchTeams {
			if !mt.OnField || mt.Dq {
				continue
			}

			if mt.Alliance == database.AllianceRed {
				redTeams = append(redTeams, mt.TeamID)
			} else {
				blueTeams = append(blueTeams, mt.TeamID)
			}

			teamSet[mt.TeamID] = struct{}{}
		}

		// Only include matches with teams on both alliances
		if len(redTeams) == 0 || len(blueTeams) == 0 {
			continue
		}

		matches = append(matches, performance.Match{
			RedTeams:      redTeams,
			BlueTeams:     blueTeams,
			RedScore:      float64(redScore.TotalPoints),
			BlueScore:     float64(blueScore.TotalPoints),
			RedPenalties:  float64(redScore.FoulPointsCommitted),
			BluePenalties: float64(blueScore.FoulPointsCommitted),
		})
	}

	// Skip if no valid matches
	if len(matches) == 0 {
		slog.Info("No valid matches found for event", "event", event.EventCode)
		return nil
	}

	// Convert teamSet to sorted slice
	eventTeams := slices.Collect(maps.Keys(teamSet))
	sort.Ints(eventTeams)

	// Calculate lambda for this event
	lambdaValue := getLambda(matches)

	slog.Info("calculating team rankings", "event", event.EventCode, "matches", len(matches), "teams", len(eventTeams), "lambda", lambdaValue)

	// Calculate performance metrics for this event
	calculator := performance.Calculator{
		Matches: matches,
		Teams:   eventTeams,
		Lambda:  lambdaValue,
	}

	opr := calculator.CalculateOPR()
	npopr := calculator.CalculateNpOPR()
	ccwm := calculator.CalculateCCWM()
	dpr := calculator.CalculateDPR()
	npdpr := calculator.CalculateNpDPR()

	// Save TeamRanking records for each team
	for _, teamID := range eventTeams {
		// Count matches for this team in this event
		matchCount := 0
		for _, m := range matches {
			if slices.Contains(m.RedTeams, teamID) || slices.Contains(m.BlueTeams, teamID) {
				matchCount++
			}
		}

		npavg := calculator.CalculateNpAVG(matches, teamID)

		teamRanking := &database.TeamRanking{
			TeamID:     teamID,
			EventID:    event.EventID,
			NumMatches: matchCount,
			CCWM:       ccwm[teamID],
			OPR:        opr[teamID],
			NpOPR:      npopr[teamID],
			DPR:        dpr[teamID],
			NpDPR:      npdpr[teamID],
			NpAvg:      npavg,
		}

		if err := db.SaveTeamRanking(teamRanking); err != nil {
			slog.Error("Failed to save team ranking", "event", event.EventCode, "team", teamID, "error", err)
			continue
		}
	}

	slog.Info("Finished calculating team rankings", "event", event.EventCode, "teamsProcessed", len(eventTeams))
	return nil
}
