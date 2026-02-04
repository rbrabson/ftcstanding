# FTC API Client

A Go client library for accessing the FIRST Tech Challenge (FTC) Events API.

[Official FTC Events API Documentation](https://ftc-events.firstinspires.org/api-docs/index.html)

## Overview

This library provides a comprehensive Go client for interacting with the FTC Events API, allowing you to retrieve information about teams, events, matches, scores, awards, and more.

## Installation

```bash
go get github.com/rbrabson/ftc
```

## Prerequisites

You need FTC API credentials to use this library. Set the following environment variables:

- `FTC_USERNAME`: Your FTC API username
- `FTC_AUTHORIZATION_KEY`: Your FTC API authorization key

You can also use a `.env` file in your project root:

``` ini
FTC_USERNAME=your_username
FTC_AUTHORIZATION_KEY=your_key
```

## Usage

```go
import "github.com/rbrabson/ftc/ftc"
```

### API Index

Get information about the FTC API server:

```go
apiIndex, err := ftc.GetApiIndex()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("API Version: %s\n", apiIndex.APIVersion)
fmt.Printf("Current Season: %d\n", apiIndex.CurrentSeason)
```

### Teams

Retrieve team information:

```go
// Get all teams for a season
teams, err := ftc.GetTeams("2024")
if err != nil {
    log.Fatal(err)
}

// Get a specific team
teams, err := ftc.GetTeams("2024", "teamNumber=12345")
if err != nil {
    log.Fatal(err)
}
```

### Events

Retrieve event information:

```go
// Get all events for a season
events, err := ftc.GetEvents("2024")
if err != nil {
    log.Fatal(err)
}

// Get events by event code
events, err := ftc.GetEvents("2024", map[string]string{"eventCode": "USNCCMP"})
if err != nil {
    log.Fatal(err)
}

// Get events by team number
events, err := ftc.GetEvents("2024", map[string]string{"teamNumber": "12345"})
if err != nil {
    log.Fatal(err)
}
```

### Matches

Get match information for an event:

```go
// Get all matches
matches, err := ftc.GetMatches("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}

// Get qualifier matches only
matches, err := ftc.GetMatches("2024", "USNCCMP", ftc.QUALIFIER)
if err != nil {
    log.Fatal(err)
}

// Get playoff matches only
matches, err := ftc.GetMatches("2024", "USNCCMP", ftc.PLAYOFF)
if err != nil {
    log.Fatal(err)
}
```

### Schedule

Retrieve event schedules:

```go
// Get schedule for an event
schedule, err := ftc.GetEventSchedule("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}

// Get hybrid schedule (for hybrid events)
hybridSchedule, err := ftc.GetHybridSchedule("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}
```

### Scores

Get match scores:

```go
scores, err := ftc.GetScores("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}

// Filter by match level
scores, err := ftc.GetScores("2024", "USNCCMP", ftc.QUALIFIER)
if err != nil {
    log.Fatal(err)
}
```

### Awards

Retrieve award information:

```go
// Get list of available awards for a season
awards, err := ftc.GetAwards("2024")
if err != nil {
    log.Fatal(err)
}

// Get awards for a specific event
eventAwards, err := ftc.GetEventAwards("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}

// Get awards won by a team
teamAwards, err := ftc.GetTeamAwards("2024", 12345)
if err != nil {
    log.Fatal(err)
}
```

### Rankings

Get event rankings:

```go
rankings, err := ftc.GetRankings("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}
```

### Alliances

Get alliance selections:

```go
alliances, err := ftc.GetAlliances("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}
```

### Advancement

Get advancement information:

```go
advancement, err := ftc.GetAdvancement("2024", "USNCCMP")
if err != nil {
    log.Fatal(err)
}
```

### Leagues

Get league information:

```go
leagues, err := ftc.GetLeagues("2024")
if err != nil {
    log.Fatal(err)
}
```

## Features

- ✅ Complete coverage of FTC Events API endpoints
- ✅ Automatic authentication using environment variables
- ✅ Type-safe Go structs for all API responses
- ✅ Comprehensive test coverage
- ✅ Support for pagination and query parameters
- ✅ Custom time parsing for FTC date formats

## API Endpoints Covered

- **API Index**: Server information and current season
- **Teams**: Team information and listings
- **Events**: Event details, schedules, and information
- **Matches**: Match results and details
- **Schedule**: Event schedules (standard and hybrid)
- **Scores**: Detailed match scores
- **Awards**: Award definitions and winners
- **Rankings**: Team rankings at events
- **Alliances**: Alliance selections
- **Advancement**: Championship advancement information
- **Leagues**: League information

## Requirements

- Go 1.22.1 or higher

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Environment variable management

## Testing

Run the test suite:

```bash
go test ./...
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
