package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/lambda"
	"github.com/rbrabson/ftcstanding/performance"
)

func LoadMatchesFromDatabase(db database.DB, eventCode string, year int) ([]performance.Match, []int, error) {
	// Generate eventID from eventCode and year
	eventID := fmt.Sprintf("%s : %d", eventCode, year)

	// First, verify the event exists
	event := db.GetEvent(eventID)
	if event == nil {
		// Event not found with this exact ID, try searching for it
		events := db.GetAllEvents(database.EventFilter{EventCodes: []string{eventCode}})
		if len(events) == 0 {
			return nil, nil, fmt.Errorf("event %s not found in database", eventCode)
		}

		// Find event matching the year
		for _, e := range events {
			if e.Year == year || e.DateStart.Year() == year {
				event = e
				eventID = e.EventID
				break
			}
		}

		if event == nil {
			availableYears := []int{}
			for _, e := range events {
				availableYears = append(availableYears, e.Year)
			}
			return nil, nil, fmt.Errorf("event %s not found for year %d (available years: %v)", eventCode, year, availableYears)
		}
	}

	// Get matches for the event
	dbMatches := db.GetMatchesByEvent(eventID)
	if len(dbMatches) == 0 {
		return nil, nil, fmt.Errorf("no matches found for event %s (%s)", event.Name, eventID)
	}

	fmt.Printf("Found event: %s (%s)\n", event.Name, eventID)

	teamSet := map[int]struct{}{}
	matches := []performance.Match{}

	for _, dbMatch := range dbMatches {
		// Get alliance scores
		redScore := db.GetMatchAllianceScore(dbMatch.MatchID, database.AllianceRed)
		blueScore := db.GetMatchAllianceScore(dbMatch.MatchID, database.AllianceBlue)

		if redScore == nil || blueScore == nil {
			continue // Skip matches without complete score data
		}

		// Get teams in the match
		matchTeams := db.GetMatchTeams(dbMatch.MatchID)

		var redTeams []int
		var blueTeams []int

		for _, mt := range matchTeams {
			if !mt.OnField || mt.Dq {
				continue // Skip surrogates and DQ'd teams
			}

			if mt.Alliance == database.AllianceRed {
				redTeams = append(redTeams, mt.TeamID)
			} else {
				blueTeams = append(blueTeams, mt.TeamID)
			}

			teamSet[mt.TeamID] = struct{}{}
		}

		// Only include matches with teams on both alliances
		if len(redTeams) == 0 || len(blueTeams) == 0 {
			continue
		}

		matches = append(matches, performance.Match{
			RedTeams:      redTeams,
			BlueTeams:     blueTeams,
			RedScore:      float64(redScore.TotalPoints),
			BlueScore:     float64(blueScore.TotalPoints),
			RedPenalties:  float64(redScore.FoulPointsCommitted),
			BluePenalties: float64(blueScore.FoulPointsCommitted),
		})
	}

	teams := []int{}
	for t := range teamSet {
		teams = append(teams, t)
	}
	sort.Ints(teams)

	return matches, teams, nil
}

func main() {
	godotenv.Load()

	// Parse command-line arguments
	eventCode := flag.String("event", "", "Event code (required)")
	year := flag.Int("year", 0, "Year (defaults to FTC_SEASON environment variable)")
	flag.Parse()

	if *eventCode == "" {
		fmt.Fprintln(os.Stderr, "Error: -event flag is required")
		fmt.Fprintln(os.Stderr, "Usage: test -event <eventCode> [-year <year>]")
		os.Exit(1)
	}

	// Get year from environment if not specified
	if *year == 0 {
		season := os.Getenv("FTC_SEASON")
		if season == "" {
			fmt.Fprintln(os.Stderr, "Error: FTC_SEASON environment variable not set and -year not provided")
			os.Exit(1)
		}
		var err error
		*year, err = strconv.Atoi(season)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid FTC_SEASON value: %s\n", season)
			os.Exit(1)
		}
	}

	// Initialize database
	db, err := database.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Load matches from database
	matches, teams, err := LoadMatchesFromDatabase(db, *eventCode, *year)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Printf("Loaded %d matches with %d teams from event %s (%d)\n\n", len(matches), len(teams), *eventCode, *year)

	lambdaValue := lambda.GetLambda(len(matches))

	calculator := performance.Calculator{
		Matches: matches,
		Teams:   teams,
	}

	opr := calculator.CalculateOPR()
	npopr := calculator.CalculateNpOPR()
	ccwm := calculator.CalculateCCWM()
	dpr := calculator.CalculateDPR()
	npdpr := calculator.CalculateNpDPR()

	fmt.Println()
	fmt.Println("Without Regularization (Lambda = 0):")
	table := tablewriter.NewTable(os.Stdout)
	table.Header([]string{"Team", "OPR", "npOPR", "CCWM", "DPR", "npDPR", "npAVG"})

	for _, t := range teams {
		npavg := calculator.CalculateNpAVG(matches, t)
		table.Append([]string{
			fmt.Sprintf("%d", t),
			fmt.Sprintf("%.2f", opr[t]),
			fmt.Sprintf("%.2f", npopr[t]),
			fmt.Sprintf("%.2f", ccwm[t]),
			fmt.Sprintf("%.2f", dpr[t]),
			fmt.Sprintf("%.2f", npdpr[t]),
			fmt.Sprintf("%.2f", npavg),
		})
	}
	table.Render()

	calculator = performance.Calculator{
		Matches: matches,
		Teams:   teams,
		Lambda:  lambdaValue,
	}
	// TODO: add all the values to a struct, and add the ability to sort them on any field captured.
	//       the teamID needs to be the key into the table, with the values that follow.
	//       to do this, add a "calculator.Calculate()" method to calculate all fields, and return a
	//       map with the team as the key, and the values as the team-specific output.
	opr = calculator.CalculateOPR()
	npopr = calculator.CalculateNpOPR()
	ccwm = calculator.CalculateCCWM()
	dpr = calculator.CalculateDPR()
	npdpr = calculator.CalculateNpDPR()

	fmt.Println()
	fmt.Printf("With Regularization (Lambda = %.6f):\n", lambdaValue)
	table = tablewriter.NewTable(os.Stdout)
	table.Header([]string{"Team", "OPR", "npOPR", "CCWM", "DPR", "npDPR", "npAVG"})

	for _, t := range teams {
		npavg := calculator.CalculateNpAVG(matches, t)
		table.Append([]string{
			fmt.Sprintf("%d", t),
			fmt.Sprintf("%.2f", opr[t]),
			fmt.Sprintf("%.2f", npopr[t]),
			fmt.Sprintf("%.2f", ccwm[t]),
			fmt.Sprintf("%.2f", dpr[t]),
			fmt.Sprintf("%.2f", npdpr[t]),
			fmt.Sprintf("%.2f", npavg),
		})
	}
	table.Render()
}
