# Uber Extractor

A CLI tool for extracting and exporting Uber trip data from your account.

## Features

- Fetch complete trip history from Uber's GraphQL API
- Export data in JSON or CSV format
- Location clustering and tracking
- Date range filtering with flexible syntax
- Summary views for quick analysis

## Installation

From source:

```bash
make install
```

Or build manually:

```bash
make build
./bin/ue
```

## Usage

### Authentication

First, log in with your Uber account:

```bash
ue login
```

You'll need to provide your Uber cookie from browser dev tools. Open browser, go to `riders.uber.com`, open developer tools (F12), go to Application > Cookies, and copy the value of the `__cf_bm` or similar cookie.

Check your login status:

```bash
ue status
```

Logout if needed:

```bash
ue logout
```

### Fetching Trips

Fetch last 7 days of trips in JSON format:

```bash
ue trips --last 7d
```

Fetch trips for a specific date range in CSV format:

```bash
ue trips --from 2024-01-01 --to 2024-01-31 --output csv
```

Show summary without fetching full details:

```bash
ue trips --last 30d --summary
```

### Viewing Locations

List all saved locations (clustered from trip data):

```bash
ue locations
```

## Development

Build:

```bash
make build
```

Run tests:

```bash
make test
```

Format code:

```bash
make fmt
```

Clean build artifacts:

```bash
make clean
```

## Project Structure

```
cmd/ue/              # CLI entry point
cmd/                 # Command implementations
internal/            # Internal packages
  auth/              # Authentication logic
  uberapi/           # Uber API client
  locations/         # Location clustering
  trips/             # Trip data models
  format/            # Output formatting (JSON, CSV)
  datetime/          # Date/time utilities
  parser/            # Data parsing
  transform/         # Data transformation
```

## Requirements

- Go 1.25.4 or later