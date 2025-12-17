# Startup Funding Database API

A gRPC/Connect-RPC API server for querying startup funding data stored in DuckDB.

## Overview

- **Database**: DuckDB
- **Framework**: Connect-RPC (gRPC-compatible)
- **Data Source**: CSV file (`startup_funding.csv`)
- **Server**: HTTP/2 on port 8080
- **Language**: Go

## Getting Started

### Running Locally

```bash
go run main.go
```

The server will:
1. Open DuckDB connection
2. Load startup funding data from `./app/startup_funding.csv` into a table
3. Start HTTP server on `:8080`

### Running with Docker

```bash
docker run -p 8080:8080 db-app
```

Then access from host:
```bash
curl http://localhost:8080
```

## API Endpoints

All endpoints use POST with JSON bodies. Connect-RPC provides gRPC reflection support.

### Get Records

```bash
curl -X POST http://localhost:8080/api.v1.QueryService/GetRecords \
  -H "Content-Type: application/json" \
  -d '{"limit": 10}' | jq
```


### Get Top Industries (All Time)

```bash
curl -X POST http://localhost:8080/api.v1.QueryService/GetTopIndustries \
  -H "Content-Type: application/json" \
  -d '{"limit": 10}' | jq
```

### Get Top Industries by Year

```bash
curl -X POST http://localhost:8080/api.v1.QueryService/GetTopIndustriesByYear \
  -H "Content-Type: application/json" \
  -d '{"year": 2025, "limit": 10}' | jq
```

### Get Top Funded Startups

```bash
curl -X POST http://localhost:8080/api.v1.QueryService/GetTopFundedStartups \
  -H "Content-Type: application/json" \
  -d '{"limit": 10}' | jq
```

### Get Top by City and Industry

```bash
curl -X POST http://localhost:8080/api.v1.QueryService/GetTopCityAndIndustries \
  -H "Content-Type: application/json" \
  -d '{"city": "Pune", "industry": "hospitality", "limit": 5}' | jq
```

### Note : With AppRunner :

curl -X POST <DEFAULT_DOMAIN>/api.v1.QueryService/GetRecords \
  -H "Content-Type: application/json" \
  -d '{"limit": 10}' | jq


## Data Schema

The startup table contains the following columns (loaded from CSV):
- `Company` (String)
- `Industry` (String)
- `Funding_Amount_USD` (Int64)
- `Investor` (String)
- `Year` (Int64)
- `City` (String)

## Architecture Notes

- Uses Connect-RPC handlers for type-safe API communication
- gRPC reflection enabled for API discovery
- DuckDB provides efficient analytical queries on CSV data
- Alternative ClickHouse implementation available (commented in code)

