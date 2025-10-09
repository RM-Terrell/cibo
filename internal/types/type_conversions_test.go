package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Given a slice of AnnualEarningRecords, verify that it is correctly converted
// to a slice of AnnualEarningRecordParquet with all data preserved.
func TestAnnualEarningsToParquet_Success(t *testing.T) {
	inputRecords := []AnnualEarningRecord{
		{Ticker: "NVDA", FiscalDateEnding: "2025-01-31", ReportedEPS: 25.50},
		{Ticker: "NVDA", FiscalDateEnding: "2024-01-31", ReportedEPS: 12.05},
		{Ticker: "AMD", FiscalDateEnding: "2024-12-31", ReportedEPS: 1.10},
	}

	expectedOutput := []AnnualEarningRecordParquet{
		{Ticker: "NVDA", FiscalDateEnding: "2025-01-31", ReportedEPS: 25.50},
		{Ticker: "NVDA", FiscalDateEnding: "2024-01-31", ReportedEPS: 12.05},
		{Ticker: "AMD", FiscalDateEnding: "2024-12-31", ReportedEPS: 1.10},
	}

	result := AnnualEarningsToParquet(inputRecords)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("AnnualEarningsToParquet() mismatch (-want +got):\n%s", diff)
	}
}

// Given an empty slice of AnnualEarningRecords, verify that an empty slice
// of AnnualEarningRecordParquet is returned.
func TestAnnualEarningsToParquet_EmptyInput(t *testing.T) {
	inputRecords := []AnnualEarningRecord{}
	expectedOutput := []AnnualEarningRecordParquet{}

	result := AnnualEarningsToParquet(inputRecords)

	if len(result) != 0 {
		t.Errorf("AnnualEarningsToParquet() expected an empty slice for empty input, but got %d elements", len(result))
	}
	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("AnnualEarningsToParquet() mismatch for empty slice (-want +got):\n%s", diff)
	}
}

// Given a nil slice of AnnualEarningRecords, verify that an empty, non-nil slice
// of AnnualEarningRecordParquet is returned.
func TestAnnualEarningsToParquet_NilInput(t *testing.T) {
	var inputRecords []AnnualEarningRecord = nil

	result := AnnualEarningsToParquet(inputRecords)

	if result == nil {
		t.Fatal("AnnualEarningsToParquet() returned a nil slice for nil input, but expected an empty slice")
	}
	if len(result) != 0 {
		t.Errorf("AnnualEarningsToParquet() expected a zero-length slice for nil input, but got %d elements", len(result))
	}
}

// Given a slice of AnnualEarningRecords containing a negative ReportedEPS,
// verify that the negative value is preserved after conversion.
func TestAnnualEarningsToParquet_NegativeValues(t *testing.T) {
	// Arrange
	inputRecords := []AnnualEarningRecord{
		{Ticker: "GROW", FiscalDateEnding: "2025-12-31", ReportedEPS: 2.50},
		{Ticker: "LOSS", FiscalDateEnding: "2025-12-31", ReportedEPS: -1.25},
	}

	expectedOutput := []AnnualEarningRecordParquet{
		{Ticker: "GROW", FiscalDateEnding: "2025-12-31", ReportedEPS: 2.50},
		{Ticker: "LOSS", FiscalDateEnding: "2025-12-31", ReportedEPS: -1.25},
	}

	result := AnnualEarningsToParquet(inputRecords)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("AnnualEarningsToParquet() did not preserve negative value (-want +got):\n%s", diff)
	}
}

// Given a slice of CombinedPriceRecords, verify it is converted correctly.
func TestCombinedPricesToParquet_Success(t *testing.T) {
	inputRecords := []CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "actual_price"},
		{Ticker: "TEST", Date: "2025-12-31", Price: 150.0, Series: "fair_value"},
	}

	expectedOutput := []CombinedPriceRecordParquet{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "actual_price"},
		{Ticker: "TEST", Date: "2025-12-31", Price: 150.0, Series: "fair_value"},
	}

	result := CombinedPricesToParquet(inputRecords)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("CombinedPricesToParquet() mismatch (-want +got):\n%s", diff)
	}
}

// Given an empty slice of CombinedPriceRecords, verify an empty slice is returned.
func TestCombinedPricesToParquet_EmptyInput(t *testing.T) {
	inputRecords := []CombinedPriceRecord{}
	expectedOutput := []CombinedPriceRecordParquet{}

	result := CombinedPricesToParquet(inputRecords)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("CombinedPricesToParquet() mismatch for empty slice (-want +got):\n%s", diff)
	}
}

// Given a nil slice of CombinedPriceRecords, verify a non-nil empty slice is returned.
func TestCombinedPricesToParquet_NilInput(t *testing.T) {
	var inputRecords []CombinedPriceRecord = nil

	result := CombinedPricesToParquet(inputRecords)

	if result == nil {
		t.Fatal("CombinedPricesToParquet() returned a nil slice for nil input, but expected an empty slice")
	}
	if len(result) != 0 {
		t.Errorf("CombinedPricesToParquet() expected a zero-length slice for nil input, but got %d elements", len(result))
	}
}

// Given a record with a negative price, verify the value is preserved.
func TestCombinedPricesToParquet_NegativeValues(t *testing.T) {
	// a negative value here is nonsensical, but this test makes sure the data is preserved during conversion
	inputRecords := []CombinedPriceRecord{
		{Ticker: "GOOD", Date: "2025-12-31", Price: 100.0, Series: "actual_price"},
		{Ticker: "BAD", Date: "2025-12-31", Price: -50.0, Series: "fair_value"},
	}

	expectedOutput := []CombinedPriceRecordParquet{
		{Ticker: "GOOD", Date: "2025-12-31", Price: 100.0, Series: "actual_price"},
		{Ticker: "BAD", Date: "2025-12-31", Price: -50.0, Series: "fair_value"},
	}

	result := CombinedPricesToParquet(inputRecords)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("CombinedPricesToParquet() did not preserve negative value (-want +got):\n%s", diff)
	}
}

// Given slices of daily and fair value prices, verify they are correctly merged.
func TestDailyAndFairPriceToCombined_Success(t *testing.T) {
	dailyPrices := []DailyStockRecord{
		{Ticker: "TEST", Date: "2025-01-01", ClosingPrice: 100.0},
		{Ticker: "TEST", Date: "2025-01-02", ClosingPrice: 102.5},
	}
	fairValuePrices := []FairValuePriceRecord{
		{Ticker: "TEST", Date: "2025-12-31", FairValuePrice: 150.0},
	}
	expectedOutput := []CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-01-02", Price: 102.5, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-12-31", Price: 150.0, Series: "fair_value"},
	}

	result := DailyAndFairPriceToCombined(dailyPrices, fairValuePrices)

	// Use a sorter to make the test robust against the order of appends.
	sorter := cmpopts.SortSlices(func(a, b CombinedPriceRecord) bool { return a.Date < b.Date })
	if diff := cmp.Diff(expectedOutput, result, sorter); diff != "" {
		t.Errorf("DailyAndFairPriceToCombined() mismatch (-want +got):\n%s", diff)
	}
}

// Given one empty input slice, verify only the non-empty slice is converted.
func TestDailyAndFairPriceToCombined_OneEmptyInput(t *testing.T) {
	dailyPrices := []DailyStockRecord{
		{Ticker: "TEST", Date: "2025-01-01", ClosingPrice: 100.0},
	}
	fairValuePrices := []FairValuePriceRecord{} // Empty slice
	expectedOutput := []CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
	}

	result := DailyAndFairPriceToCombined(dailyPrices, fairValuePrices)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("DailyAndFairPriceToCombined() mismatch (-want +got):\n%s", diff)
	}
}

// Given two empty input slices, verify an empty slice is returned.
func TestDailyAndFairPriceToCombined_BothEmpty(t *testing.T) {
	dailyPrices := []DailyStockRecord{}
	fairValuePrices := []FairValuePriceRecord{}
	expectedOutput := []CombinedPriceRecord{}

	result := DailyAndFairPriceToCombined(dailyPrices, fairValuePrices)

	if diff := cmp.Diff(expectedOutput, result); diff != "" {
		t.Errorf("DailyAndFairPriceToCombined() mismatch (-want +got):\n%s", diff)
	}
}

// Given two nil input slices, verify a non-nil, empty slice is returned.
func TestDailyAndFairPriceToCombined_BothNil(t *testing.T) {
	var dailyPrices []DailyStockRecord = nil
	var fairValuePrices []FairValuePriceRecord = nil

	result := DailyAndFairPriceToCombined(dailyPrices, fairValuePrices)

	if result == nil {
		t.Fatal("DailyAndFairPriceToCombined() returned nil for nil inputs, expected empty slice")
	}
	if len(result) != 0 {
		t.Errorf("DailyAndFairPriceToCombined() expected zero-length slice for nil inputs, got %d", len(result))
	}
}
