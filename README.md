# Startup Funding Database API

A Connect-RPC API server for querying startup funding data stored in ClickHouse.

## Overview

- **Database**: ClickHouse (analytical OLAP database)
- **Framework**: Connect-RPC (gRPC-compatible protocol)
- **API Protocol**: HTTP/2 with Protocol Buffers
- **Port**: 8080 (configurable via `PORT` env var)
- **Language**: Go 1.21+

## Getting Started

### Prerequisites

- ClickHouse server running (locally or remote)
- Go 1.21 or higher
- Environment variables configured

### Configuration

Set environment variables before running:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MYAPP_USER` | Yes | - | ClickHouse username |
| `MYAPP_PASSWORD` | Yes | - | ClickHouse password |
| `MYAPP_ADDR` | Yes | - | ClickHouse address (host:port) |
| `PORT` | No | 8080 | API server port |

```bash
export MYAPP_USER=user
export MYAPP_PASSWORD=password
export MYAPP_ADDR=localhost:9000
export PORT=8080  # optional, defaults to 8080
```

### Setup ClickHouse

Start a ClickHouse server using Docker:

```bash
docker run --rm -d \
    --name clickhouse-server \
    -e CLICKHOUSE_USER=user \
    -e CLICKHOUSE_PASSWORD=password \
    -p 9000:9000 \
    -p 8123:8123 \
    clickhouse/clickhouse-server
```

### Create Table and Load Data

Connect to ClickHouse:

```bash
clickhouse-client -h localhost -u user --password password
```

Create the startup table:

```sql
CREATE TABLE startup (
    Company String,
    Industry LowCardinality(String),
    Funding_Amount_USD Int64,
    Investor String,
    Year Int64,
    City LowCardinality(String)
  )
ENGINE = MergeTree
PRIMARY KEY (Industry, Year)
```

Insert data from CSV:

```bash
cat app/startup_funding.csv | clickhouse-client \
    --query "INSERT INTO startup (Company, Industry, Funding_Amount_USD, Investor, Year, City) FORMAT CSV" \
    -u user --password password
```

Finally start server
The server will:
1. Load configuration from environment variables
2. Connect to ClickHouse
3. Start HTTP/2 API server on the configured port

```bash
go run main.go
```


### BETTER OPTION: Running with Docker

Run with ClickHouse:

```bash
docker run -p 8080:8080 ruhika0817/db-app
```
Next, starting hitting API Endpoints.

## API Endpoints

All endpoints use POST with JSON bodies. Connect-RPC provides gRPC reflection support.

### Get Records

Retrieve all startup records with optional limit:

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
  -d '{"year": 2020, "limit": 10}' | jq
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

## Data Schema

The `startup` table contains:
- `Company` (String): Company name
- `Industry` (LowCardinality String): Industry category
- `Funding_Amount_USD` (Int64): Funding amount in USD
- `Investor` (String): Primary investor
- `Year` (Int64): Funding year
- `City` (LowCardinality String): City location

**Storage**: MergeTree engine with primary key on (Industry, Year) for optimal query performance.

## Architecture

- **Type-safe API**: Protocol Buffers with code generation via Buf
- **Connect-RPC Handlers**: Adapter layer for request/response wrapping
- **gRPC Reflection**: Enabled for API discovery and debugging
- **Environment Configuration**: 12-factor app compliant
- **Error Handling**: Structured logging and error propagation
- **ClickHouse Integration**: Efficient analytical queries on large datasets

## Development

### Generate API Code

```bash
buf generate
```

### Build

```bash
go build -o funding-analytics
```

## Deployment

### AWS Deployment with Terraform

To deploy to AWS, infrastructure-as-code blocks are available in `main.tf`. 
The application is by default deployed to AWS App Runner.

```bash
terraform apply
```

This will provision the necessary AWS resources for the application including compute, networking, and database connectivity.


## Troubleshooting

**Connection refused**: Ensure ClickHouse is running and accessible at the configured address.

**Table not found**: Create the table and insert data as described in the Setup section.

**Health check timeout**: Increase the timeout or check ClickHouse logs for issues.
