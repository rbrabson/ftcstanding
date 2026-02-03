# ftcstanding

A Go application for managing FTC (FIRST Tech Challenge) competition standings, teams, matches, events, and awards.

## Project Structure

``` code
ftcstanding/
├── cmd/
│   └── ftc/
│       └── main.go           # Application entry point
├── database/
│   └── db.go                 # Database connection and prepared statement management
├── dbmodel/
│   ├── award.go              # Award data model and database operations
│   ├── event.go              # Event data model and database operations
│   ├── match.go              # Match data model and database operations
│   └── team.go               # Team data model and database operations
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
├── .env                      # Environment configuration (not in git)
├── LICENSE
└── README.md
```

## Features

- **Database Connection Management**: MySQL database with connection pooling
- **Prepared Statements**: All database operations use prepared statements for performance and security
- **Data Models**:
  - Teams: Manage team information including name, location, and rookie year
  - Events: Track competition events with dates, locations, and details
  - Matches: Record match results, alliance scores, and team participation
  - Awards: Manage awards and track which teams received them at events

## Prerequisites

- Go 1.24.0 or later
- MySQL database server

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

3. Create a `.env` file in the project root with your database connection string:

   ``` ini
   DATA_SOURCE_NAME=user:password@tcp(localhost:3306)/dbname
   ```

## Database Setup

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

## Usage

Build and run the application:

```bash
go build -o ftcstanding ./cmd/ftc
./ftcstanding
```

Or run directly:

```bash
go run ./cmd/ftc
```

## Database Operations

All database operations use prepared statements which are initialized at application startup:

### Teams

- `GetTeam(teamID)` - Retrieve a specific team
- `GetAllTeams()` - Retrieve all teams
- `SaveTeam(team)` - Insert or update a team

### Events

- `GetEvent(eventID)` - Retrieve a specific event
- `SaveEvent(event)` - Insert or update an event
- `GetEventAwards(eventID)` - Get awards for an event
- `SaveEventAward(eventAward)` - Record an award
- `GetEventRankings(eventID)` - Get team rankings
- `SaveEventRanking(ranking)` - Update rankings
- `GetEventAdvancements(eventID)` - Get advancing teams
- `SaveEventAdvancement(advancement)` - Record advancement

### Matches

- `GetMatch(matchID)` - Retrieve a specific match
- `GetAllMatches()` - Retrieve all matches
- `SaveMatch(match)` - Insert or update a match
- `GetMatchAllianceScore(matchID, alliance)` - Get alliance score
- `SaveMatchAllianceScore(score)` - Update alliance score
- `GetMatchTeams(matchID)` - Get teams in a match
- `SaveMatchTeam(matchTeam)` - Record team participation

### Awards

- `GetAward(awardID)` - Retrieve a specific award
- `GetAllAwards()` - Retrieve all awards
- `SaveAward(award)` - Insert or update an award

## Development

### Building

```bash
go build ./...
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
