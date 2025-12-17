package handlers

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	queryv1 "github.com/ruhika006/funding-analytics/gen/api/v1"
)

// GetTopIndustries retrieves top funded industries across all time
func (s *Server) GetTopIndustries(ctx context.Context, req *queryv1.GetTopIndustriesRequest) (*queryv1.GetTopIndustriesResponse, error) {
	
	// add Validation.
	if err := protovalidate.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	limit := req.Limit

	query := fmt.Sprintf(`
		SELECT
			Industry,
			CAST(SUM(Funding_Amount_USD) AS VARCHAR) AS total_funds
		FROM startup
		GROUP BY Industry
		ORDER BY SUM(Funding_Amount_USD) DESC
		LIMIT %d
	`, limit)

	rows, err := s.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var industries []*queryv1.IndustryFunding
	for rows.Next() {
		var industry,totalFunds string
		if err := rows.Scan(&industry, &totalFunds); err != nil {
			return nil, s.ScanError(err, "GetTopIndustries")
		}
		industries = append(industries, &queryv1.IndustryFunding{
			Industry:   industry,
			TotalFunds: totalFunds,
		})
	}

	if err := s.CheckRowsError(rows); err != nil {
		return nil, err
	}

	return &queryv1.GetTopIndustriesResponse{
		Industries: industries,
	}, nil
}

// GetTopIndustriesByYear retrieves top industries for a specific year
func (s *Server) GetTopIndustriesByYear(ctx context.Context, req *queryv1.GetTopIndustriesByYearRequest) (*queryv1.GetTopIndustriesByYearResponse, error) {
	year := req.Year
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	query := fmt.Sprintf(`
		SELECT
			Industry,
			CAST(SUM(Funding_Amount_USD) AS VARCHAR) AS total_funds
		FROM startup
		WHERE Year = %d
		GROUP BY Industry
		ORDER BY SUM(Funding_Amount_USD) DESC
		LIMIT %d
	`, year, limit)

	rows, err := s.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var industries []*queryv1.IndustryFunding
	for rows.Next() {
		var industry, totalFunds string
		if err := rows.Scan(&industry, &totalFunds); err != nil {
			return nil, s.ScanError(err, "GetTopIndustriesByYear")
		}
		industries = append(industries, &queryv1.IndustryFunding{
			Industry:   industry,
			TotalFunds: totalFunds,
		})
	}

	if err := s.CheckRowsError(rows); err != nil {
		return nil, err
	}

	return &queryv1.GetTopIndustriesByYearResponse{
		Industries: industries,
	}, nil
}
