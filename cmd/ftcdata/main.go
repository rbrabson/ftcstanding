package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rbrabson/ftcstanding/request"
)

var db database.DB

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

// printUsage prints the usage information for the command
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "FTC Standing data synchronization tool")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintln(os.Stderr, "  # Sync all data for the season")
	fmt.Fprintln(os.Stderr, "  ftc -season 2024")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  # Sync data for a specific region")
	fmt.Fprintln(os.Stderr, "  ftc -season 2024 -region USCHS")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  # Sync data for a specific event")
	fmt.Fprintln(os.Stderr, "  ftc -season 2024 -event USNCCOQ")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  # Force refresh all data")
	fmt.Fprintln(os.Stderr, "  ftc -season 2024 -refresh")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Environment Variables:")
	fmt.Fprintln(os.Stderr, "  FTC_SEASON   - Default season year if -season not provided")
	fmt.Fprintln(os.Stderr, "  LOG_LEVEL    - Logging level (debug, info, warn, error)")
	fmt.Fprintln(os.Stderr, "")
}

func main() {
	godotenv.Load()

	setLogLevelFromEnv()

	// Set custom usage function
	flag.Usage = printUsage

	// Define CLI flags
	helpFlag := flag.Bool("help", false, "Show usage information")
	regionFlag := flag.String("region", "", "Region code to filter events (e.g., USCHS)")
	eventFlag := flag.String("event", "", "Event code to process (e.g., USNCCOQ)")
	seasonFlag := flag.String("season", "", "Season year (defaults to FTC_SEASON environment variable)")
	refreshFlag := flag.Bool("refresh", false, "Force refresh of all data")

	flag.Parse()

	// Show help if requested
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Determine season
	season := *seasonFlag
	if season == "" {
		season = os.Getenv("FTC_SEASON")
		if season == "" {
			fmt.Fprintln(os.Stderr, "Error: Season not specified. Use -season flag or set FTC_SEASON environment variable")
			os.Exit(1)
		}
	}

	var err error
	db, err = database.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	request.Init(db)
	query.Init(db)

	// Handle different modes based on flags
	switch {
	case *eventFlag != "":
		// Process single event
		processEvent(season, *eventFlag)
	case *regionFlag != "":
		// Process region
		processRegion(season, *regionFlag, *refreshFlag)
	default:
		// Default behavior: process all
		request.RequestAndSaveAll(season, *refreshFlag)
	}
}

// processEvent processes a single event
func processEvent(season, eventCode string) {
	slog.Info("Processing single event", "eventCode", eventCode, "season", season)

	// Get the event
	filter := database.EventFilter{
		EventCodes: []string{eventCode},
	}
	events := db.GetAllEvents(filter)

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
	teams := db.GetAllTeams()
	if refresh || len(teams) == 0 {
		teams = request.RequestAndSaveTeams(season)
	}

	awards := db.GetAllAwards()
	if refresh || len(awards) == 0 {
		awards = request.RequestAndSaveAwards(season)
	}

	// Get events for the region
	filter := database.EventFilter{
		RegionCodes: []string{regionCode},
	}
	events := db.GetAllEvents(filter)

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

	// Process each event in the region
	for i, event := range events {
		slog.Info("Processing event", "eventNumber", i+1, "totalEvents", len(events), "event", event.EventCode)

		request.RequestAndSaveEventAwards(event)
		request.RequestAndSaveEventRankings(event)
		request.RequestAndSaveEventAdvancements(event)
		request.RequestAndSaveMatches(event)
		request.RequestAndSaveTeamsInEvent(event)

		slog.Info("Finished processing event", "eventCode", event.EventCode)
	}

	slog.Info("Finished processing region", "regionCode", regionCode)
}
