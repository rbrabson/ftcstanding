package query

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rbrabson/ftc"
	"github.com/rbrabson/ftcstanding/database"
)

// TeamAdvancement represents a team's advancement information from an event.
type TeamAdvancement struct {
	Rank                int
	Team                *database.Team
	Ranking             *database.EventRanking
	TotalPoints         int
	JudgingPoints       int
	PlayoffPoints       int
	SelectionPoints     int
	QualificationPoints int
	AdvancementNumber   string // Rank by total points for advancing teams, or "-"
	Advances            bool
	Status              string // Status from EventAdvancement (e.g., "already advanced")
}

// AdvancementReport represents an event with all team advancement information.
type AdvancementReport struct {
	Event            *database.Event
	TeamAdvancements []*TeamAdvancement
}

// AdvancementReportQuery retrieves advancement information for all teams at an event.
// It returns an AdvancementReport with teams sorted by their ranking.
func AdvancementReportQuery(eventCode string, year int) *AdvancementReport {
	// Get the event details
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events := db.GetAllEvents(filter)
	if len(events) == 0 {
		return nil
	}

	// Find the event matching the year
	var event *database.Event
	for _, e := range events {
		if e.Year == year {
			event = e
			break
		}
	}

	if event == nil {
		return nil
	}

	// Get rankings for the event
	rankings := db.GetEventRankings(event.EventID)
	if len(rankings) == 0 {
		return &AdvancementReport{
			Event:            event,
			TeamAdvancements: []*TeamAdvancement{},
		}
	}

	// Get advancements for the event
	advancements := db.GetEventAdvancements(event.EventID)
	advancementMap := make(map[int]bool)
	advancementStatusMap := make(map[int]string)
	for _, adv := range advancements {
		advancementMap[adv.TeamID] = true
		advancementStatusMap[adv.TeamID] = adv.Status
	}

	// Get awards for judging points calculation
	awards := db.GetEventAwards(event.EventID)
	judgingPointsMap := calculateJudgingPoints(awards)
	playoffPointsMap := calculatePlayoffPoints(event)
	selectionPointsMap := calculateSelectionPoints(event)
	qualificationPointsMap := calculateQualificationPoints(rankings)

	// Build team advancement records
	var teamAdvancements []*TeamAdvancement
	for _, ranking := range rankings {
		team := db.GetTeam(ranking.TeamID)
		if team == nil {
			continue
		}

		// Get qualification points for this team
		qualificationPoints := qualificationPointsMap[ranking.TeamID]

		// Get judging points for this team
		judgingPoints := judgingPointsMap[ranking.TeamID]

		// Get playoff points for this team
		playoffPoints := playoffPointsMap[ranking.TeamID]

		// Get selection points for this team
		selectionPoints := selectionPointsMap[ranking.TeamID]

		// Calculate total points
		totalPoints := judgingPoints + qualificationPoints + playoffPoints + selectionPoints

		teamAdv := &TeamAdvancement{
			Rank:                ranking.Rank,
			Team:                team,
			Ranking:             ranking,
			TotalPoints:         totalPoints,
			JudgingPoints:       judgingPoints,
			PlayoffPoints:       playoffPoints,
			SelectionPoints:     selectionPoints,
			QualificationPoints: qualificationPoints,
			Advances:            advancementMap[ranking.TeamID],
			Status:              advancementStatusMap[ranking.TeamID],
		}

		teamAdvancements = append(teamAdvancements, teamAdv)
	}

	// Sort by total points (descending) to assign advancement numbers
	slices.SortFunc(teamAdvancements, func(a, b *TeamAdvancement) int {
		if a.TotalPoints != b.TotalPoints {
			if a.TotalPoints > b.TotalPoints {
				return -1
			}
			return 1
		}
		// Tie-breaker: use original qualification rank
		return a.Ranking.Rank - b.Ranking.Rank
	})

	// Assign advancement numbers and update rank to match sorted order
	advancementRank := 1
	for i, ta := range teamAdvancements {
		// Update rank to match sorted position
		ta.Rank = i + 1

		// Assign advancement number if team advances
		if ta.Advances {
			// Skip teams that have "already advanced" status when assigning numbers
			if ta.Status == "already advanced" {
				ta.AdvancementNumber = "-"
			} else {
				ta.AdvancementNumber = fmt.Sprintf("%d", advancementRank)
				advancementRank++
			}
		} else {
			ta.AdvancementNumber = "-"
		}
	}

	// Teams remain sorted by total points (descending) for display

	return &AdvancementReport{
		Event:            event,
		TeamAdvancements: teamAdvancements,
	}
}

// calculateJudgingPoints calculates judging points based on awards.
// Points are awarded as follows:
// - Inspire 1: 60 points, Inspire 2: 30 points, Inspire 3: 15 points
// - Other judged awards: 1st place (series 1): 12 points, 2nd place (series 2): 6 points, 3rd place (series 3): 3 points
func calculateJudgingPoints(awards []*database.EventAward) map[int]int {
	pointsMap := make(map[int]int)

	for _, award := range awards {
		// Skip playoff awards (winning/finalist alliance)
		if isPlayoffAward(award.Name) {
			continue
		}

		var points int
		awardNameLower := award.Name

		// Assign points based on award type and series
		if containsIgnoreCase(awardNameLower, "inspire") {
			// Inspire awards have special point values
			switch award.Series {
			case 1:
				points = 60
			case 2:
				points = 30
			case 3:
				points = 15
			}
		} else if isJudgedAward(awardNameLower) {
			// Other judged awards use standard point scale based on series
			switch award.Series {
			case 1:
				points = 12
			case 2:
				points = 6
			case 3:
				points = 3
			}
		}

		pointsMap[award.TeamID] += points
	}

	return pointsMap
}

// calculatePlayoffPoints calculates playoff points based on how far teams progress in the playoff bracket.
// Points are awarded as follows:
// - Winning Alliance: 40 points
// - Finalist Alliance: 20 points
// - 3rd Place: 10 points (highest scoring losing semifinalist)
// - 4th Place: 5 points (lowest scoring losing semifinalist)
//
// This handles both single-elimination and modified double-elimination (winners/losers bracket) formats.
func calculatePlayoffPoints(event *database.Event) map[int]int {
	pointsMap := make(map[int]int)

	// Get all matches for the event
	matches := db.GetMatchesByEvent(event.EventID)

	// Filter for playoff matches only
	var playoffMatches []*database.Match
	for _, match := range matches {
		if strings.EqualFold(match.TournamentLevel, string(ftc.PLAYOFF)) {
			playoffMatches = append(playoffMatches, match)
		}
	}

	if len(playoffMatches) == 0 {
		return pointsMap
	}

	// Sort playoff matches by match number to identify finals (highest number)
	slices.SortFunc(playoffMatches, func(a, b *database.Match) int {
		return b.MatchNumber - a.MatchNumber // Descending order
	})

	// The first match after sorting should be the finals (championship)
	finalsMatch := playoffMatches[0]

	// Get alliance scores for finals
	redScore := db.GetMatchAllianceScore(finalsMatch.MatchID, database.AllianceRed)
	blueScore := db.GetMatchAllianceScore(finalsMatch.MatchID, database.AllianceBlue)

	if redScore != nil && blueScore != nil {
		var winningAlliance, losingAlliance string

		if redScore.TotalPoints > blueScore.TotalPoints {
			winningAlliance = database.AllianceRed
			losingAlliance = database.AllianceBlue
		} else {
			winningAlliance = database.AllianceBlue
			losingAlliance = database.AllianceRed
		}

		// Assign 40 points to winning alliance teams
		winningTeams := db.GetMatchTeams(finalsMatch.MatchID)
		for _, mt := range winningTeams {
			if mt.Alliance == winningAlliance {
				pointsMap[mt.TeamID] = 40
			}
		}

		// Assign 20 points to finalist alliance teams
		finalistTeams := db.GetMatchTeams(finalsMatch.MatchID)
		for _, mt := range finalistTeams {
			if mt.Alliance == losingAlliance {
				pointsMap[mt.TeamID] = 20
			}
		}
	}

	// Find the losing alliances from semifinal matches for 3rd and 4th place
	// In a winners/losers bracket:
	// - The finals teams already have their points
	// - We need to find the losing alliances from earlier playoff rounds

	type SemifinalistAlliance struct {
		alliance string
		matchID  string
		score    int
		teams    []int
	}
	var semifinalLosers []SemifinalistAlliance

	// Get the team IDs that played in finals
	finalsTeamIDs := make(map[int]bool)
	finalsTeams := db.GetMatchTeams(finalsMatch.MatchID)
	for _, mt := range finalsTeams {
		finalsTeamIDs[mt.TeamID] = true
	}

	// Look through all playoff matches (except finals) to find losing alliances
	// that didn't make it to the championship finals
	for i := 1; i < len(playoffMatches); i++ {
		match := playoffMatches[i]

		// Get teams and scores for this match
		matchTeams := db.GetMatchTeams(match.MatchID)
		redScore := db.GetMatchAllianceScore(match.MatchID, database.AllianceRed)
		blueScore := db.GetMatchAllianceScore(match.MatchID, database.AllianceBlue)

		if redScore == nil || blueScore == nil {
			continue
		}

		// Determine winning and losing alliance
		var losingAlliance string
		var losingScore int

		if redScore.TotalPoints < blueScore.TotalPoints {
			losingAlliance = database.AllianceRed
			losingScore = redScore.TotalPoints
		} else {
			losingAlliance = database.AllianceBlue
			losingScore = blueScore.TotalPoints
		}

		// Get the team IDs from the losing alliance
		var losingTeamIDs []int
		for _, mt := range matchTeams {
			if mt.Alliance == losingAlliance {
				losingTeamIDs = append(losingTeamIDs, mt.TeamID)
			}
		}

		// Check if any of the losing teams made it to finals
		// If none of them made finals, this is a semifinal loser eligible for 3rd/4th place
		madeItToFinals := false
		for _, teamID := range losingTeamIDs {
			if finalsTeamIDs[teamID] {
				madeItToFinals = true
				break
			}
		}

		if !madeItToFinals && len(losingTeamIDs) > 0 {
			semifinalLosers = append(semifinalLosers, SemifinalistAlliance{
				alliance: losingAlliance,
				matchID:  match.MatchID,
				score:    losingScore,
				teams:    losingTeamIDs,
			})
		}
	}

	// Sort semifinal losers by score to determine 3rd vs 4th place
	// Higher score gets 3rd place (10 points), lower score gets 4th place (5 points)
	slices.SortFunc(semifinalLosers, func(a, b SemifinalistAlliance) int {
		// Descending order: higher scores first
		if a.score > b.score {
			return -1
		}
		if a.score < b.score {
			return 1
		}
		return 0
	})

	// Assign points to the top 2 semifinal losers (3rd and 4th place)
	for i := 0; i < len(semifinalLosers) && i < 2; i++ {
		loser := semifinalLosers[i]
		points := 10 // 3rd place
		if i == 1 {
			points = 5 // 4th place
		}

		// Assign points to all teams in the losing alliance
		for _, teamID := range loser.teams {
			// Only assign points if team doesn't already have playoff points
			if _, exists := pointsMap[teamID]; !exists {
				pointsMap[teamID] = points
			}
		}
	}

	return pointsMap
}

// calculateSelectionPoints calculates selection points based on alliance selection.
// Points are awarded as follows:
// - 1st alliance: 20 points
// - 2nd alliance: 19 points
// - 3rd alliance: 18 points
// - And so on: 20 - (alliance_number - 1)
func calculateSelectionPoints(event *database.Event) map[int]int {
	pointsMap := make(map[int]int)

	// Fetch alliance data from FTC API
	alliances, err := ftc.GetEventAlliances(strconv.Itoa(event.Year), event.EventCode)
	if err != nil {
		slog.Warn("Failed to fetch alliances for selection points", "eventCode", event.EventCode, "year", event.Year, "error", err)
		return pointsMap
	}

	// Assign points based on alliance number
	for _, alliance := range alliances {
		if alliance.Number <= 0 {
			continue
		}

		// Calculate points: 20 for 1st alliance, 19 for 2nd, etc.
		points := 20 - (alliance.Number - 1)
		if points < 0 {
			points = 0
		}

		// Assign points to all teams in the alliance
		if alliance.Captain > 0 {
			pointsMap[alliance.Captain] = points
		}
		if alliance.Round1 > 0 {
			pointsMap[alliance.Round1] = points
		}
		if alliance.Round2 > 0 {
			pointsMap[alliance.Round2] = points
		}
		if alliance.Round3 != nil && *alliance.Round3 > 0 {
			pointsMap[*alliance.Round3] = points
		}
	}

	return pointsMap
}

// calculateQualificationPoints calculates qualification points based on ranking scores.
// Points are awarded as follows:
// - Highest ranking score: 16 points
// - Each lower ranking score: 1 point less
// - Lowest ranking score: minimum 2 points
// - Teams with the same ranking score get the same points
// - After multiple teams with the same score, the next lower score only loses 1 point (not skipping)
func calculateQualificationPoints(rankings []*database.EventRanking) map[int]int {
	pointsMap := make(map[int]int)

	if len(rankings) == 0 {
		return pointsMap
	}

	// Sort rankings by ranking score (SortOrder1) in descending order
	sortedRankings := make([]*database.EventRanking, len(rankings))
	copy(sortedRankings, rankings)
	slices.SortFunc(sortedRankings, func(a, b *database.EventRanking) int {
		if a.SortOrder1 > b.SortOrder1 {
			return -1
		}
		if a.SortOrder1 < b.SortOrder1 {
			return 1
		}
		return 0
	})

	// Assign points starting at 16, decreasing by 1 for each unique ranking score, minimum 2
	maxPoints := 16
	minPoints := 2
	var previousSortOrder1 float64
	var previousScore int
	numTeams := len(sortedRankings)

	multiplier := float64(maxPoints-minPoints) / float64(numTeams)

	for i, ranking := range sortedRankings {
		rank := numTeams - i
		score := int(math.Floor(multiplier*float64(rank))) + minPoints

		if previousScore != 0 && ranking.SortOrder1 == previousSortOrder1 {
			score = previousScore
		} else {
			previousScore = score
		}
		pointsMap[ranking.TeamID] = score

		previousScore = score
		previousSortOrder1 = ranking.SortOrder1
	}

	currentScure := maxPoints + 1
	previousSortOrder1 = 0
	for _, ranking := range sortedRankings {
		var score int
		if ranking.SortOrder1 == previousSortOrder1 {
			score = previousScore
		} else {
			currentScure--
			score = currentScure
		}
		if score < minPoints {
			score = minPoints
		}
		pointsMap[ranking.TeamID] = score

		previousScore = score
		previousSortOrder1 = ranking.SortOrder1
	}

	return pointsMap
}

// isPlayoffAward returns true if the award is a playoff-related award.
func isPlayoffAward(awardName string) bool {
	return containsIgnoreCase(awardName, "winning alliance") ||
		containsIgnoreCase(awardName, "finalist alliance")
}

// isJudgedAward returns true if the award is a judged award (not alliance/playoff awards).
func isJudgedAward(awardName string) bool {
	// Check if it's a playoff award first
	if isPlayoffAward(awardName) {
		return false
	}

	// Other awards are typically judged
	return containsIgnoreCase(awardName, "award") ||
		containsIgnoreCase(awardName, "innovate") ||
		containsIgnoreCase(awardName, "design") ||
		containsIgnoreCase(awardName, "control") ||
		containsIgnoreCase(awardName, "motivate") ||
		containsIgnoreCase(awardName, "compass") ||
		containsIgnoreCase(awardName, "promote") ||
		containsIgnoreCase(awardName, "think") ||
		containsIgnoreCase(awardName, "connect") ||
		containsIgnoreCase(awardName, "sustain") ||
		containsIgnoreCase(awardName, "reach")
}

// containsIgnoreCase checks if a string contains a substring (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
