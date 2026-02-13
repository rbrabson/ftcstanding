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

func NewServer(db database.DB) *Server {
	s := &Server{
		db:     db,
		mux:    http.NewServeMux(),
		logger: slog.Default(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/v1/", s.handleV1Routes)
	s.mux.HandleFunc("/health", s.handleHealth)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

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

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

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

	if limit > 0 && limit < len(eventTeams.Teams) {
		eventTeams.Teams = eventTeams.Teams[:limit]
	}

	s.writeJSON(w, http.StatusOK, eventTeams)
}

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

	if limit > 0 && limit < len(rankings.TeamRankings) {
		rankings.TeamRankings = rankings.TeamRankings[:limit]
	}

	s.writeJSON(w, http.StatusOK, rankings)
}

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

	if limit > 0 && limit < len(awards.Awards) {
		awards.Awards = awards.Awards[:limit]
	}

	s.writeJSON(w, http.StatusOK, awards)
}

func (s *Server) handleEventAdvancement(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	advancement := query.AdvancementReportQuery(eventCode, year)
	s.writeJSON(w, http.StatusOK, advancement)
}

func (s *Server) handleEventMatches(w http.ResponseWriter, r *http.Request, year int, eventCode string) {
	limit, err := s.parseLimit(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	teamIDStr := r.URL.Query().Get("team")
	if teamIDStr != "" {
		teamID, err := strconv.Atoi(teamIDStr)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid team parameter: %s", teamIDStr))
			return
		}
		matches := query.MatchesByEventAndTeamQuery(eventCode, teamID, year)
		if limit > 0 && limit < len(matches) {
			matches = matches[:limit]
		}
		s.writeJSON(w, http.StatusOK, matches)
	} else {
		matches := query.MatchesByEventQuery(eventCode, year)
		if limit > 0 && limit < len(matches) {
			matches = matches[:limit]
		}
		s.writeJSON(w, http.StatusOK, matches)
	}
}

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

	if limit > 0 && limit < len(performances) {
		performances = performances[:limit]
	}

	s.writeJSON(w, http.StatusOK, performances)
}

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

func (s *Server) handleRegionAdvancement(w http.ResponseWriter, r *http.Request, year int, regionCode string) {
	advancement := query.RegionAdvancementQuery(regionCode, year)
	s.writeJSON(w, http.StatusOK, advancement)
}

func (s *Server) handleAllAdvancement(w http.ResponseWriter, r *http.Request, year int, parts []string) {
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "ALL"
	}
	advancement := query.EventAdvancementSummaryQuery(region, year)
	s.writeJSON(w, http.StatusOK, advancement)
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to encode JSON response", "error", err)
	}
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}
