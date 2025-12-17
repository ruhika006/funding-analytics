package handlers

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	queryv1 "github.com/ruhika006/funding-analytics/gen/api/v1"
)

// GetTopFundedStartups retrieves most funded companies across all time
func (s *Server) GetTopFundedStartups(ctx context.Context, req *queryv1.GetTopFundedStartupsRequest) (*queryv1.GetTopFundedStartupsResponse, error) {
	
	// add Validation.
	if err := protovalidate.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	limit := req.Limit

	query := fmt.Sprintf(`
		SELECT
			Company,
			CAST(SUM(Funding_Amount_USD) AS VARCHAR) AS total_funds
		FROM startup
		GROUP BY Company
		ORDER BY SUM(Funding_Amount_USD) DESC
		LIMIT %d
	`, limit)

	rows, err := s.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var startups []*queryv1.StartupFunding
	for rows.Next() {
		var company, totalFunds string
		if err := rows.Scan(&company, &totalFunds); err != nil {
			return nil, s.ScanError(err, "GetTopFundedStartups")
		}
		startups = append(startups, &queryv1.StartupFunding{
			Company:    company,
			TotalFunds: totalFunds,
		})
	}

	if err := s.CheckRowsError(rows); err != nil {
		return nil, err
	}

	return &queryv1.GetTopFundedStartupsResponse{
		Startups: startups,
	}, nil
}
