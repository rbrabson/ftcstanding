package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

// sqldb wraps a sql.DB and provides prepared statements for common operations.
type sqldb struct {
	ctx   context.Context
	sqldb *sql.DB
	stmts map[string]*sql.Stmt
}

// InitDB initializes the database connection.
func InitSQLDB() (*sqldb, error) {
	godotenv.Load()
	dsn := os.Getenv("DATA_SOURCE_NAME")
	if dsn == "" {
		return nil, errors.New("DATA_SOURCE_NAME environment variable not set")
	}

	ctx := context.Background()
	var err error
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	// Set database connection pool settings
	sqlDB.SetConnMaxLifetime(time.Minute * 3) // Make it less than 5 minutes to avoid timeouts
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10) // Make it the same as MaxOpenConns

	db := &sqldb{
		ctx:   ctx,
		sqldb: sqlDB,
		stmts: make(map[string]*sql.Stmt),
	}
	db.initStatements()

	return db, nil

}

// CloseDB closes all prepared statements.
func (db *sqldb) Close() {
	for _, stmt := range db.stmts {
		stmt.Close()
	}
	db.stmts = make(map[string]*sql.Stmt)
}

// InitStatements initializes all prepared statements for the dbmodel package.
func (db *sqldb) initStatements() error {
	if err := db.initEventStatements(); err != nil {
		return err
	}
	if err := db.initAwardStatements(); err != nil {
		return err
	}
	if err := db.initMatchStatements(); err != nil {
		return err
	}
	if err := db.initTeamStatements(); err != nil {
		return err
	}

	return nil
}

// InitAwardStatements prepares all SQL statements for award operations.
func (db *sqldb) initAwardStatements() error {
	queries := map[string]string{
		"getAward":     "SELECT award_id, name, description, for_person FROM awards WHERE award_id = ?",
		"getAllAwards": "SELECT award_id, name, description, for_person FROM awards",
		"saveAward":    "INSERT INTO awards (award_id, name, description, for_person) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description), for_person = VALUES(for_person)",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// InitEventStatements prepares all SQL statements for event operations.
func (db *sqldb) initEventStatements() error {
	queries := map[string]string{
		"getEvent":                "SELECT event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end FROM events WHERE event_id = ?",
		"saveEvent":               "INSERT INTO events (event_id, event_code, year, name, type, division_code, region_code, league_code, venue, address, city, state_prov, country, timezone, date_start, date_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE event_code = VALUES(event_code), year = VALUES(year), name = VALUES(name), type = VALUES(type), division_code = VALUES(division_code), region_code = VALUES(region_code), league_code = VALUES(league_code), venue = VALUES(venue), address = VALUES(address), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), timezone = VALUES(timezone), date_start = VALUES(date_start), date_end = VALUES(date_end)",
		"getEventAwards":          "SELECT event_id, team_id, award_id FROM event_awards WHERE event_id = ?",
		"saveEventAward":          "INSERT INTO event_awards (event_id, team_id, award_id) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE event_id = event_id",
		"getTeamAwardsByEvent":    "SELECT event_id, team_id, award_id FROM event_awards WHERE event_id = ? AND team_id = ?",
		"getAllTeamAwards":        "SELECT event_id, team_id, award_id FROM event_awards WHERE team_id = ? ORDER BY event_id",
		"getEventRankings":        "SELECT event_id, team_id, rank, sort_order1, sort_order2, sort_order3, sort_order4, sort_order5, sort_order6, wins, losses, ties, dq, matches_played, matches_counted FROM event_rankings WHERE event_id = ?",
		"saveEventRanking":        "INSERT INTO event_rankings (event_id, team_id, rank, sort_order1, sort_order2, sort_order3, sort_order4, sort_order5, sort_order6, wins, losses, ties, dq, matches_played, matches_counted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE rank = VALUES(rank), sort_order1 = VALUES(sort_order1), sort_order2 = VALUES(sort_order2), sort_order3 = VALUES(sort_order3), sort_order4 = VALUES(sort_order4), sort_order5 = VALUES(sort_order5), sort_order6 = VALUES(sort_order6), wins = VALUES(wins), losses = VALUES(losses), ties = VALUES(ties), dq = VALUES(dq), matches_played = VALUES(matches_played), matches_counted = VALUES(matches_counted)",
		"getEventAdvancements":    "SELECT event_id, team_id FROM event_advancements WHERE event_id = ?",
		"saveEventAdvancement":    "INSERT INTO event_advancements (event_id, team_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE event_id = event_id",
		"getAllAdvancements":      "SELECT event_id, team_id FROM event_advancements ORDER BY event_id, team_id",
		"getRegionCodes":          "SELECT DISTINCT region_code FROM events WHERE region_code IS NOT NULL AND region_code != '' ORDER BY region_code",
		"getEventCodesByRegion":   "SELECT DISTINCT event_code FROM events WHERE region_code = ? ORDER BY event_code",
		"getAdvancementsByRegion": "SELECT ea.event_id, ea.team_id FROM event_advancements ea INNER JOIN events e ON ea.event_id = e.event_id WHERE e.region_code = ? ORDER BY ea.event_id, ea.team_id",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// InitMatchStatements prepares all SQL statements for match operations.
func (db *sqldb) initMatchStatements() error {
	queries := map[string]string{
		"getMatch":               "SELECT match_id, event_id, match_number, actual_start_time, description, tournament_level FROM matches WHERE match_id = ?",
		"getAllMatches":          "SELECT match_id, event_id, match_number, actual_start_time, description, tournament_level FROM matches",
		"getMatchesByEvent":      "SELECT match_id, event_id, match_number, actual_start_time, description, tournament_level FROM matches WHERE event_id = ? ORDER BY match_number",
		"saveMatch":              "INSERT INTO matches (match_id, event_id, match_number, actual_start_time, description, tournament_level) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE event_id = VALUES(event_id), match_number = VALUES(match_number), actual_start_time = VALUES(actual_start_time), description = VALUES(description), tournament_level = VALUES(tournament_level)",
		"getMatchAllianceScore":  "SELECT match_id, alliance, auto_points, teleop_points, foul_points_committed, pre_foul_total, total_points, major_fouls, minor_fouls FROM match_alliance_scores WHERE match_id = ? AND alliance = ?",
		"saveMatchAllianceScore": "INSERT INTO match_alliance_scores (match_id, alliance, auto_points, teleop_points, foul_points_committed, pre_foul_total, total_points, major_fouls, minor_fouls) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE auto_points = VALUES(auto_points), teleop_points = VALUES(teleop_points), foul_points_committed = VALUES(foul_points_committed), pre_foul_total = VALUES(pre_foul_total), total_points = VALUES(total_points), major_fouls = VALUES(major_fouls), minor_fouls = VALUES(minor_fouls)",
		"getMatchTeams":          "SELECT match_id, team_id, alliance, dq, on_field FROM match_teams WHERE match_id = ?",
		"getTeamsByEvent":        "SELECT DISTINCT mt.team_id FROM match_teams mt INNER JOIN matches m ON mt.match_id = m.match_id WHERE m.event_id = ? ORDER BY mt.team_id",
		"saveMatchTeam":          "INSERT INTO match_teams (match_id, team_id, alliance, dq, on_field) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE alliance = VALUES(alliance), dq = VALUES(dq), on_field = VALUES(on_field)",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// InitTeamStatements prepares all SQL statements for team operations.
func (db *sqldb) initTeamStatements() error {
	queries := map[string]string{
		"getTeam":          "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE team_id = ?",
		"getAllTeams":      "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams",
		"getTeamsByRegion": "SELECT team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name FROM teams WHERE home_region = ? ORDER BY team_id",
		"saveTeam":         "INSERT INTO teams (team_id, name, full_name, city, state_prov, country, website, rookie_year, home_region, robot_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), full_name = VALUES(full_name), city = VALUES(city), state_prov = VALUES(state_prov), country = VALUES(country), website = VALUES(website), rookie_year = VALUES(rookie_year), home_region = VALUES(home_region), robot_name = VALUES(robot_name)",
	}

	for name, query := range queries {
		if err := db.prepareStatement(name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}
	return nil
}

// PrepareStatement prepares and caches a SQL statement.
func (db *sqldb) prepareStatement(name, query string) error {
	stmt, err := db.sqldb.Prepare(query)
	if err != nil {
		return err
	}
	db.stmts[name] = stmt
	return nil
}

// GetStatement retrieves a prepared statement by name.
func (db *sqldb) getStatement(name string) *sql.Stmt {
	return db.stmts[name]
}
