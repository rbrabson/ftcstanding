package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx   = context.Background()
	DB    *sql.DB
	stmts = make(map[string]*sql.Stmt)
)

// InitDB initializes the database connection.
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	if err := DB.PingContext(ctx); err != nil {
		return err
	}
	// Set database connection pool settings
	DB.SetConnMaxLifetime(time.Minute * 3) // Make it less than 5 minutes to avoid timeouts
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	return nil
}

// PrepareStatement prepares and caches a SQL statement.
func PrepareStatement(name, query string) error {
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	stmts[name] = stmt
	return nil
}

// GetStatement retrieves a prepared statement by name.
func GetStatement(name string) *sql.Stmt {
	return stmts[name]
}

// CloseStatements closes all prepared statements.
func CloseStatements() {
	for _, stmt := range stmts {
		stmt.Close()
	}
	stmts = make(map[string]*sql.Stmt)
}
