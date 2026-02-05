package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
// The filedb is thread-safe and uses read-write locks to allow concurrent
// reads while ensuring exclusive access for writes.
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
	mu      sync.RWMutex

	awards            map[int]*Award
	teams             map[int]*Team
	events            map[string]*Event
	eventAwards       map[string][]*EventAward       // keyed by eventID
	eventRankings     map[string][]*EventRanking     // keyed by eventID
	eventAdvancements map[string][]*EventAdvancement // keyed by eventID
	matches           map[string]*Match
	matchScores       map[string]map[string]*MatchAllianceScore // matchID -> alliance -> score
	matchTeams        map[string][]*MatchTeam                   // keyed by matchID
}

// InitFileDB initializes a file-based database.
// dataDir is the directory where data files will be stored.
// If dataDir is empty, it defaults to "./data"
//
// The function creates the data directory if it doesn't exist and loads
// any existing data from JSON files in that directory. If the directory
// is empty or files don't exist, the database starts with empty datasets.
func initFileDB() (*filedb, error) {
	godotenv.Load()
	dataDir := os.Getenv("FILEDB_DATA_DIR")
	if dataDir == "" {
		return nil, errors.New("FILEDB_DATA_DIR environment variable not set")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db := &filedb{
		dataDir:           dataDir,
		awards:            make(map[int]*Award),
		teams:             make(map[int]*Team),
		events:            make(map[string]*Event),
		eventAwards:       make(map[string][]*EventAward),
		eventRankings:     make(map[string][]*EventRanking),
		eventAdvancements: make(map[string][]*EventAdvancement),
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

// Close implements the DB interface. For file-based DB, this saves all data.
func (db *filedb) Close() {
	db.saveAll()
}

// loadAll loads all data from JSON files.
func (db *filedb) loadAll() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Load awards
	if err := db.loadJSONFile("awards.json", &db.awards); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Load teams
	if err := db.loadJSONFile("teams.json", &db.teams); err != nil && !os.IsNotExist(err) {
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
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.saveJSONFile("awards.json", db.awards); err != nil {
		return err
	}

	if err := db.saveJSONFile("teams.json", db.teams); err != nil {
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

// loadJSONFile loads data from a JSON file.
func (db *filedb) loadJSONFile(filename string, v interface{}) error {
	path := filepath.Join(db.dataDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// saveJSONFile saves data to a JSON file.
func (db *filedb) saveJSONFile(filename string, v interface{}) error {
	path := filepath.Join(db.dataDir, filename)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
