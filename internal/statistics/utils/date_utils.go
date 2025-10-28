package utils

import (
	"cibo/internal/types"
	"fmt"
	"time"
)

const layout = "2006-01-02" // The reference layout for YYYY-MM-DD in Go.

// Filters a slice of DailyStockRecord based on a start and end date and returns the data WITHIN those dates.
// If startDateStr or endDateStr are empty, they are ignored and all data returned beyonds those values.
func FilterDailyPricesWithinDateRange(records []types.DailyStockRecord, startDateStr, endDateStr string) ([]types.DailyStockRecord, error) {
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(layout, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(layout, endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
	}

	filteredRecords := []types.DailyStockRecord{}

	for _, record := range records {
		recordDate, err := time.Parse(layout, record.Date)
		if err != nil {
			continue // Skip records with un-parse-able dates
		}

		isAfterStartDate := startDate.IsZero() || !recordDate.Before(startDate)
		isBeforeEndDate := endDate.IsZero() || !recordDate.After(endDate)

		if isAfterStartDate && isBeforeEndDate {
			filteredRecords = append(filteredRecords, record)
		}
	}

	return filteredRecords, nil
}

// Filters a slice of AnnualEarningRecord based on a start and end date and returns the data WITHIN those dates.
// If startDateStr or endDateStr are empty, they are ignored and all data returned beyonds those values.
func FilterAnnualEarningsWithinDateRange(records []types.AnnualEarningRecord, startDateStr, endDateStr string) ([]types.AnnualEarningRecord, error) {
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(layout, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(layout, endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
	}

	filteredRecords := []types.AnnualEarningRecord{}
	for _, record := range records {
		recordDate, err := time.Parse(layout, record.FiscalDateEnding)
		if err != nil {
			continue // Skip records with un-parse-able dates
		}

		isAfterStartDate := startDate.IsZero() || !recordDate.Before(startDate)
		isBeforeEndDate := endDate.IsZero() || !recordDate.After(endDate)

		if isAfterStartDate && isBeforeEndDate {
			filteredRecords = append(filteredRecords, record)
		}
	}

	return filteredRecords, nil
}
