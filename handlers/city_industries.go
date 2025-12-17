package handlers

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	queryv1 "github.com/ruhika006/funding-analytics/gen/api/v1"
)

// GetTopCityAndIndustries retrieves top industries by city with optional filters
func (s *Server) GetTopCityAndIndustries(ctx context.Context, req *queryv1.GetTopCityAndIndustriesRequest) (*queryv1.GetTopCityAndIndustriesResponse, error) {
	
	// add Validation.
	if err := protovalidate.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	limit := req.Limit

	// whereClause := buildWhereClause(req.City, req.Industry)

	query := fmt.Sprintf(`
		SELECT
			Industry,
			City,
			CAST(SUM(Funding_Amount_USD) AS VARCHAR) AS total_funds
		FROM startup
		GROUP BY Industry, City
		ORDER BY SUM(Funding_Amount_USD) DESC
		LIMIT %d
	`,limit)

	rows, err := s.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*queryv1.CityIndustryFunding
	for rows.Next() {
		var industry, city, totalFunds string
		if err := rows.Scan(&industry, &city, &totalFunds); err != nil {
			return nil, s.ScanError(err, "GetTopCityAndIndustries")
		}
		results = append(results, &queryv1.CityIndustryFunding{
			Industry:   industry,
			City:       city,
			TotalFunds: totalFunds,
		})
	}

	if err := s.CheckRowsError(rows); err != nil {
		return nil, err
	}

	return &queryv1.GetTopCityAndIndustriesResponse{
		Results: results,
	}, nil
}

// buildWhereClause constructs the WHERE clause based on filters
func buildWhereClause(city, industry string) string {
	if city != "" && industry != "" {
		return fmt.Sprintf(" WHERE lower(City) = lower('%s') AND lower(Industry) = lower('%s')", city, industry)
	}
	if city != "" {
		return fmt.Sprintf(" WHERE lower(City) = lower('%s')", city)
	}
	if industry != "" {
		return fmt.Sprintf(" WHERE lower(Industry) = lower('%s')", industry)
	}
	return ""
}
