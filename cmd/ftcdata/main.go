package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rbrabson/ftcstanding/request"
	"github.com/spf13/cobra"
)

var (
	db          database.DB
	allFlag     bool
	regionFlag  string
	eventFlag   string
	seasonFlag  string
	refreshFlag bool
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

var rootCmd = &cobra.Command{
	Use:   "ftcdata",
	Short: "FTC Standing data synchronization tool",
	Long:  `A tool to synchronize FTC (FIRST Tech Challenge) standing data including teams, events, matches, awards, and rankings.`,
	Example: `  # Sync all data for the season
  ftcdata --season 2025 --all

  # Sync data for a specific region
  ftcdata --season 2025 --region USNC

  # Sync data for a specific event
  ftcdata --season 2025 --event USNCRAQ

  # Force refresh all data
  ftcdata --season 2025 --all --refresh`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no action flags are specified, show help
		if !allFlag && eventFlag == "" && regionFlag == "" {
			return cmd.Help()
		}

		// Determine season
		season := seasonFlag
		if season == "" {
			season = os.Getenv("FTC_SEASON")
			if season == "" {
				return fmt.Errorf("season not specified. Use --season flag or set FTC_SEASON environment variable")
			}
		}

		var err error
		db, err = database.Init(season)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		request.Init(db)
		query.Init(db)

		// Handle different modes based on flags
		switch {
		case eventFlag != "":
			// Process single event
			processEvent(season, eventFlag)
		case regionFlag != "":
			// Process region
			processRegion(season, regionFlag, refreshFlag)
		case allFlag:
			// Process all data
			request.RequestAndSaveAll(season, refreshFlag)
		}

		return nil
	},
}

func init() {
	// Load environment variables
	godotenv.Load()
	setLogLevelFromEnv()

	// Define flags
	rootCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Sync all data for the season")
	rootCmd.Flags().StringVarP(&regionFlag, "region", "r", "", "Region code to filter events (e.g., USCHS)")
	rootCmd.Flags().StringVarP(&eventFlag, "event", "e", "", "Event code to process (e.g., USNCCOQ)")
	rootCmd.Flags().StringVarP(&seasonFlag, "season", "s", "", "Season year (defaults to FTC_SEASON environment variable)")
	rootCmd.Flags().BoolVar(&refreshFlag, "refresh", false, "Force refresh of all data")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// processEvent processes a single event
func processEvent(season, eventCode string) {
	slog.Info("Processing single event", "eventCode", eventCode, "season", season)

	// Get the event
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events, err := db.GetAllEvents(filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load events: %v\n", err)
		os.Exit(1)
	}

	var event *database.Event
	for _, e := range events {
		if e.EventCode == eventCode {
			event = e
			break
		}
	}

	if event == nil {
		fmt.Fprintf(os.Stderr, "Error: Event %s not found\n", eventCode)
		os.Exit(1)
	}

	// Process event details
	request.RequestAndSaveEventAwards(event)
	request.RequestAndSaveEventRankings(event)
	request.RequestAndSaveEventAdvancements(event)
	request.RequestAndSaveMatches(event)
	request.RequestAndSaveTeamsInEvent(event)

	slog.Info("Finished processing event", "eventCode", eventCode)
}

// processRegion processes all events in a region
func processRegion(season, regionCode string, refresh bool) {
	slog.Info("Processing region", "regionCode", regionCode, "season", season)

	// Get or refresh teams and awards
	teams, err := db.GetAllTeams()
	if err != nil {
		slog.Warn("failed to load teams", "error", err)
	}
	if refresh || len(teams) == 0 {
		teams = request.RequestAndSaveTeams(season)
	}

	awards, err := db.GetAllAwards()
	if err != nil {
		slog.Warn("failed to load awards", "error", err)
	}
	if refresh || len(awards) == 0 {
		awards = request.RequestAndSaveAwards(season)
	}

	// Get events for the region
	filter := database.EventFilter{
		RegionCodes: []string{regionCode},
	}
	events, err := db.GetAllEvents(filter)
	if err != nil {
		slog.Warn("failed to load region events", "regionCode", regionCode, "error", err)
	}

	if refresh || len(events) == 0 {
		// Refresh all events and filter
		allEvents := request.RequestAndSaveEvents(season)
		events = nil
		for _, e := range allEvents {
			if e.RegionCode == regionCode {
				events = append(events, e)
			}
		}
	}

	slog.Info("Found events in region", "regionCode", regionCode, "eventCount", len(events))

	// If not refresh, filter events to only those in the past 24 hours
	filteredEvents := events
	if !refresh {
		now := time.Now()
		var recentEvents []*database.Event
		for _, event := range events {
			if event.DateStart.Before(now) && event.DateStart.After(now.Add(-24*time.Hour)) {
				recentEvents = append(recentEvents, event)
			}
		}
		filteredEvents = recentEvents
	}

	for i, event := range filteredEvents {
		slog.Info("Processing event", "eventNumber", i+1, "totalEvents", len(filteredEvents), "event", event.EventCode)

		request.RequestAndSaveEventAwards(event)
		request.RequestAndSaveEventRankings(event)
		request.RequestAndSaveEventAdvancements(event)
		request.RequestAndSaveMatches(event)
		request.RequestAndSaveTeamsInEvent(event)

		slog.Info("Finished processing event", "eventCode", event.EventCode)
	}

	slog.Info("Finished processing region", "regionCode", regionCode)
}
