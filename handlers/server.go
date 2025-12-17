package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Server implements QueryService
type Server struct {
	DBClient *sql.DB
}

// NewServer creates a new query handler server
func NewServer(dbClient *sql.DB) *Server {
	return &Server{
		DBClient: dbClient,
	}
}

// ExecuteQuery runs a SQL query and handles common error logging
func (s *Server) ExecuteQuery(ctx context.Context, query string) (*sql.Rows, error) {
	rows, err := s.DBClient.Query(query)
	if err != nil {
		log.Printf("DuckDB query error: %v\nQuery: %s", err, query)
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	return rows, nil
}

// CheckRowsError checks for errors after row iteration completes
func (s *Server) CheckRowsError(rows *sql.Rows) error {
	if err := rows.Err(); err != nil {
		log.Printf("DuckDB rows iteration error: %v", err)
		return fmt.Errorf("rows iteration failed: %w", err)
	}
	return nil
}

// ScanError logs scan errors with context
func (s *Server) ScanError(err error, context string) error {
	log.Printf("DuckDB scan error in %s: %v", context, err)
	return fmt.Errorf("failed to scan row in %s: %w", context, err)
}
