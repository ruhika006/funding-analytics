package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Server implements QueryService
type Server struct {
	client clickhouse.Conn
}

// NewServer creates a new query handler server
func NewServer(client clickhouse.Conn) *Server {
	return &Server{
		client: client,
	}
}

// ExecuteQuery runs a SQL query and handles common error logging
func (s *Server) ExecuteQuery(ctx context.Context, query string) (driver.Rows, error) {
	rows, err := s.client.Query(ctx, query)
	if err != nil {
		log.Printf("ClickHouse query error: %v\nQuery: %s", err, query)
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	return rows, nil
}

// CheckRowsError checks for errors after row iteration completes
func (s *Server) CheckRowsError(rows driver.Rows) error {
	if err := rows.Err(); err != nil {
		log.Printf("ClickHouse rows iteration error: %v", err)
		return fmt.Errorf("rows iteration failed: %w", err)
	}
	return nil
}

// ScanError logs scan errors with context
func (s *Server) ScanError(err error, context string) error {
	log.Printf("ClickHouse scan error in %s: %v", context, err)
	return fmt.Errorf("failed to scan row in %s: %w", context, err)
}
