package handlers

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	queryv1 "github.com/ruhika006/funding-analytics/gen/api/v1"
)

// GetRecords retrieves all startup records with limit
func (s *Server) GetRecords(ctx context.Context, req *queryv1.GetRecordsRequest) (*queryv1.GetRecordsResponse, error) {
	limit := req.Limit

	// add Validation.
	if err := protovalidate.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	query := fmt.Sprintf(
		`SELECT Company, Industry, Funding_Amount_USD, Investor, Year, City 
		FROM startup
		LIMIT %d`,limit,
	)

	rows, err := s.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*queryv1.StartupRecord
	for rows.Next() {
		var company, industry, funding, investor, city string
		var year int64
		if err := rows.Scan(&company, &industry, &funding, &investor, &year, &city); err != nil {
			return nil, s.ScanError(err, "GetRecords")
		}
		records = append(records, &queryv1.StartupRecord{
			Company:       company,
			Industry:      industry,
			FundingAmount: funding,
			Investor:      investor,
			Year:          int32(year),
			City:          city,
		})
	}

	if err := s.CheckRowsError(rows); err != nil {
		return nil, err
	}

	return &queryv1.GetRecordsResponse{
		Records: records,
	}, nil
}
