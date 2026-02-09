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

	for _, match := range playoffMatches {
		// Get alliance scores for finals
		redScore := db.GetMatchAllianceScore(match.MatchID, database.AllianceRed)
		blueScore := db.GetMatchAllianceScore(match.MatchID, database.AllianceBlue)

		if redScore != nil && blueScore != nil {
			var winningAlliance string
			if redScore.TotalPoints > blueScore.TotalPoints {
				winningAlliance = database.AllianceRed
			} else {
				winningAlliance = database.AllianceBlue
			}

			var winningPoints, losingPoints int
			switch len(pointsMap) {
			case 0:
				winningPoints = 40
				losingPoints = 20
			case 4:
				losingPoints = 10
			case 6:
				losingPoints = 5
			default:
				break
			}

			// Assign 40 points to winning alliance teams
			teams := db.GetMatchTeams(match.MatchID)
			for _, mt := range teams {
				if pointsMap[mt.TeamID] == 0 {
					if mt.Alliance == winningAlliance {
						pointsMap[mt.TeamID] = winningPoints
					} else {
						pointsMap[mt.TeamID] = losingPoints
					}
				}
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

	N := len(sortedRankings)
	for i, ranking := range sortedRankings {
		R := i + 1
		pointsMap[ranking.TeamID] = ftcQualificationPoints(R, N)
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

// ftcQualificationPoints computes FTC Qualification Phase Performance points
func ftcQualificationPoints(rank, teams int) int {
	alpha := 1.07

	r := float64(rank)
	n := float64(teams)

	x := (n - 2*r + 2) / (alpha * n)

	scale := 7.0 / math.Erfinv(1.0/alpha)
	points := math.Erfinv(x)*scale + 9.0

	return int(math.Ceil(points))
}

// EventParticipation represents a team's participation at an event including awards.
type EventParticipation struct {
	Event  *database.Event
	Awards []*database.EventAward
}

// RegionTeamAdvancement represents a team's advancement information across all events in a region.
type RegionTeamAdvancement struct {
	Team                     *database.Team
	AdvancingEvent           *database.Event        // The event from which the team advanced
	AdvancingEventAwards     []*database.EventAward // Awards from the advancing event
	OtherEventParticipations []*EventParticipation  // Other events the team participated in
}

// RegionAdvancementReport represents all teams advancing from events in a region.
type RegionAdvancementReport struct {
	RegionCode       string
	Year             int
	TeamAdvancements []*RegionTeamAdvancement
}

// RegionAdvancementQuery retrieves advancement information for all teams advancing in a region.
// It returns a RegionAdvancementReport with teams sorted by team number.
func RegionAdvancementQuery(regionCode string, year int) *RegionAdvancementReport {
	// Get all events in the region for the given year
	filter := database.EventFilter{
		RegionCodes: []string{regionCode},
	}
	allEvents := db.GetAllEvents(filter)

	// Filter events by year
	var events []*database.Event
	for _, e := range allEvents {
		if e.Year == year {
			events = append(events, e)
		}
	}

	if len(events) == 0 {
		return &RegionAdvancementReport{
			RegionCode:       regionCode,
			Year:             year,
			TeamAdvancements: []*RegionTeamAdvancement{},
		}
	}

	// Sort events by date to ensure we process them chronologically
	slices.SortFunc(events, func(a, b *database.Event) int {
		return a.DateStart.Compare(b.DateStart)
	})

	// Track all advancement records for each team
	type teamAdvancementRecord struct {
		event  *database.Event
		status string
	}
	teamAdvancementRecords := make(map[int][]teamAdvancementRecord)

	// Track all events each team participated in
	teamEventParticipationMap := make(map[int][]*database.Event)

	// Track awards for each team at each event
	teamEventAwardsMap := make(map[int]map[string][]*database.EventAward) // teamID -> eventID -> awards

	// First pass: collect all advancements, participations, and awards
	for _, event := range events {
		// Get advancements for this event
		advancements := db.GetEventAdvancements(event.EventID)
		for _, adv := range advancements {
			// Track all advancement records for this team
			teamAdvancementRecords[adv.TeamID] = append(teamAdvancementRecords[adv.TeamID], teamAdvancementRecord{
				event:  event,
				status: adv.Status,
			})
		}

		// Get all teams that participated in this event
		eventTeams := db.GetEventTeams(event.EventID)
		for _, et := range eventTeams {
			teamEventParticipationMap[et.TeamID] = append(teamEventParticipationMap[et.TeamID], event)
		}

		// Get awards for this event
		awards := db.GetEventAwards(event.EventID)
		for _, award := range awards {
			if teamEventAwardsMap[award.TeamID] == nil {
				teamEventAwardsMap[award.TeamID] = make(map[string][]*database.EventAward)
			}
			teamEventAwardsMap[award.TeamID][event.EventID] = append(teamEventAwardsMap[award.TeamID][event.EventID], award)
		}
	}

	// Determine the advancing event for each team
	teamAdvancingEventMap := make(map[int]*database.Event)
	for teamID, records := range teamAdvancementRecords {
		// Find ALL events where status is NOT "already_advancing"
		var advancingEvents []*database.Event
		for _, record := range records {
			if record.status != "already_advancing" {
				advancingEvents = append(advancingEvents, record.event)
			}
		}

		// If we found advancing events, use the earliest one
		if len(advancingEvents) > 0 {
			// Sort by date to find the earliest
			slices.SortFunc(advancingEvents, func(a, b *database.Event) int {
				return a.DateStart.Compare(b.DateStart)
			})
			teamAdvancingEventMap[teamID] = advancingEvents[0]
		} else if len(records) > 0 {
			// If no non-"already_advancing" event found, use the earliest event overall
			var allEvents []*database.Event
			for _, record := range records {
				allEvents = append(allEvents, record.event)
			}
			slices.SortFunc(allEvents, func(a, b *database.Event) int {
				return a.DateStart.Compare(b.DateStart)
			})
			teamAdvancingEventMap[teamID] = allEvents[0]
		}
	}

	// Build RegionTeamAdvancement records for advancing teams
	var teamAdvancements []*RegionTeamAdvancement
	for teamID, advancingEvent := range teamAdvancingEventMap {
		team := db.GetTeam(teamID)
		if team == nil {
			continue
		}

		// Get awards from the advancing event
		var advancingEventAwards []*database.EventAward
		if teamEventAwardsMap[teamID] != nil {
			advancingEventAwards = teamEventAwardsMap[teamID][advancingEvent.EventID]
			// Sort awards alphabetically by name
			slices.SortFunc(advancingEventAwards, func(a, b *database.EventAward) int {
				return strings.Compare(a.Name, b.Name)
			})
		}

		// Get all other events this team participated in
		var otherParticipations []*EventParticipation
		allParticipations := teamEventParticipationMap[teamID]
		for _, event := range allParticipations {
			// Skip the advancing event
			if event.EventID == advancingEvent.EventID {
				continue
			}

			// Get awards from this event
			var eventAwards []*database.EventAward
			if teamEventAwardsMap[teamID] != nil {
				eventAwards = teamEventAwardsMap[teamID][event.EventID]
				// Sort awards alphabetically by name
				slices.SortFunc(eventAwards, func(a, b *database.EventAward) int {
					return strings.Compare(a.Name, b.Name)
				})
			}

			otherParticipations = append(otherParticipations, &EventParticipation{
				Event:  event,
				Awards: eventAwards,
			})
		}

		// Sort other participations by event date
		slices.SortFunc(otherParticipations, func(a, b *EventParticipation) int {
			return a.Event.DateStart.Compare(b.Event.DateStart)
		})

		teamAdvancements = append(teamAdvancements, &RegionTeamAdvancement{
			Team:                     team,
			AdvancingEvent:           advancingEvent,
			AdvancingEventAwards:     advancingEventAwards,
			OtherEventParticipations: otherParticipations,
		})
	}

	// Sort by team number
	slices.SortFunc(teamAdvancements, func(a, b *RegionTeamAdvancement) int {
		return a.Team.TeamID - b.Team.TeamID
	})

	return &RegionAdvancementReport{
		RegionCode:       regionCode,
		Year:             year,
		TeamAdvancements: teamAdvancements,
	}
}
