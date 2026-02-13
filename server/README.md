# FTC Standing HTTP API Documentation

The FTC Standing HTTP API provides REST endpoints to access FTC (FIRST Tech Challenge) competition data including teams, events, matches, awards, and rankings.

## Getting Started

### Starting the Server

```bash
# Start on default port 8080
ftcserver

# Start on custom port
ftcserver --port 3000

# Specify a default season (optional)
ftcserver --season 2024
```

### Health Check

``` http
GET /health
```

Returns server health status.

## API Endpoints

All API endpoints are versioned under `/v1` and require a season parameter in the path.

### Base URL Format

``` http
http://localhost:8080/v1/{season}/{resource}
```

### Teams

#### Get Team Details

``` http
GET /v1/{season}/team/{teamID}
```

Returns detailed information about a specific team.

**Example:**

``` http
GET /v1/2024/team/12345
```

#### List Teams

``` http
GET /v1/{season}/teams?limit={limit}
GET /v1/{season}/teams/{region}?limit={limit}
```

Returns all teams. If region is specified, filters to teams in that region.

**Query Parameters:**

- `limit` (optional): Limit number of results

**Examples:**

``` http
# All teams
GET /v1/2024/teams

# Teams in a specific region
GET /v1/2024/teams/USCHS

# First 100 teams
GET /v1/2024/teams?limit=100
```

### Events

#### Get Event Teams

``` http
GET /v1/{season}/events/{eventCode}/teams?limit={limit}
```

Returns event information with an embedded array of teams that participated in the event.

**Response structure:**

```json
{
  "event": {
    "event_id": "...",
    "event_code": "...",
    "name": "...",
    ...
    "teams": [...]
  }
}
```

**Query Parameters:**

- `limit` (optional): Limit number of results

**Examples:**

``` http
GET /v1/2024/events/USNCCOQ/teams

# First 20 teams
GET /v1/2024/events/USNCCOQ/teams?limit=20
```

#### Get Event Rankings

``` http
GET /v1/{season}/events/{eventCode}/rankings?limit={limit}
```

Returns event information along with an array of team rankings at the event.

**Response structure:**

```json
{
  "event": {...},
  "rankings": [...]
}
```

**Query Parameters:**

- `limit` (optional): Limit number of results

**Examples:**

``` http
GET /v1/2024/events/USNCCOQ/rankings

# Top 10 ranked teams
GET /v1/2024/events/USNCCOQ/rankings?limit=10
```

#### Get Event Awards

``` http
GET /v1/{season}/events/{eventCode}/awards?limit={limit}
```

Returns event information along with an array of awards given at the event.

**Response structure:**

```json
{
  "event": {...},
  "awards": [...]
}
```

**Note:** Award objects do not include `event_id` since the event information is already provided at the top level.

**Query Parameters:**

- `limit` (optional): Limit number of results

**Examples:**

``` http
GET /v1/2024/events/USNCCOQ/awards

# First 5 awards
GET /v1/2024/events/USNCCOQ/awards?limit=5
```

#### Get Event Advancement

``` http
GET /v1/{season}/events/{eventCode}/advancement
```

Returns advancement report for an event showing which teams advanced.

**Example:**

``` http
GET /v1/2024/events/USNCCOQ/advancement
```

#### Get Event Matches

``` http
GET /v1/{season}/events/{eventCode}/matches?team={teamID}&limit={limit}
```

Returns event information along with an array of matches at the event. Optional `team` query parameter filters to matches for a specific team.

**Response structure:**

```json
{
  "event": {...},
  "matches": [...]
}
```

**Query Parameters:**

- `team` (optional): Filter to matches for a specific team
- `limit` (optional): Limit number of results

**Examples:**

``` http
# All matches at an event
GET /v1/2024/events/USNCCOQ/matches

# Matches for a specific team at an event
GET /v1/2024/events/USNCCOQ/matches?team=12345

# First 10 matches at an event
GET /v1/2024/events/USNCCOQ/matches?limit=10

# First 5 matches for a specific team
GET /v1/2024/events/USNCCOQ/matches?team=12345&limit=5
```

### Team Performance Rankings

#### Get Team Rankings (Consolidated)

``` http
GET /v1/{season}/team-rankings?region={region}&country={country}&event={eventCode}&limit={limit}
```

Returns team performance rankings consolidated across all events.

**Query Parameters:**

- `region` (optional): Filter by region code
- `country` (optional): Filter by country
- `event` (optional): Filter by specific event
- `limit` (optional): Limit number of results

**Examples:**

``` http
# All teams for a season
GET /v1/2024/team-rankings

# Teams in a specific region
GET /v1/2024/team-rankings?region=USCHS

# Top 50 teams
GET /v1/2024/team-rankings?limit=50

# Teams at a specific event
GET /v1/2024/team-rankings?event=USNCCOQ
```

#### Get Team Event Rankings (By Event)

``` http
GET /v1/{season}/team-event-rankings?region={region}&country={country}&event={eventCode}&limit={limit}
```

Returns team performance rankings by individual event (not consolidated).

**Query Parameters:**

- `region` (optional): Filter by region code
- `country` (optional): Filter by country
- `event` (optional): Filter by specific event
- `limit` (optional): Limit number of results

**Example:**

``` http
GET /v1/2024/team-event-rankings?region=USCHS&limit=100
```

### Regional Advancement

#### Get Region Advancement

``` http
GET /v1/{season}/regions/{regionCode}/advancement
```

Returns advancement information for all teams in a region.

**Example:**

``` http
GET /v1/2024/regions/USCHS/advancement
```

### Event Advancement Summary

#### Get All Advancement

``` http
GET /v1/{season}/advancement?region={region}
```

Returns advancement organized by qualifying events.

**Query Parameters:**

- `region` (optional): Filter by region code, defaults to "ALL"

**Examples:**

``` http
# All advancement
GET /v1/2024/advancement

# Advancement for a specific region
GET /v1/2024/advancement?region=USCHS
```

## Response Format

All successful responses return JSON with the appropriate data structure. Errors return JSON with an `error` field:

```json
{
  "error": "error message here"
}
```

## HTTP Status Codes

- `200 OK` - Successful request
- `400 Bad Request` - Invalid parameters or missing required fields
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Examples

### cURL Examples

```bash
# Get team details
curl http://localhost:8080/v1/2024/team/12345

# Get event rankings
curl http://localhost:8080/v1/2024/events/USNCCOQ/rankings

# Get team rankings for a region with limit
curl "http://localhost:8080/v1/2024/team-rankings?region=USCHS&limit=20"

# Get matches for a specific team at an event
curl "http://localhost:8080/v1/2024/events/USNCCOQ/matches?team=12345"

# Get first 100 teams
curl "http://localhost:8080/v1/2024/teams?limit=100"

# Get top 10 event rankings with limit
curl "http://localhost:8080/v1/2024/events/USNCCOQ/rankings?limit=10"
```

### JavaScript/Fetch Example

```javascript
// Get team rankings
fetch('http://localhost:8080/v1/2024/team-rankings?region=USCHS&limit=10')
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

## Environment Variables

The server respects the following environment variables:

- `FTC_SEASON` - Default season to use if not specified in the command line
- `LOG_LEVEL` - Logging level (debug, info, warm, error)
- `DB_TYPE` - Database type (sql or file)
- `DATA_SOURCE_NAME` - Database connection string (for SQL databases)
- `FILEDB_DATA_DIR` - Base directory for file-based database

## Server Configuration

The server includes the following features:

- Graceful shutdown on SIGINT/SIGTERM
- Configurable read/write timeouts (15 seconds)
- Configurable idle timeout (60 seconds)
- Health check endpoint for monitoring
