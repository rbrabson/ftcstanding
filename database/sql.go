package database

import (
	"context"
	"database/sql"
	"errors"
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
