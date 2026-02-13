package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
)

type Server struct {
	db     database.DB
	mux    *http.ServeMux
	logger *slog.Logger
}

// Response types for event resources - grouped under event

// EventResponse represents an event without event_id
type EventResponse struct {
	EventCode    string `json:"event_code"`
	Year         int    `json:"year"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	DivisionCode string `json:"division_code"`
	RegionCode   string `json:"region_code"`
	LeagueCode   string `json:"league_code"`
	Venue        string `json:"venue"`
	Address      string `json:"address"`
	City         string `json:"city"`
	StateProv    string `json:"state_prov"`
	Country      string `json:"country"`
	Timezone     string `json:"timezone"`
	DateStart    string `json:"date_start"`
	DateEnd      string `json:"date_end"`
}

// MatchResponse represents a match without event_id
type MatchResponse struct {
	MatchType       string `json:"matchType"`
	MatchNumber     int    `json:"matchNumber"`
	ActualStartTime string `json:"actualStartTime"`
	Description     string `json:"description"`
	TournamentLevel string `json:"tournamentLevel"`
}

type AwardResponse struct {
	Name string         `json:"name"`
	Team *database.Team `json:"team"`
}

type RankingResponse struct {
	Team           *database.Team `json:"team"`
	Year           int            `json:"year"`
	SortOrder1     float64        `json:"sort_order1"`
	SortOrder2     float64        `json:"sort_order2"`
	SortOrder3     float64        `json:"sort_order3"`
	SortOrder4     float64        `json:"sort_order4"`
	SortOrder5     float64        `json:"sort_order5"`
	SortOrder6     float64        `json:"sort_order6"`
	Wins           int            `json:"wins"`
	Losses         int            `json:"losses"`
	Ties           int            `json:"ties"`
	Dq             int            `json:"dq"`
	MatchesPlayed  int            `json:"matches_played"`
	MatchesCounted int            `json:"matches_counted"`
	HighMatchScore int            `json:"high_match_score"`
}

// EventWithTeams represents an event along with its teams
type EventWithTeams struct {
	*EventResponse
	Teams []*database.Team `json:"teams"`
}

// EventTeamsResponse represents the response for an event's teams endpoint
type EventTeamsResponse struct {
	Event *EventWithTeams `json:"event"`
}

// EventRankingsResponse represents the response for an event's rankings endpoint
type EventRankingsResponse struct {
	Event    *EventResponse    `json:"event"`
	Rankings []RankingResponse `json:"rankings"`
}

// EventWithAwards represents an event along with its awards
type EventWithAwards struct {
	*EventResponse
	Awards []AwardResponse `json:"awards"`
}

// EventAwardsResponse represents the response for an event's awards endpoint
type EventAwardsResponse struct {
	Event *EventWithAwards `json:"event"`
}

// MatchAllianceScoreResponse represents the score details for an alliance in a match
type MatchAllianceScoreResponse struct {
	AutoPoints          int `json:"auto_points"`
	TeleopPoints        int `json:"teleop_points"`
	FoulPointsCommitted int `json:"foul_points_committed"`
	PreFoulTotal        int `json:"pre_foul_total"`
	TotalPoints         int `json:"total_points"`
	MajorFouls          int `json:"major_fouls"`
	MinorFouls          int `json:"minor_fouls"`
}

// MatchAllianceDetailsResponse represents the details of an alliance in a match, including score and teams
type MatchAllianceDetailsResponse struct {
	Alliance string                      `json:"alliance"`
	Score    *MatchAllianceScoreResponse `json:"score"`
	Teams    []*database.Team            `json:"teams"`
}

// MatchWithAlliancesResponse represents a match along with its alliance details
type MatchWithAlliancesResponse struct {
	MatchType       string                        `json:"matchType"`
	MatchNumber     int                           `json:"matchNumber"`
	ActualStartTime string                        `json:"actualStartTime"`
	Description     string                        `json:"description"`
	TournamentLevel string                        `json:"tournamentLevel"`
	RedAlliance     *MatchAllianceDetailsResponse `json:"red_alliance"`
	BlueAlliance    *MatchAllianceDetailsResponse `json:"blue_alliance"`
}

// TeamMatchResultResponse represents a match result for a specific team, including alliance details and result
type TeamMatchResultResponse struct {
	MatchType       string                        `json:"matchType"`
	MatchNumber     int                           `json:"matchNumber"`
	ActualStartTime string                        `json:"actualStartTime"`
	Description     string                        `json:"description"`
	TournamentLevel string                        `json:"tournamentLevel"`
	RedAlliance     *MatchAllianceDetailsResponse `json:"red_alliance"`
	BlueAlliance    *MatchAllianceDetailsResponse `json:"blue_alliance"`
	Team            *database.Team                `json:"team"`
	Result          string                        `json:"result"`
}

// EventWithMatches represents an event along with its matches (with alliance details if team filter is not applied)
type EventWithMatches struct {
	*EventResponse
	Matches interface{} `json:"matches"`
}

// EventMatchesResponse represents the response for an event's matches endpoint
type EventMatchesResponse struct {
	Event *EventWithMatches `json:"event"`
}

// EventAdvancementResponse represents the response for an event's advancement endpoint, including team advancements
type EventAdvancementResponse struct {
	Event            *EventResponse           `json:"event"`
	TeamAdvancements []*query.TeamAdvancement `json:"team_advancements"`
}

// TeamPerformanceResponse represents the performance metrics for a team across events in a season
type PerformanceResponse struct {
	TeamID   int     `json:"team_id"`
	TeamName string  `json:"team_name"`
	Region   string  `json:"region"`
	OPR      float64 `json:"opr"`
	NpOPR    float64 `json:"np_opr"`
	CCWM     float64 `json:"ccwm"`
	DPR      float64 `json:"dpr"`
	NpDPR    float64 `json:"np_dpr"`
	NpAVG    float64 `json:"np_avg"`
	Matches  int     `json:"matches"`
}

// EventPerformanceResponse represents the performance metrics for a team at a specific event in a season
type EventPerformanceResponse struct {
	TeamID    int     `json:"team_id"`
	TeamName  string  `json:"team_name"`
	Region    string  `json:"region"`
	Year      int     `json:"year"`
	EventCode string  `json:"event_code"`
	EventName string  `json:"event_name"`
	OPR       float64 `json:"opr"`
	NpOPR     float64 `json:"np_opr"`
	CCWM      float64 `json:"ccwm"`
	DPR       float64 `json:"dpr"`
	NpDPR     float64 `json:"np_dpr"`
	NpAVG     float64 `json:"np_avg"`
	Matches   int     `json:"matches"`
}

// QualifiedTeamResponse represents a team that qualified for advancement, along with whether they were already qualified and the first event they qualified at
type QualifiedTeamResponse struct {
	Team              *database.Team `json:"team"`
	AlreadyQualified  bool           `json:"already_qualified"`
	FirstQualifyEvent string         `json:"first_qualify_event"`
}

// EventQualifiedTeamsResponse represents an event along with the teams that qualified for advancement from that event
type EventQualifiedTeamsResponse struct {
	Event          *EventResponse           `json:"event"`
	QualifiedTeams []*QualifiedTeamResponse `json:"qualified_teams"`
}

// EventAdvancementSummaryResponse represents the advancement summary for a region and season, including summaries for each event in the region
type EventAdvancementSummaryResponse struct {
	RegionCode     string                         `json:"region_code"`
	Year           int                            `json:"year"`
	EventSummaries []*EventQualifiedTeamsResponse `json:"event_summaries"`
}

// Helper functions to convert database types to response types

// toEventResponse converts a database.Event to an EventResponse, which is used in API responses without exposing internal event_id
func toEventResponse(e *database.Event) *EventResponse {
	if e == nil {
		return nil
	}
	return &EventResponse{
		EventCode:    e.EventCode,
		Year:         e.Year,
		Name:         e.Name,
		Type:         e.Type,
		DivisionCode: e.DivisionCode,
		RegionCode:   e.RegionCode,
		LeagueCode:   e.LeagueCode,
		Venue:        e.Venue,
		Address:      e.Address,
		City:         e.City,
		StateProv:    e.StateProv,
		Country:      e.Country,
		Timezone:     e.Timezone,
		DateStart:    e.DateStart.Format("2006-01-02T15:04:05Z07:00"),
		DateEnd:      e.DateEnd.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// toMatchResponse converts a database.Match to a MatchResponse, which is used in API responses without exposing internal event_id
func toMatchAllianceScoreResponse(mas *database.MatchAllianceScore) *MatchAllianceScoreResponse {
	if mas == nil {
		return nil
	}
	return &MatchAllianceScoreResponse{
		AutoPoints:          mas.AutoPoints,
		TeleopPoints:        mas.TeleopPoints,
		FoulPointsCommitted: mas.FoulPointsCommitted,
		PreFoulTotal:        mas.PreFoulTotal,
		TotalPoints:         mas.TotalPoints,
		MajorFouls:          mas.MajorFouls,
		MinorFouls:          mas.MinorFouls,
	}
}

// toMatchAllianceDetailsResponse converts a query.MatchAllianceDetails to a MatchAllianceDetailsResponse, which includes alliance name, score details, and teams
func toMatchAllianceDetailsResponse(mad *query.MatchAllianceDetails) *MatchAllianceDetailsResponse {
	if mad == nil {
		return nil
	}
	return &MatchAllianceDetailsResponse{
		Alliance: mad.Alliance,
		Score:    toMatchAllianceScoreResponse(mad.Score),
		Teams:    mad.Teams,
	}
}

// toMatchWithAlliancesResponse converts a database.Match along with its alliance details to a MatchWithAlliancesResponse, which is used in API responses without exposing internal event_id
func toMatchWithAlliancesResponse(m *database.Match, red, blue *query.MatchAllianceDetails) *MatchWithAlliancesResponse {
	if m == nil {
		return nil
	}
	return &MatchWithAlliancesResponse{
		MatchType:       m.MatchType,
		MatchNumber:     m.MatchNumber,
		ActualStartTime: m.ActualStartTime,
		Description:     m.Description,
		TournamentLevel: m.TournamentLevel,
		RedAlliance:     toMatchAllianceDetailsResponse(red),
		BlueAlliance:    toMatchAllianceDetailsResponse(blue),
	}
}

// toTeamMatchResultResponse converts a query.TeamMatchResult to a TeamMatchResultResponse, which includes match details, alliance details, and the team's result in the match
func toTeamMatchResultResponse(tmr *query.TeamMatchResult) *TeamMatchResultResponse {
	if tmr == nil {
		return nil
	}
	return &TeamMatchResultResponse{
		MatchType:       tmr.Match.MatchType,
		MatchNumber:     tmr.Match.MatchNumber,
		ActualStartTime: tmr.Match.ActualStartTime,
		Description:     tmr.Match.Description,
		TournamentLevel: tmr.Match.TournamentLevel,
		RedAlliance:     toMatchAllianceDetailsResponse(tmr.TeamAlliance),
		BlueAlliance:    toMatchAllianceDetailsResponse(tmr.OpponentAlliance),
		Team:            tmr.Team,
		Result:          tmr.Result,
	}
}

// toEventAdvancementSummaryResponse converts a query.EventAdvancementSummary to an EventAdvancementSummaryResponse, which includes the region code, year, and summaries for each event in the region
func toEventAdvancementSummaryResponse(summary *query.EventAdvancementSummary) *EventAdvancementSummaryResponse {
	if summary == nil {
		return nil
	}

	eventSummaries := make([]*EventQualifiedTeamsResponse, 0, len(summary.EventSummaries))
	for _, es := range summary.EventSummaries {
		qualifiedTeams := make([]*QualifiedTeamResponse, 0, len(es.QualifiedTeams))
		for _, qt := range es.QualifiedTeams {
			qualifiedTeams = append(qualifiedTeams, &QualifiedTeamResponse{
				Team:              qt.Team,
				AlreadyQualified:  qt.AlreadyQualified,
				FirstQualifyEvent: qt.FirstQualifyEvent,
			})
		}
		eventSummaries = append(eventSummaries, &EventQualifiedTeamsResponse{
			Event:          toEventResponse(es.Event),
			QualifiedTeams: qualifiedTeams,
		})
	}

	return &EventAdvancementSummaryResponse{
		RegionCode:     summary.RegionCode,
		Year:           summary.Year,
		EventSummaries: eventSummaries,
	}
}

// NewServer creates a new Server instance with the given database connection and sets up the routes
func NewServer(db database.DB) *Server {
	s := &Server{
		db:     db,
		mux:    http.NewServeMux(),
		logger: slog.Default(),
	}
	s.setupRoutes()
	return s
}

// setupRoutes registers the HTTP handlers for the server's endpoints
func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/v1/", s.handleV1Routes)
	s.mux.HandleFunc("/health", s.handleHealth)
}

// ServeHTTP allows Server to satisfy the http.Handler interface by delegating to the internal ServeMux
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// parseLimit extracts the 'limit' query parameter from the request and converts it to an integer. It returns an error if the limit is invalid.
func (s *Server) parseLimit(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return 0, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("invalid limit: %s", limitStr)
	}
	if limit < 0 {
		return 0, fmt.Errorf("limit must be non-negative")
	}
	return limit, nil
}

// handleHealth responds with a simple JSON indicating the server is healthy. This can be used for health checks.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleV1Routes is the main handler for all /v1/ API routes. It parses the URL path to determine the requested resource and delegates to the appropriate handler function based on the resource type.
func (s *Server) handleV1Routes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		s.writeError(w, http.StatusBadRequest, "season is required in path")
		return
	}

	season := parts[0]
	year, err := strconv.Atoi(season)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid season: %s", season))
		return
	}

	if len(parts) < 2 {
		s.writeError(w, http.StatusBadRequest, "resource type is required")
		return
	}

	resource := parts[1]

	switch resource {
	case "team":
		s.handleTeam(w, r, year, parts[2:])
	case "teams":
		s.handleTeams(w, r, year, parts[2:])
	case "events":
		s.handleEvents(w, r, year, parts[2:])
	case "team-rankings":
		s.handleTeamRankings(w, r, year, parts[2:])
	case "team-event-rankings":
		s.handleTeamEventRankings(w, r, year, parts[2:])
	case "regions":
		s.handleRegions(w, r, year, parts[2:])
	case "advancement":
		s.handleAllAdvancement(w, r, year, parts[2:])
	default:
		s.writeError(w, http.StatusNotFound, fmt.Sprintf("unknown resource: %s", resource))
	}
}

// handleTeam handles requests for a specific team's details. It expects the team ID to be provided in the URL path and returns the team's information in JSON format.
func (s *Server) handleTeam(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	if len(parts) < 1 {
		s.writeError(w, http.StatusBadRequest, "teamID is required")
		return
	}

	teamID, err := strconv.Atoi(parts[0])
	if err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid teamID: %s", parts[0]))
		return
	}

	details := query.TeamDetailsQuery(teamID)
	if details == nil {
		s.writeError(w, http.StatusNotFound, fmt.Sprintf("team %d not found", teamID))
		return
	}

	s.writeJSON(w, http.StatusOK, details)
}

// handleTeams handles requests for teams, optionally filtered by region. It supports a 'limit' query parameter to limit the number of teams returned. If a region is specified in the URL path, it filters teams by that region; otherwise, it returns all teams.
func (s *Server) handleTeams(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var teams []*database.Team
	if len(parts) > 0 && parts[0] != "" {
		// Region specified - filter by region
		region := parts[0]
		teamsFilter := database.TeamFilter{
			HomeRegions: []string{region},
		}
		teams = query.TeamsQuery(teamsFilter)
	} else {
		// No region specified - get all teams
		teams = query.TeamsQuery(database.TeamFilter{})
	}

	if limit > 0 && limit < len(teams) {
		teams = teams[:limit]
	}

	s.writeJSON(w, http.StatusOK, teams)
}

// handleEvents handles requests for events, optionally filtered by event code. It expects the event code to be provided in the URL path and supports a 'limit' query parameter to limit the number of events returned. It delegates to specific handlers for different event resources such as teams, rankings, awards, advancement, and matches based on the second part of the URL path.
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	if len(parts) < 1 {
		s.writeError(w, http.StatusBadRequest, "eventCode is required")
		return
	}

	eventCode := parts[0]

	if len(parts) < 2 {
		s.writeError(w, http.StatusBadRequest, "event resource type is required")
		return
	}

	resource := parts[1]

	switch resource {
	case "teams":
		s.handleEventTeams(w, r, year, eventCode)
	case "rankings":
		s.handleEventRankings(w, r, year, eventCode)
	case "awards":
		s.handleEventAwards(w, r, year, eventCode)
	case "advancement":
		s.handleEventAdvancement(w, r, year, eventCode)
	case "matches":
		s.handleEventMatches(w, r, year, eventCode)
	default:
		s.writeError(w, http.StatusNotFound, fmt.Sprintf("unknown event resource: %s", resource))
	}
}

// handleEventTeams handles requests for the teams participating in a specific event. It expects the event code to be provided in the URL path and supports a 'limit' query parameter to limit the number of teams returned. It returns the event details along with the list of teams in JSON format.
func (s *Server) handleEventTeams(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	eventTeams := query.TeamsByEventQuery(eventCode, year)
	if eventTeams == nil {
		s.writeError(w, http.StatusNotFound, "event not found")
		return
	}

	teams := eventTeams.Teams
	if limit > 0 && limit < len(teams) {
		teams = teams[:limit]
	}

	response := EventTeamsResponse{
		Event: &EventWithTeams{
			EventResponse: toEventResponse(eventTeams.Event),
			Teams:         teams,
		},
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleEventRankings handles requests for the team rankings of a specific event. It expects the event code to be provided in the URL path and supports a 'limit' query parameter to limit the number of rankings returned. It returns the event details along with the list of team rankings in JSON format.
func (s *Server) handleEventRankings(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	rankings := query.EventTeamRankingQuery(eventCode, year)
	if rankings == nil {
		s.writeError(w, http.StatusNotFound, "event not found")
		return
	}

	// Convert to clean response format
	rankingList := make([]RankingResponse, 0, len(rankings.TeamRankings))
	for _, tr := range rankings.TeamRankings {
		rankingList = append(rankingList, RankingResponse{
			Team:           tr.Team,
			Year:           rankings.Event.Year,
			SortOrder1:     tr.Ranking.SortOrder1,
			SortOrder2:     tr.Ranking.SortOrder2,
			SortOrder3:     tr.Ranking.SortOrder3,
			SortOrder4:     tr.Ranking.SortOrder4,
			SortOrder5:     tr.Ranking.SortOrder5,
			SortOrder6:     tr.Ranking.SortOrder6,
			Wins:           tr.Ranking.Wins,
			Losses:         tr.Ranking.Losses,
			Ties:           tr.Ranking.Ties,
			Dq:             tr.Ranking.Dq,
			MatchesPlayed:  tr.Ranking.MatchesPlayed,
			MatchesCounted: tr.Ranking.MatchesCounted,
			HighMatchScore: tr.HighMatchScore,
		})
	}

	if limit > 0 && limit < len(rankingList) {
		rankingList = rankingList[:limit]
	}

	response := EventRankingsResponse{
		Event:    toEventResponse(rankings.Event),
		Rankings: rankingList,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleEventAwards handles requests for the awards given at a specific event. It expects the event code to be provided in the URL path and supports a 'limit' query parameter to limit the number of awards returned. It returns the event details along with the list of awards in JSON format.
func (s *Server) handleEventAwards(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	awards := query.AwardsByEventQuery(eventCode, year)
	if awards == nil {
		s.writeError(w, http.StatusNotFound, "event not found")
		return
	}

	// Convert to clean response format without event_id
	awardList := make([]AwardResponse, 0, len(awards.Awards))
	for _, ta := range awards.Awards {
		awardList = append(awardList, AwardResponse{
			Name: ta.Award.Name,
			Team: ta.Team,
		})
	}

	if limit > 0 && limit < len(awardList) {
		awardList = awardList[:limit]
	}

	response := EventAwardsResponse{
		Event: &EventWithAwards{
			EventResponse: toEventResponse(awards.Event),
			Awards:        awardList,
		},
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleEventAdvancement handles requests for the advancement details of a specific event. It expects the event code to be provided in the URL path and returns the event details along with the team advancements in JSON format.
func (s *Server) handleEventAdvancement(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	advancement := query.AdvancementReportQuery(eventCode, year)
	if advancement == nil || advancement.Event == nil {
		s.writeError(w, http.StatusNotFound, "event not found")
		return
	}

	response := EventAdvancementResponse{
		Event:            toEventResponse(advancement.Event),
		TeamAdvancements: advancement.TeamAdvancements,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleEventMatches handles requests for the matches of a specific event. It expects the event code to be provided in the URL path and supports an optional 'team' query parameter to filter matches by a specific team. It also supports a 'limit' query parameter to limit the number of matches returned. It returns the event details along with the list of matches (with alliance details if team filter is not applied) in JSON format.
func (s *Server) handleEventMatches(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	teamIDStr := r.URL.Query().Get("team")
	var matches interface{}
	var event *database.Event

	if teamIDStr != "" {
		teamID, err := strconv.Atoi(teamIDStr)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid team parameter: %s", teamIDStr))
			return
		}
		matchList := query.MatchesByEventAndTeamQuery(eventCode, teamID, year)
		if len(matchList) > 0 {
			event = matchList[0].Event
		}
		if limit > 0 && limit < len(matchList) {
			matchList = matchList[:limit]
		}
		// Convert to TeamMatchResultResponse list
		convertedMatches := make([]*TeamMatchResultResponse, 0, len(matchList))
		for _, m := range matchList {
			convertedMatches = append(convertedMatches, toTeamMatchResultResponse(m))
		}
		matches = convertedMatches
	} else {
		matchList := query.MatchesByEventQuery(eventCode, year)
		if len(matchList) > 0 {
			event = matchList[0].Event
		}
		if limit > 0 && limit < len(matchList) {
			matchList = matchList[:limit]
		}
		// Convert to MatchWithAlliancesResponse list
		convertedMatches := make([]*MatchWithAlliancesResponse, 0, len(matchList))
		for _, m := range matchList {
			convertedMatches = append(convertedMatches, toMatchWithAlliancesResponse(m.Match, m.RedAlliance, m.BlueAlliance))
		}
		matches = convertedMatches
	}

	if event == nil {
		s.writeError(w, http.StatusNotFound, "event not found")
		return
	}

	response := EventMatchesResponse{
		Event: &EventWithMatches{
			EventResponse: toEventResponse(event),
			Matches:       matches,
		},
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleTeamRankings handles requests for the overall team rankings for a specific season. It supports optional query parameters for region, country, and event code to filter the rankings. It also supports a 'limit' query parameter to limit the number of rankings returned. It returns a list of team performances in JSON format.
func (s *Server) handleTeamRankings(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	region := r.URL.Query().Get("region")
	country := r.URL.Query().Get("country")
	eventCode := r.URL.Query().Get("event")

	performances, err := query.TeamRankingsQuery(region, country, eventCode, year)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if limit > 0 && limit < len(performances) {
		performances = performances[:limit]
	}

	s.writeJSON(w, http.StatusOK, performances)
}

// handleTeamEventRankings handles requests for the team rankings at specific events for a specific season. It supports optional query parameters for region, country, and event code to filter the rankings. It also supports a 'limit' query parameter to limit the number of rankings returned. It returns a list of team performances at events in JSON format.
func (s *Server) handleTeamEventRankings(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	region := r.URL.Query().Get("region")
	country := r.URL.Query().Get("country")
	eventCode := r.URL.Query().Get("event")

	performances, err := query.TeamEventRankingsQuery(region, country, eventCode, year)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert to EventPerformanceResponse (without event_id, with year)
	responses := make([]EventPerformanceResponse, 0, len(performances))
	for _, p := range performances {
		responses = append(responses, EventPerformanceResponse{
			TeamID:    p.TeamID,
			TeamName:  p.TeamName,
			Region:    p.Region,
			Year:      year,
			EventCode: p.EventCode,
			EventName: p.EventName,
			OPR:       p.OPR,
			NpOPR:     p.NpOPR,
			CCWM:      p.CCWM,
			DPR:       p.DPR,
			NpDPR:     p.NpDPR,
			NpAVG:     p.NpAVG,
			Matches:   p.Matches,
		})
	}

	if limit > 0 && limit < len(responses) {
		responses = responses[:limit]
	}

	s.writeJSON(w, http.StatusOK, responses)
}

// handleRegions handles requests for region-specific resources. It expects the region code to be provided in the URL path and delegates to specific handlers for different region resources such as advancement based on the second part of the URL path.
func (s *Server) handleRegions(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	if len(parts) < 1 {
		s.writeError(w, http.StatusBadRequest, "region code is required")
		return
	}

	regionCode := parts[0]

	if len(parts) < 2 {
		s.writeError(w, http.StatusBadRequest, "region resource type is required")
		return
	}

	resource := parts[1]

	switch resource {
	case "advancement":
		s.handleRegionAdvancement(w, r, year, regionCode)
	default:
		s.writeError(w, http.StatusNotFound, fmt.Sprintf("unknown region resource: %s", resource))
	}
}

// handleRegionAdvancement handles requests for the advancement summary of a specific region and season. It expects the region code to be provided in the URL path and returns the advancement summary for that region and season in JSON format.
func (s *Server) handleRegionAdvancement(w http.ResponseWriter, r *http.Request, year int, regionCode string) {
	advancement := query.RegionAdvancementQuery(regionCode, year)
	s.writeJSON(w, http.StatusOK, advancement)
}

// handleAllAdvancement handles requests for the advancement summary of all regions for a specific season. It supports an optional 'region' query parameter to filter the summary by a specific region. It returns the advancement summary for the specified region (or all regions if no region is specified) and season in JSON format.
func (s *Server) handleAllAdvancement(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "ALL"
	}
	advancement := query.EventAdvancementSummaryQuery(region, year)
	response := toEventAdvancementSummaryResponse(advancement)
	s.writeJSON(w, http.StatusOK, response)
}

// writeJSON is a helper function to write a JSON response with the given status code and data. It sets the appropriate content type header and encodes the data as JSON. If encoding fails, it logs an error.
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to encode JSON response", "error", err)
	}
}

// writeError is a helper function to write an error response in JSON format with the given status code and error message. It uses the writeJSON helper to send a standardized error response.
func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}
