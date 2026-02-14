package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// filedb implements a file-based database that stores data in JSON files.
// This implementation provides a lightweight alternative to SQL databases
// for development, testing, or deployments where a full database server
// is not available or desired.
//
// Data is stored in memory during runtime and persisted to JSON files
// in the specified data directory. Each entity type (awards, teams, events,
// matches, etc.) is stored in a separate JSON file for easy inspection and
// manual editing if needed.
//
// The filedb is thread-safe and uses table-level read-write locks to allow
// concurrent operations on different tables while ensuring exclusive access
// for writes within each table.
//
// Usage:
//
//	db, err := database.InitFileDB("./data")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	// Use db like any other DB interface implementation
//	team := db.GetTeam(12345)
type filedb struct {
	dataDir string

	fileStateMu sync.Mutex
	fileStates  map[string]fileState

	// Table-level locks for fine-grained concurrency control
	awardsMu            sync.RWMutex
	teamsMu             sync.RWMutex
	teamRankingsMu      sync.RWMutex
	eventsMu            sync.RWMutex
	eventAwardsMu       sync.RWMutex
	eventRankingsMu     sync.RWMutex
	eventAdvancementsMu sync.RWMutex
	eventTeamsMu        sync.RWMutex
	matchesMu           sync.RWMutex
	matchScoresMu       sync.RWMutex
	matchTeamsMu        sync.RWMutex

	awards            map[int]*Award
	teams             map[int]*Team
	teamRankings      map[string]map[int]*TeamRanking // eventID -> teamID -> ranking
	events            map[string]*Event
	eventAwards       map[string][]*EventAward       // keyed by eventID
	eventRankings     map[string][]*EventRanking     // keyed by eventID
	eventAdvancements map[string][]*EventAdvancement // keyed by eventID
	eventTeams        map[string][]*EventTeam        // keyed by eventID
	matches           map[string]*Match
	matchScores       map[string]map[string]*MatchAllianceScore // matchID -> alliance -> score
	matchTeams        map[string][]*MatchTeam                   // keyed by matchID
}

type fileState struct {
	exists  bool
	modTime time.Time
	size    int64
}

// InitFileDB initializes a file-based database.
// season is an optional parameter. If provided, it will be used to construct the data directory path.
// If not provided, the FTC_SEASON environment variable will be used.
//
// The function creates the data directory if it doesn't exist and loads
// any existing data from JSON files in that directory. If the directory
// is empty or files don't exist, the database starts with empty datasets.
func initFileDB(season ...string) (*filedb, error) {
	godotenv.Load()
	baseDir := os.Getenv("FILEDB_DATA_DIR")
	if baseDir == "" {
		return nil, errors.New("FILEDB_DATA_DIR environment variable not set")
	}

	var year string
	if len(season) > 0 && season[0] != "" {
		year = season[0]
	} else {
		year = os.Getenv("FTC_SEASON")
		if year == "" {
			return nil, errors.New("FTC_SEASON environment variable not set")
		}
	}
	dataDir := filepath.Join(baseDir, year)

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db := &filedb{
		dataDir:           dataDir,
		fileStates:        make(map[string]fileState),
		awards:            make(map[int]*Award),
		teams:             make(map[int]*Team),
		teamRankings:      make(map[string]map[int]*TeamRanking),
		events:            make(map[string]*Event),
		eventAwards:       make(map[string][]*EventAward),
		eventRankings:     make(map[string][]*EventRanking),
		eventAdvancements: make(map[string][]*EventAdvancement),
		eventTeams:        make(map[string][]*EventTeam),
		matches:           make(map[string]*Match),
		matchScores:       make(map[string]map[string]*MatchAllianceScore),
		matchTeams:        make(map[string][]*MatchTeam),
	}

	// Load existing data
	if err := db.loadAll(); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return db, nil
}

func (db *filedb) refreshAllIfChanged() error {
	if err := db.refreshAwardsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshTeamsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshTeamRankingsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshEventsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshEventAwardsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshEventRankingsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshEventAdvancementsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshEventTeamsIfChanged(); err != nil {
		return err
	}
	if err := db.refreshMatchesIfChanged(); err != nil {
		return err
	}
	if err := db.refreshMatchScoresIfChanged(); err != nil {
		return err
	}
	if err := db.refreshMatchTeamsIfChanged(); err != nil {
		return err
	}

	return nil
}

// Close implements the DB interface. For file-based DB, this saves all data.
func (db *filedb) Close() {
	db.saveAll()
}

// loadAll loads all data from JSON files.
func (db *filedb) loadAll() error {
	// Lock all tables for loading
	db.awardsMu.Lock()
	defer db.awardsMu.Unlock()
	db.teamsMu.Lock()
	defer db.teamsMu.Unlock()
	db.teamRankingsMu.Lock()
	defer db.teamRankingsMu.Unlock()
	db.eventsMu.Lock()
	defer db.eventsMu.Unlock()
	db.eventAwardsMu.Lock()
	defer db.eventAwardsMu.Unlock()
	db.eventRankingsMu.Lock()
	defer db.eventRankingsMu.Unlock()
	db.eventAdvancementsMu.Lock()
	defer db.eventAdvancementsMu.Unlock()
	db.eventTeamsMu.Lock()
	defer db.eventTeamsMu.Unlock()
	db.matchesMu.Lock()
	defer db.matchesMu.Unlock()
	db.matchScoresMu.Lock()
	defer db.matchScoresMu.Unlock()
	db.matchTeamsMu.Lock()
	defer db.matchTeamsMu.Unlock()

	// Load awards
	if err := db.loadJSONFile("awards.json", &db.awards); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load teams
	if err := db.loadJSONFile("teams.json", &db.teams); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load team rankings
	if err := db.loadJSONFile("team_rankings.json", &db.teamRankings); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load events
	if err := db.loadJSONFile("events.json", &db.events); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load event awards
	if err := db.loadJSONFile("event_awards.json", &db.eventAwards); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load event rankings
	if err := db.loadJSONFile("event_rankings.json", &db.eventRankings); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load event advancements
	if err := db.loadJSONFile("event_advancements.json", &db.eventAdvancements); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load event teams
	if err := db.loadJSONFile("event_teams.json", &db.eventTeams); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load matches
	if err := db.loadJSONFile("matches.json", &db.matches); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load match scores
	if err := db.loadJSONFile("match_scores.json", &db.matchScores); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load match teams
	if err := db.loadJSONFile("match_teams.json", &db.matchTeams); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// saveAll saves all data to JSON files.
func (db *filedb) saveAll() error {
	if err := db.refreshAllIfChanged(); err != nil {
		return err
	}

	// Lock all tables for saving (read locks since we're only reading the data structures to save)
	db.awardsMu.RLock()
	defer db.awardsMu.RUnlock()
	db.teamsMu.RLock()
	defer db.teamsMu.RUnlock()
	db.teamRankingsMu.RLock()
	defer db.teamRankingsMu.RUnlock()
	db.eventsMu.RLock()
	defer db.eventsMu.RUnlock()
	db.eventAwardsMu.RLock()
	defer db.eventAwardsMu.RUnlock()
	db.eventRankingsMu.RLock()
	defer db.eventRankingsMu.RUnlock()
	db.eventAdvancementsMu.RLock()
	defer db.eventAdvancementsMu.RUnlock()
	db.eventTeamsMu.RLock()
	defer db.eventTeamsMu.RUnlock()
	db.matchesMu.RLock()
	defer db.matchesMu.RUnlock()
	db.matchScoresMu.RLock()
	defer db.matchScoresMu.RUnlock()
	db.matchTeamsMu.RLock()
	defer db.matchTeamsMu.RUnlock()

	if err := db.saveJSONFile("awards.json", db.awards); err != nil {
		return err
	}

	if err := db.saveJSONFile("teams.json", db.teams); err != nil {
		return err
	}

	if err := db.saveJSONFile("team_rankings.json", db.teamRankings); err != nil {
		return err
	}

	if err := db.saveJSONFile("events.json", db.events); err != nil {
		return err
	}

	if err := db.saveJSONFile("event_awards.json", db.eventAwards); err != nil {
		return err
	}

	if err := db.saveJSONFile("event_rankings.json", db.eventRankings); err != nil {
		return err
	}

	if err := db.saveJSONFile("event_advancements.json", db.eventAdvancements); err != nil {
		return err
	}

	if err := db.saveJSONFile("event_teams.json", db.eventTeams); err != nil {
		return err
	}

	if err := db.saveJSONFile("matches.json", db.matches); err != nil {
		return err
	}

	if err := db.saveJSONFile("match_scores.json", db.matchScores); err != nil {
		return err
	}

	if err := db.saveJSONFile("match_teams.json", db.matchTeams); err != nil {
		return err
	}

	return nil
}

func (db *filedb) refreshAwardsIfChanged() error {
	return db.refreshJSONFileIfChanged("awards.json", &db.awardsMu, &db.awards)
}

func (db *filedb) refreshTeamsIfChanged() error {
	return db.refreshJSONFileIfChanged("teams.json", &db.teamsMu, &db.teams)
}

func (db *filedb) refreshTeamRankingsIfChanged() error {
	return db.refreshJSONFileIfChanged("team_rankings.json", &db.teamRankingsMu, &db.teamRankings)
}

func (db *filedb) refreshEventsIfChanged() error {
	return db.refreshJSONFileIfChanged("events.json", &db.eventsMu, &db.events)
}

func (db *filedb) refreshEventAwardsIfChanged() error {
	return db.refreshJSONFileIfChanged("event_awards.json", &db.eventAwardsMu, &db.eventAwards)
}

func (db *filedb) refreshEventRankingsIfChanged() error {
	return db.refreshJSONFileIfChanged("event_rankings.json", &db.eventRankingsMu, &db.eventRankings)
}

func (db *filedb) refreshEventAdvancementsIfChanged() error {
	return db.refreshJSONFileIfChanged("event_advancements.json", &db.eventAdvancementsMu, &db.eventAdvancements)
}

func (db *filedb) refreshEventTeamsIfChanged() error {
	return db.refreshJSONFileIfChanged("event_teams.json", &db.eventTeamsMu, &db.eventTeams)
}

func (db *filedb) refreshMatchesIfChanged() error {
	return db.refreshJSONFileIfChanged("matches.json", &db.matchesMu, &db.matches)
}

func (db *filedb) refreshMatchScoresIfChanged() error {
	return db.refreshJSONFileIfChanged("match_scores.json", &db.matchScoresMu, &db.matchScores)
}

func (db *filedb) refreshMatchTeamsIfChanged() error {
	return db.refreshJSONFileIfChanged("match_teams.json", &db.matchTeamsMu, &db.matchTeams)
}

func (db *filedb) refreshJSONFileIfChanged(filename string, mu *sync.RWMutex, target interface{}) error {
	changed, err := db.hasFileChanged(filename)
	if err != nil || !changed {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	changed, err = db.hasFileChanged(filename)
	if err != nil || !changed {
		return err
	}

	if err := db.loadJSONFile(filename, target); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := resetTarget(target); err != nil {
			return err
		}
	}

	return nil
}

func (db *filedb) hasFileChanged(filename string) (bool, error) {
	known, ok := db.getKnownFileState(filename)
	if !ok {
		return true, nil
	}

	current, err := db.currentFileState(filename)
	if err != nil {
		return false, err
	}

	return !sameFileState(current, known), nil
}

func (db *filedb) currentFileState(filename string) (fileState, error) {
	path := filepath.Join(db.dataDir, filename)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fileState{exists: false}, nil
		}
		return fileState{}, err
	}

	return fileState{
		exists:  true,
		modTime: info.ModTime(),
		size:    info.Size(),
	}, nil
}

func (db *filedb) getKnownFileState(filename string) (fileState, bool) {
	db.fileStateMu.Lock()
	defer db.fileStateMu.Unlock()

	state, ok := db.fileStates[filename]
	return state, ok
}

func (db *filedb) setKnownFileState(filename string, state fileState) {
	db.fileStateMu.Lock()
	defer db.fileStateMu.Unlock()

	db.fileStates[filename] = state
}

func sameFileState(a, b fileState) bool {
	if a.exists != b.exists {
		return false
	}
	if !a.exists {
		return true
	}

	return a.size == b.size && a.modTime.Equal(b.modTime)
}

func resetTarget(target interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	elem := v.Elem()
	switch elem.Kind() {
	case reflect.Map:
		elem.Set(reflect.MakeMap(elem.Type()))
	case reflect.Slice:
		elem.Set(reflect.MakeSlice(elem.Type(), 0, 0))
	default:
		elem.Set(reflect.Zero(elem.Type()))
	}

	return nil
}

// loadJSONFile loads data from a JSON file.
func (db *filedb) loadJSONFile(filename string, v interface{}) error {
	path := filepath.Join(db.dataDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			db.setKnownFileState(filename, fileState{exists: false})
		}
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	state, err := db.currentFileState(filename)
	if err != nil {
		return err
	}
	db.setKnownFileState(filename, state)

	return nil
}

// saveJSONFile saves data to a JSON file.
func (db *filedb) saveJSONFile(filename string, v interface{}) error {
	path := filepath.Join(db.dataDir, filename)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	state, err := db.currentFileState(filename)
	if err != nil {
		return err
	}
	db.setKnownFileState(filename, state)

	return nil
}
