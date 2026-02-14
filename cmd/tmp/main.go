package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/request"
)

// setLogLevelFromEnv sets the log level from the LOG_LEVEL environment variable.
func setLogLevelFromEnv() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")

	var logLevel slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	return logLevel
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	godotenv.Load()
	setLogLevelFromEnv()

	// Initialize database
	db, err := database.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Initialize request package
	request.Init(db)

	// Get all events
	events, err := db.GetAllEvents()
	if err != nil {
		return fmt.Errorf("failed to load events: %w", err)
	}
	slog.Info("Processing all events", "totalEvents", len(events))

	// Process each event
	for i, event := range events {
		slog.Info("Processing event", "eventNumber", i+1, "totalEvents", len(events), "event", event.EventCode)

		if err := request.RequestAndSaveTeamRankings(event); err != nil {
			slog.Error("Failed to process event", "event", event.EventCode, "error", err)
			continue
		}

		slog.Info("Completed event", "event", event.EventCode)
	}

	slog.Info("Finished processing all events", "totalEvents", len(events))
	return nil
}
