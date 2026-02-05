package request

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
)

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
		slog.Error("Error requesting match results:", "error", err)
		return nil
	}
	slog.Debug("Requesting match results...", "count", len(ftcMatches))

	ftcScores, err := ftc.GetEventScores(strconv.Itoa(event.Year), event.EventCode, matchType)
	if err != nil {
		slog.Error("failed to get event scores", "error", err)
	}

	matches := make([]*database.Match, 0, len(ftcMatches))
	for _, ftcMatch := range ftcMatches {
		match := getMatch(event, ftcMatch)
		matches = append(matches, match)

		var ftcScore *ftc.MatchScores
		for _, score := range ftcScores {
			if score.MatchNumber == ftcMatch.MatchNumber {
				ftcScore = score
				break
			}
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
	slog.Info("Finished requesting match results", "count", len(matches))
	return matches
}

// getMatch creates a database.Match from an ftc.Match.
func getMatch(event *database.Event, ftcMatch *ftc.Match) *database.Match {
	match := &database.Match{
		EventID:         event.EventID,
		MatchID:         database.GetMatchID(event, ftcMatch),
		MatchNumber:     ftcMatch.MatchNumber,
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
	}
	blueScore = &database.MatchAllianceScore{
		MatchID:             match.MatchID,
		Alliance:            database.AllianceBlue,
		AutoPoints:          ftcMatch.ScoreBlueAuto,
		TeleopPoints:        ftcMatch.ScoreBlueFinal - ftcMatch.ScoreBlueAuto,
		FoulPointsCommitted: ftcMatch.ScoreBlueFoul,
	}

	if ftcScore != nil {
		for _, allianceScore := range ftcScore.Alliances {
			if allianceScore.Alliance == database.AllianceRed {
				redScore.AutoPoints = allianceScore.AutoPoints
				redScore.TeleopPoints = allianceScore.TeleopPoints
				redScore.FoulPointsCommitted = allianceScore.FoulPointsCommitted
				redScore.MinorFouls = allianceScore.MinorFouls
				redScore.MajorFouls = allianceScore.MajorFouls
				redScore.PreFoulTotal = allianceScore.PreFoulTotal
				redScore.TotalPoints = allianceScore.TotalPoints
			} else {
				blueScore.AutoPoints = allianceScore.AutoPoints
				blueScore.TeleopPoints = allianceScore.TeleopPoints
				blueScore.FoulPointsCommitted = allianceScore.FoulPointsCommitted
				blueScore.MinorFouls = allianceScore.MinorFouls
				blueScore.MajorFouls = allianceScore.MajorFouls
				blueScore.PreFoulTotal = allianceScore.PreFoulTotal
				blueScore.TotalPoints = allianceScore.TotalPoints
			}
		}
	}
	slog.Debug("Finished requesting match scores", "redScore", redScore, "blueScore", blueScore)
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
	slog.Debug("Finished requesting match teams", "redTeams", redTeams, "blueTeams", blueTeams)
	return redTeams, blueTeams
}
