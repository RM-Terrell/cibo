package utils

import (
	"cibo/internal/types"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var testDailyPrices = []types.DailyStockRecord{
	{Ticker: "TEST", Date: "2023-12-31", ClosingPrice: 100.0},
	{Ticker: "TEST", Date: "2024-01-01", ClosingPrice: 101.0},
	{Ticker: "TEST", Date: "2024-06-15", ClosingPrice: 102.0},
	{Ticker: "TEST", Date: "2024-12-31", ClosingPrice: 103.0},
	{Ticker: "TEST", Date: "2025-01-01", ClosingPrice: 104.0},
}

var testAnnualEarnings = []types.AnnualEarningRecord{
	{Ticker: "TEST", FiscalDateEnding: "2022-12-31", ReportedEPS: 8.0},
	{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 9.0},
	{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0},
}

// Given a valid start and end date, verify that only records within that range are returned.
func TestFilterDailyPricesHappyPath(t *testing.T) {
	startDate, endDate := "2024-01-01", "2024-12-31"
	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-01-01", ClosingPrice: 101.0},
		{Ticker: "TEST", Date: "2024-06-15", ClosingPrice: 102.0},
		{Ticker: "TEST", Date: "2024-12-31", ClosingPrice: 103.0},
	}

	result, err := FilterDailyPricesWithinDateRange(testDailyPrices, startDate, endDate)
	if err != nil {
		t.Fatalf("Function returned an unexpected error: %v", err)
	}

	sorter := cmpopts.SortSlices(func(a, b types.DailyStockRecord) bool { return a.Date < b.Date })
	if diff := cmp.Diff(expected, result, sorter); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given only a start date, verify that all records from that date forward are returned.
func TestFilterDailyPricesOnlyStartDate(t *testing.T) {
	startDate := "2024-06-15"
	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-06-15", ClosingPrice: 102.0},
		{Ticker: "TEST", Date: "2024-12-31", ClosingPrice: 103.0},
		{Ticker: "TEST", Date: "2025-01-01", ClosingPrice: 104.0},
	}

	result, err := FilterDailyPricesWithinDateRange(testDailyPrices, startDate, "")
	if err != nil {
		t.Fatalf("Function returned an unexpected error: %v", err)
	}

	sorter := cmpopts.SortSlices(func(a, b types.DailyStockRecord) bool { return a.Date < b.Date })
	if diff := cmp.Diff(expected, result, sorter); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given only an end date, verify that all records up to and including that date are returned.
func TestFilterDailyPricesOnlyEndDate(t *testing.T) {
	endDate := "2024-06-15"
	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2023-12-31", ClosingPrice: 100.0},
		{Ticker: "TEST", Date: "2024-01-01", ClosingPrice: 101.0},
		{Ticker: "TEST", Date: "2024-06-15", ClosingPrice: 102.0},
	}

	result, err := FilterDailyPricesWithinDateRange(testDailyPrices, "", endDate)
	if err != nil {
		t.Fatalf("Function returned an unexpected error: %v", err)
	}

	sorter := cmpopts.SortSlices(func(a, b types.DailyStockRecord) bool { return a.Date < b.Date })
	if diff := cmp.Diff(expected, result, sorter); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given a date range that does not include any records, verify that an empty slice is returned.
func TestFilterDailyPricesNoResults(t *testing.T) {
	startDate, endDate := "2022-01-01", "2022-12-31"
	expected := []types.DailyStockRecord{}

	result, err := FilterDailyPricesWithinDateRange(testDailyPrices, startDate, endDate)
	if err != nil {
		t.Fatalf("Function returned an unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected an empty slice, but got %d records", len(result))
	}
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given an invalid start date format, verify that an error is returned.
func TestFilterDailyPricesInvalidStartDate(t *testing.T) {
	_, err := FilterDailyPricesWithinDateRange(testDailyPrices, "2024/01/01", "2024-12-31")
	if err == nil {
		t.Fatal("Expected an error for invalid start date format, but got nil")
	}
}

// Given an invalid end date format, verify that an error is returned.
func TestFilterDailyPricesInvalidEndDate(t *testing.T) {
	_, err := FilterDailyPricesWithinDateRange(testDailyPrices, "2024-01-01", "2024/12/31")
	if err == nil {
		t.Fatal("Expected an error for invalid end date format, but got nil")
	}
}

// Given an empty slice of records, verify that an empty slice is returned without error.
func TestFilterDailyPricesEmptyInput(t *testing.T) {
	emptyRecords := []types.DailyStockRecord{}
	result, err := FilterDailyPricesWithinDateRange(emptyRecords, "2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatalf("Function returned an unexpected error for empty input: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected an empty slice for empty input, but got %d records", len(result))
	}
}

// Given a slice of daily prices containing a record with a malformed date,
// verify that the record is skipped and the valid records are returned correctly.
func TestFilterDailyPricesSkipsMalformedRecordDate(t *testing.T) {
	recordsWithBadData := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-01-01", ClosingPrice: 101.0},
		{Ticker: "TEST", Date: "not-a-date", ClosingPrice: 999.0}, // simulating data quality issue
		{Ticker: "TEST", Date: "2024-01-03", ClosingPrice: 103.0},
	}

	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-01-01", ClosingPrice: 101.0},
		{Ticker: "TEST", Date: "2024-01-03", ClosingPrice: 103.0},
	}

	result, err := FilterDailyPricesWithinDateRange(recordsWithBadData, "2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatalf("Function returned an unexpected error when it should have skipped bad data: %v", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given a valid start and end date, verify that only annual earnings within that range are returned.
func TestFilterAnnualEarningsHappyPath(t *testing.T) {
	startDate, endDate := "2023-01-01", "2024-12-31"
	expected := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 9.0},
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0},
	}

	result, err := FilterAnnualEarningsWithinDateRange(testAnnualEarnings, startDate, endDate)
	if err != nil {
		t.Fatalf("Function returned an unexpected error: %v", err)
	}

	sorter := cmpopts.SortSlices(func(a, b types.AnnualEarningRecord) bool { return a.FiscalDateEnding < b.FiscalDateEnding })
	if diff := cmp.Diff(expected, result, sorter); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given no start or end dates, verify that all original annual earnings records are returned.
func TestFilterAnnualEarningsNoDates(t *testing.T) {
	expected := testAnnualEarnings

	result, err := FilterAnnualEarningsWithinDateRange(testAnnualEarnings, "", "")
	if err != nil {
		t.Fatalf("Function returned an unexpected error: %v", err)
	}

	sorter := cmpopts.SortSlices(func(a, b types.AnnualEarningRecord) bool { return a.FiscalDateEnding < b.FiscalDateEnding })
	if diff := cmp.Diff(expected, result, sorter); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}

// Given an invalid end date format, verify that an error is returned.
func TestFilterAnnualEarningsInvalidEndDate(t *testing.T) {
	_, err := FilterAnnualEarningsWithinDateRange(testAnnualEarnings, "2023-01-01", "not-a-date")

	if err == nil {
		t.Fatal("Expected an error for invalid end date format, but got nil")
	}
}

// Given an invalid start date format, verify that an error is returned.
func TestFilterAnnualEarningsInvalidStartDate(t *testing.T) {
	_, err := FilterAnnualEarningsWithinDateRange(testAnnualEarnings, "not-a-date", "2024-12-31")
	if err == nil {
		t.Fatal("Expected an error for invalid start date format, but got nil")
	}
}

// Given a slice of annual earnings containing a record with a malformed date,
// verify that the record is skipped and the valid records are returned correctly.
func TestFilterAnnualEarningsSkipsMalformedRecordDate(t *testing.T) {
	recordsWithBadData := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 9.0},
		{Ticker: "TEST", FiscalDateEnding: "Jan 1st 2024", ReportedEPS: 999.0}, // simulating data quality issue
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0},
	}

	expected := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 9.0},
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0},
	}

	result, err := FilterAnnualEarningsWithinDateRange(recordsWithBadData, "2023-01-01", "2025-12-31")
	if err != nil {
		t.Fatalf("Function returned an unexpected error when it should have skipped bad data: %v", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Filtered records mismatch (-want +got):\n%s", diff)
	}
}
