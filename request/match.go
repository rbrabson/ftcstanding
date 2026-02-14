package request

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

// GetAndSaveMatches retrieves all matches for an event and saves them to the database.
func RequestAndSaveMatches(event *database.Event) []*database.Match {
	qualifierMatches := RequestAndSaveMatchesByType(event, ftc.QUALIFIER)
	playoffMatches := RequestAndSaveMatchesByType(event, ftc.PLAYOFF)
	return append(qualifierMatches, playoffMatches...)
}

// GetMatches retrieves all matches for an event.
func RequestMatches(event *database.Event) []*database.Match {
	qualifierMatches := RequestMatchesByType(event, ftc.QUALIFIER)
	playoffMatches := RequestMatchesByType(event, ftc.PLAYOFF)
	return append(qualifierMatches, playoffMatches...)
}

// GetAndSaveMatchesByType retrieves all qualification matches for an event and saves them to the database.
func RequestAndSaveMatchesByType(event *database.Event, matchType ftc.MatchType) []*database.Match {
	matches := RequestMatchesByType(event, matchType)
	for _, match := range matches {
		_ = db.SaveMatch(match)
	}
	return matches
}

// GetMatchesByType retrieves all qualification matches for an event.
func RequestMatchesByType(event *database.Event, matchType ftc.MatchType) []*database.Match {
	ftcMatches, err := ftc.GetMatchResults(strconv.Itoa(event.Year), event.EventCode, matchType)
	if err != nil {
		slog.Error("Error requesting match results:", "year", event.Year, "eventCode", event.EventCode, "matchType", matchType, "error", err)
		return nil
	}
	slog.Info("Retrieved match results...", "count", len(ftcMatches))

	ftcScores, err := ftc.GetEventScores(strconv.Itoa(event.Year), event.EventCode, matchType)
	if err != nil {
		slog.Error("failed to get event scores", "year", event.Year, "eventCode", event.EventCode, "matchType", matchType, "error", err)
		return nil
	}
	slog.Info("Retrieved event scores...", "count", len(ftcScores))

	matches := make([]*database.Match, 0, len(ftcMatches))
	for _, ftcMatch := range ftcMatches {
		match := getMatch(event, ftcMatch)
		matches = append(matches, match)

		// TODO: this is wrong
		var ftcScore *ftc.MatchScores
		for _, score := range ftcScores {
			var matchNumber int
			if strings.EqualFold(string(ftc.PLAYOFF), ftcMatch.TournamentLevel) {
				matchNumber = ftcMatch.Series
			} else {
				matchNumber = ftcMatch.MatchNumber
			}
			if score.MatchNumber == matchNumber {
				ftcScore = score
				break
			}
		}
		if ftcScore == nil {
			slog.Info("No match scores available", "year", event.Year, "eventCode", event.EventCode, "matchType", matchType)
		}

		redScore, blueScore := getMatchScores(match, ftcMatch, ftcScore)
		_ = db.SaveMatchAllianceScore(redScore)
		_ = db.SaveMatchAllianceScore(blueScore)

		redTeams, blueTeams := getMatchTeams(match, ftcMatch)
		for _, team := range redTeams {
			_ = db.SaveMatchTeam(team)
		}
		for _, team := range blueTeams {
			_ = db.SaveMatchTeam(team)
		}
	}
	slog.Info("Finished processing match results and event results", "count", len(matches))
	return matches
}

// getMatch creates a database.Match from an ftc.Match.
func getMatch(event *database.Event, ftcMatch *ftc.Match) *database.Match {
	tournamentLevel := titleCaser.String(ftcMatch.TournamentLevel)
	var matchNumber int
	if strings.EqualFold(tournamentLevel, string(ftc.PLAYOFF)) {
		matchNumber = ftcMatch.Series
	} else {
		matchNumber = ftcMatch.MatchNumber
	}

	match := &database.Match{
		EventID:         event.EventID,
		MatchID:         database.GetMatchID(event, ftcMatch.TournamentLevel, matchNumber),
		MatchType:       tournamentLevel,
		MatchNumber:     matchNumber,
		ActualStartTime: ftcMatch.ActualStartTime,
		Description:     ftcMatch.Description,
		TournamentLevel: ftcMatch.TournamentLevel,
	}

	return match
}

// getMatchScores creates database.MatchAllianceScore objects from an ftc.Match.
func getMatchScores(match *database.Match, ftcMatch *ftc.Match, ftcScore *ftc.MatchScores) (redScore, blueScore *database.MatchAllianceScore) {
	redScore = &database.MatchAllianceScore{
		MatchID:             match.MatchID,
		Alliance:            database.AllianceRed,
		AutoPoints:          ftcMatch.ScoreRedAuto,
		TeleopPoints:        ftcMatch.ScoreRedFinal - ftcMatch.ScoreRedAuto,
		FoulPointsCommitted: ftcMatch.ScoreRedFoul,
		TotalPoints:         ftcMatch.ScoreRedFinal,
	}
	blueScore = &database.MatchAllianceScore{
		MatchID:             match.MatchID,
		Alliance:            database.AllianceBlue,
		AutoPoints:          ftcMatch.ScoreBlueAuto,
		TeleopPoints:        ftcMatch.ScoreBlueFinal - ftcMatch.ScoreBlueAuto,
		FoulPointsCommitted: ftcMatch.ScoreBlueFoul,
		TotalPoints:         ftcMatch.ScoreBlueFinal,
	}

	if ftcScore != nil {
		for _, allianceScore := range ftcScore.Alliances {
			if strings.EqualFold(allianceScore.Alliance, database.AllianceRed) {
				redScore.AutoPoints = allianceScore.AutoPoints
				redScore.TeleopPoints = allianceScore.TeleopPoints
				redScore.FoulPointsCommitted = allianceScore.FoulPointsCommitted
				redScore.MinorFouls = allianceScore.MinorFouls
				redScore.MajorFouls = allianceScore.MajorFouls
				redScore.PreFoulTotal = allianceScore.PreFoulTotal
			} else {
				blueScore.AutoPoints = allianceScore.AutoPoints
				blueScore.TeleopPoints = allianceScore.TeleopPoints
				blueScore.FoulPointsCommitted = allianceScore.FoulPointsCommitted
				blueScore.MinorFouls = allianceScore.MinorFouls
				blueScore.MajorFouls = allianceScore.MajorFouls
				blueScore.PreFoulTotal = allianceScore.PreFoulTotal
			}
		}
	}
	slog.Debug("Finished processing match scores", "redScore", redScore, "blueScore", blueScore)
	return redScore, blueScore
}

// getMatchTeams creates database.MatchTeam objects from an ftc.Match.
func getMatchTeams(match *database.Match, ftcMatch *ftc.Match) (redTeams, blueTeams []*database.MatchTeam) {
	redTeams = make([]*database.MatchTeam, 0, len(ftcMatch.Teams)/2)
	blueTeams = make([]*database.MatchTeam, 0, len(ftcMatch.Teams)/2)
	for _, team := range ftcMatch.Teams {
		var alliance string
		if strings.HasPrefix(strings.ToLower(team.Station), strings.ToLower(database.AllianceRed)) {
			alliance = database.AllianceRed
		} else {
			alliance = database.AllianceBlue
		}
		matchTeam := &database.MatchTeam{
			MatchID:  match.MatchID,
			TeamID:   team.TeamNumber,
			Alliance: alliance,
			Dq:       team.DQ,
			OnField:  team.OnField,
		}
		if alliance == database.AllianceRed {
			redTeams = append(redTeams, matchTeam)
		} else {
			blueTeams = append(blueTeams, matchTeam)
		}
	}
	slog.Debug("Finished processing match teams", "redTeams", redTeams, "blueTeams", blueTeams)
	return redTeams, blueTeams
}

// StoreEventTeamsFromMatches extracts all unique teams from matches and stores them as EventTeam entries.
// This should be called after matches have been retrieved and saved to ensure the event_teams table is populated.
func StoreEventTeamsFromMatches(event *database.Event) error {
	// Get all matches for the event from the database
	matches, err := db.GetMatchesByEvent(event.EventID)
	if err != nil {
		slog.Error("failed to load matches for event", "eventID", event.EventID, "error", err)
		return err
	}
	if len(matches) == 0 {
		slog.Warn("no matches found for event", "eventID", event.EventID)
		return nil
	}

	// Collect all unique team IDs from matches
	teamIDsMap := make(map[int]bool)
	for _, match := range matches {
		matchTeams, err := db.GetMatchTeams(match.MatchID)
		if err != nil {
			slog.Error("failed to load match teams", "matchID", match.MatchID, "error", err)
			continue
		}
		for _, mt := range matchTeams {
			teamIDsMap[mt.TeamID] = true
		}
	}

	// Store EventTeam entries for all unique teams
	for teamID := range teamIDsMap {
		eventTeam := &database.EventTeam{
			EventID: event.EventID,
			TeamID:  teamID,
		}
		if err := db.SaveEventTeam(eventTeam); err != nil {
			slog.Error("failed to save event team", "eventID", event.EventID, "teamID", teamID, "error", err)
			return err
		}
	}

	slog.Info("stored event teams from matches", "eventID", event.EventID, "eventCode", event.EventCode, "teamCount", len(teamIDsMap))
	return nil
}
