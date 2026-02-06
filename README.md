# ftcstanding

A Go application for managing FTC (FIRST Tech Challenge) competition standings, teams, matches, events, and awards.

## Project Structure

``` code
ftcstanding/
├── cmd/
│   └── ftc/
│       └── main.go           # Application entry point
├── database/
│   ├── db.go                 # Database interface definition
│   ├── sql.go                # SQL database connection and prepared statement management
│   ├── sql_*.go              # SQL implementation for each entity type
│   ├── filedb.go             # File-based database implementation
│   ├── filedb_*.go           # File-based implementation for each entity type
│   ├── award.go              # Award data model with SQL query constants
│   ├── event.go              # Event data model with SQL query constants
│   ├── match.go              # Match data model with SQL query constants
│   ├── team.go               # Team data model with SQL query constants
│   └── statements.go         # Statement initialization
├── Configfile                # Configuration for Makefile
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
├── Makefile                  # Build automation for multiple platforms
├── .env                      # Environment configuration (not in git)
├── LICENSE
└── README.md
```

## Features

- **Multiple Database Backends**:
  - SQL database (MySQL) with connection pooling
  - File-based database using JSON storage (for development/testing)
- **Flexible Filtering**: Query data with optional filters for teams, events, matches, and advancements
  - Filter teams by ID, country, or home region
  - Filter events by event code, region code, or country
  - Filter matches by event IDs
  - Filter advancements by country or region code
  - Combine multiple filter criteria with intuitive OR/AND logic
- **Prepared Statements**: All SQL operations use prepared statements for performance and security
- **SQL Query Constants**: All SQL queries are defined as package-level constants for maintainability
- **String Representations**: All data models implement the `fmt.Stringer` interface for easy debugging and logging
- **Thread-Safe Operations**: File-based database includes mutex protection for concurrent access
- **Data Models**:
  - Teams: Manage team information including name, location, and rookie year
  - Events: Track competition events with dates, locations, and details
  - Matches: Record match results, alliance scores, and team participation
  - Awards: Manage awards and track which teams received them at events

## Prerequisites

- Go 1.24.0 or later
- MySQL database server (for SQL backend) OR
- File system access (for file-based backend)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/rbrabson/ftcstanding.git
   cd ftcstanding
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Configure your database backend:

   **For SQL Database (MySQL):**

   Create a `.env` file in the project root with your database connection string:

   ``` ini
   DATA_SOURCE_NAME=user:password@tcp(localhost:3306)/dbname
   ```

   **For File-Based Database:**

   No configuration needed. Data will be stored in JSON files in the `./data` directory by default.

## Database Backends

### SQL Database (MySQL)

The SQL backend uses MySQL with connection pooling and prepared statements for optimal performance and security. All SQL queries are defined as constants for easy maintenance.

To initialize:

```go
db, err := database.InitSQLDB()
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### File-Based Database (OS File System)

The file-based database provides a lightweight alternative that stores data in JSON files. This is ideal for:

- Development and testing
- Deployments without database servers
- Small datasets
- Easy data inspection and manual editing

Features:

- Thread-safe with read-write locks
- Automatic persistence on each save operation
- Human-readable JSON format
- Separate files for each entity type

To initialize:

```go
db, err := database.InitFileDB("./data")
if err != nil {
    log.Fatal(err)
}
defer db.Close() // Ensures all data is persisted
```

Both implementations satisfy the `database.DB` interface, so they can be used interchangeably.

## Database Setup

### SQL Database

The application expects a MySQL database with the following tables:

- `teams` - Team information
- `events` - Competition events
- `matches` - Match information
- `match_alliance_scores` - Alliance scores for matches
- `match_teams` - Team participation in matches
- `event_awards` - Awards given at events
- `event_rankings` - Team rankings within events
- `event_advancements` - Teams advancing from events
- `awards` - Award definitions

### File-Based Database

No setup required. The database will automatically create the following JSON files in the data directory:

- `awards.json` - Award definitions
- `teams.json` - Team information
- `events.json` - Competition events
- `matches.json` - Match information
- `match_scores.json` - Alliance scores for matches
- `match_teams.json` - Team participation in matches
- `event_awards.json` - Awards given at events
- `event_rankings.json` - Team rankings within events
- `event_advancements.json` - Teams advancing from events

## Usage

After building (see Development section), run the appropriate binary for your platform:

```bash
# macOS ARM (Apple Silicon)
./bin/macos/arm64/rank

# macOS Intel
./bin/macos/amd64/rank

# Linux
./bin/linux/amd64/rank

# Windows
.\bin\windows\amd64\rank
```

The `database.DB` interface provides a consistent API for both SQL and file-based backends. All database operations are available through this interface.

### Using Filters

The database supports flexible filtering for querying data. Filters use optional variadic parameters:

```go
// Get all teams (no filter)
allTeams := db.GetAllTeams()

// Get teams from specific countries
usaCanadaTeams := db.GetAllTeams(database.TeamFilter{
    Countries: []string{"USA", "Canada"},
})

// Get specific teams by ID
selectedTeams := db.GetAllTeams(database.TeamFilter{
    TeamIDs: []int{12345, 67890},
})

// Combine filters (AND logic between fields)
californiaUSATeams := db.GetAllTeams(database.TeamFilter{
    Countries: []string{"USA"},
    HomeRegions: []string{"California"},
})

// Filter events by region
regionalEvents := db.GetAllEvents(database.EventFilter{
    RegionCodes: []string{"USCALA", "USTXHO"},
})

// Filter matches by event
eventMatches := db.GetAllMatches(database.MatchFilter{
    EventIDs: []string{"EVENT-123 : 2024", "EVENT-456 : 2024"},
})

// Filter advancements by country
usaAdvancements := db.GetAllAdvancements(database.AdvancementFilter{
    Countries: []string{"USA"},
})
```

### String Representations

All data models implement the `fmt.Stringer` interface for convenient logging and debugging:

```go
team := db.GetTeam(12345)
fmt.Println(team) // Output: Team{ID: 12345, Name: "Example Team", City: Boston, MA, Region: US-MA}

award := db.GetAward(1)
fmt.Println(award) // Output: Award{ID: 1, Name: "Inspire Award", Type: Team}
```

Or run directly without building:

```bash
go run ./cmd/ftc
```

## Database Operations

All database operations use prepared statements which are initialized at application startup:

### Teams

- `GetTeam(teamID)` - Retrieve a specific team
- `GetAllTeams(filters...)` - Retrieve all teams with optional filtering
  - Filter by `TeamIDs`, `Countries`, or `HomeRegions`
  - Example: `GetAllTeams(TeamFilter{Countries: []string{"USA", "Canada"}})`
- `GetTeamsByRegion(region)` - Retrieve all teams in a specific home region
- `SaveTeam(team)` - Insert or update a team

### Events

- `GetEvent(eventID)` - Retrieve a specific event
- `GetAllEvents(filters...)` - Retrieve all events with optional filtering
  - Filter by `EventCodes`, `RegionCodes`, or `Countries`
  - Example: `GetAllEvents(EventFilter{Countries: []string{"USA"}})`
- `SaveEvent(event)` - Insert or update an event
- `GetRegionCodes()` - Get all unique region codes
- `GetEventCodesByRegion(regionCode)` - Get all event codes for a specific region
- `GetEventAwards(eventID)` - Get awards for an event
- `GetTeamAwardsByEvent(eventID, teamID)` - Get all awards for a specific team at a specific event
- `GetAllTeamAwards(teamID)` - Get all awards for a specific team across all events
- `SaveEventAward(eventAward)` - Record an award
- `GetEventRankings(eventID)` - Get team rankings for an event
- `SaveEventRanking(ranking)` - Update rankings
- `GetEventAdvancements(eventID)` - Get advancing teams from an event
- `GetAdvancementsByRegion(regionCode)` - Get all advancements from events in a specific region
- `GetAllAdvancements(filters...)` - Get all advancements from all events with optional filtering
  - Filter by `Countries` or `RegionCodes`
  - Example: `GetAllAdvancements(AdvancementFilter{RegionCodes: []string{"USCALA"}})`
- `SaveEventAdvancement(advancement)` - Record advancement

### Matches

- `GetMatch(matchID)` - Retrieve a specific match
- `GetAllMatches(filters...)` - Retrieve all matches with optional filtering
  - Filter by `EventIDs`
  - Example: `GetAllMatches(MatchFilter{EventIDs: []string{"EVENT-123 : 2024"}})`
- `GetMatchesByEvent(eventID)` - Retrieve all matches for a specific event
- `SaveMatch(match)` - Insert or update a match
- `GetMatchAllianceScore(matchID, alliance)` - Get alliance score
- `SaveMatchAllianceScore(score)` - Update alliance score
- `GetMatchTeams(matchID)` - Get teams in a match
- `GetTeamsByEvent(eventID)` - Get all unique team IDs that participated in a specific event
- `SaveMatchTeam(matchTeam)` - Record team participation

### Awards

- `GetAward(awardID)` - Retrieve a specific award
- `GetAllAwards()` - Retrieve all awards
- `SaveAward(award)` - Insert or update an award

## Filter Types

The database supports flexible filtering for retrieving data:

### TeamFilter

```go
type TeamFilter struct {
    TeamIDs     []int    // Filter by team IDs
    Countries   []string // Filter by countries
    HomeRegions []string // Filter by home regions
}
```

### EventFilter

```go
type EventFilter struct {
    EventCodes  []string // Filter by event codes
    RegionCodes []string // Filter by region codes
    Countries   []string // Filter by countries
}
```

### MatchFilter

```go
type MatchFilter struct {
    EventIDs []string // Filter by event IDs
}
```

### AdvancementFilter

```go
type AdvancementFilter struct {
    Countries   []string // Filter by countries
    RegionCodes []string // Filter by region codes
}
```

**Filter Logic:**

- Multiple values within the same field use OR logic (e.g., `Countries: []string{"USA", "Canada"}` matches USA OR Canada)
- Multiple fields use AND logic (e.g., filtering by both Country AND Region requires both to match)
- Omitting a filter returns all records

## Development

### Building

The project includes a Makefile for building cross-platform binaries. See the [Makefile](Makefile) for build targets.

Build for all platforms:

```bash
make build
```

 and database connection

- **database/db.go**: Database interface definition
- **database/sql.go**: SQL database implementation with connection pooling
- **database/sql_*.go**: SQL-specific operations for each entity type
- **database/filedb.go**: File-based database implementation
- **database/filedb_*.go**: File-based operations for each entity type
- **database/award.go, event.go, match.go, team.go**: Data models with SQL query constants and String() methods
- **database/statements.go**: Initialization of prepared statements

### Architecture Highlights

1. **SQL Query Constants**: All SQL queries are defined as package-level constants in their respective model files (e.g., `getTeamQuery`, `saveEventQuery`), making them easy to find, update, and maintain.

2. **Interface-Based Design**: The `DB` interface allows seamless switching between database backends without changing application code.

3. **String Representations**: All models have pointer-receiver String() methods that provide formatted output for logging and debugging:
   - `Award`: Shows ID, name, and type (Team/Person)
   - `Team`: Shows ID, name, city, state, and region
   - `Event`: Shows ID, code, name, year, and location
   - `EventAward`: Shows event ID, team ID, and award ID
   - `EventRanking`: Shows event ID, team ID, rank, and win-loss-tie record
   - `EventAdvancement`: Shows event ID and team ID
   - `Match`: Shows ID, event ID, number, and tournament level
   - `MatchAllianceScore`: Shows match ID, alliance, and point breakdown
   - `MatchTeam`: Shows match ID, team ID, alliance, and status (DQ/Surrogate)

```bash
make build-linux      # Linux AMD64
make build-mac-amd    # macOS Intel
make build-mac-arm    # macOS ARM (Apple Silicon)
make build-windows    # Windows AMD64
```

Binaries will be output to the `bin/` directory under the respective platform subdirectories.

Clean build artifacts:

```bash
make clean
```

### Testing

```bash
go test ./...
```

### Code Organization

- **cmd/ftc/main.go**: Application initialization, database connection, and prepared statement setup
- **database/db.go**: Database connection management and prepared statement caching
- **dbmodel/**: Data models and database operations for each entity type

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
