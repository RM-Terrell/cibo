package algos

import (
	"testing"

	"cibo/internal/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// --- Test Data and Mocks ---

// A common set of earnings data used across multiple tests.
// It is intentionally unsorted and includes negative early-year earnings.
var mockEarnings = []types.AnnualEarningRecord{
	{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0},
	{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 8.0},
	{Ticker: "TEST", FiscalDateEnding: "2021-12-31", ReportedEPS: 5.0},
	{Ticker: "TEST", FiscalDateEnding: "2020-12-31", ReportedEPS: 2.5},
	{Ticker: "TEST", FiscalDateEnding: "2019-12-31", ReportedEPS: 1.0}, // This should be the start based on date
	{Ticker: "TEST", FiscalDateEnding: "2018-12-31", ReportedEPS: -0.5},
	{Ticker: "TEST", FiscalDateEnding: "2017-12-31", ReportedEPS: -1.0},
	{Ticker: "TEST", FiscalDateEnding: "2022-12-31", ReportedEPS: 6.0},
}

// --- Tests for Graham Lynch equations ---

// Given a valid slice of earnings with negative values, verify that the correct start (first positive) and end points are found.
func TestEarningsEndpoints_Success(t *testing.T) {
	expectedStart := types.AnnualEarningRecord{Ticker: "TEST", FiscalDateEnding: "2019-12-31", ReportedEPS: 1.0}
	expectedEnd := types.AnnualEarningRecord{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0}

	start, end, err := ProfitableEarningsStartingAndEnding(mockEarnings)

	if err != nil {
		t.Fatalf("EarningsEndpoints() returned an unexpected error: %v", err)
	}
	if diff := cmp.Diff(expectedStart, start); diff != "" {
		t.Errorf("EarningsEndpoints() start mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(expectedEnd, end); diff != "" {
		t.Errorf("EarningsEndpoints() end mismatch (-want +got):\n%s", diff)
	}
}

// Given a slice with fewer than the minimum required earnings, verify that an error is returned.
func TestEarningsEndpoints_NotEnoughData(t *testing.T) {
	earnings := []types.AnnualEarningRecord{{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 10.0}}
	_, _, err := ProfitableEarningsStartingAndEnding(earnings)
	if err == nil {
		t.Fatal("EarningsEndpoints() expected an error for insufficient data, but got none")
	}
}

// Given a slice containing only negative earnings, verify that an error is returned because there is
// no valid earnings data for an EPS calculation
func TestEarningsEndpoints_AllNegativeEPS(t *testing.T) {
	earnings := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: -1.0},
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: -0.5},
	}
	_, _, err := ProfitableEarningsStartingAndEnding(earnings)
	if err == nil {
		t.Fatal("EarningsEndpoints() expected an error for all negative EPS, but got none")
	}
}

// Given a slice where the final chronological earning is negative, verify that an error is returned.
func TestEarningsEndpoints_NegativeEndingEPS(t *testing.T) {
	earnings := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: -1.0},
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 5.0},
	}
	_, _, err := ProfitableEarningsStartingAndEnding(earnings)
	if err == nil {
		t.Fatal("EarningsEndpoints() expected an error for negative ending EPS, but got none")
	}
}

// Given a valid slice of earnings, verify that the Compound Annual Growth Rate is calculated correctly.
func TestCAGR_Success(t *testing.T) {
	// See source code for CAGR equation
	expectedCAGR := 0.58459
	tolerance := 0.0001

	cagr, err := CAGR(mockEarnings)

	if err != nil {
		t.Fatalf("CAGR() returned an unexpected error: %v", err)
	}
	if diff := cmp.Diff(expectedCAGR, cagr, cmpopts.EquateApprox(0, tolerance)); diff != "" {
		t.Errorf("CAGR() mismatch (-want +got):\n%s", diff)
	}
}

// Given a slice of earnings that is too short, verify that CAGR propagates the error from EarningsEndpoints.
func TestCAGR_ErrorPropagation(t *testing.T) {
	earnings := []types.AnnualEarningRecord{}
	_, err := CAGR(earnings)
	if err == nil {
		t.Fatal("CAGR() expected an error for insufficient data, but got none")
	}
}

// Given a slice where the start and end dates are less than a year apart, verify that an error is returned.
func TestCAGR_PeriodTooShort(t *testing.T) {
	earnings := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2024-01-01", ReportedEPS: 1.0},
		{Ticker: "TEST", FiscalDateEnding: "2024-06-01", ReportedEPS: 2.0},
	}
	_, err := CAGR(earnings)
	if err == nil {
		t.Fatal("CAGR() expected an error for a period less than one year, but got none")
	}
}

// Given a positive CAGR, verify that the correct Fair Value P/E ratio is calculated.
func TestFairValuePE_PositiveCAGR(t *testing.T) {
	cagr := 0.15 // 15% growth
	expectedPE := 15.0

	pe := FairValuePE(cagr)

	if diff := cmp.Diff(expectedPE, pe); diff != "" {
		t.Errorf("FairValuePE() mismatch (-want +got):\n%s", diff)
	}
}

// Given a zero CAGR, verify that the Fair Value P/E ratio is zero.
func TestFairValuePE_ZeroCAGR(t *testing.T) {
	cagr := 0.0
	expectedPE := 0.0

	pe := FairValuePE(cagr)

	if diff := cmp.Diff(expectedPE, pe); diff != "" {
		t.Errorf("FairValuePE() mismatch (-want +got):\n%s", diff)
	}
}

// Given a fair value P/E and a history of earnings, verify that a correct history of prices is generated, ignoring negative EPS years.
func TestFairValuePriceHistory_Success(t *testing.T) {
	fairValuePE := 20.0
	earnings := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: 2.0},
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: 1.5},
		{Ticker: "TEST", FiscalDateEnding: "2022-12-31", ReportedEPS: -0.5}, // Should be ignored
		{Ticker: "TEST", FiscalDateEnding: "2021-12-31", ReportedEPS: 1.0},
	}

	expectedHistory := []types.FairValuePriceRecord{
		{Ticker: "TEST", Date: "2021-12-31", FairValuePrice: 20.0}, // 1.0 * 20
		{Ticker: "TEST", Date: "2023-12-31", FairValuePrice: 30.0}, // 1.5 * 20
		{Ticker: "TEST", Date: "2024-12-31", FairValuePrice: 40.0}, // 2.0 * 20
	}

	history := FairValuePriceHistory(fairValuePE, earnings)

	// The function doesn't guarantee order, so we sort both slices before comparing
	sorter := func(a, b types.FairValuePriceRecord) bool {
		return a.Date < b.Date
	}
	if diff := cmp.Diff(expectedHistory, history, cmpopts.SortSlices(sorter)); diff != "" {
		t.Errorf("FairValuePriceHistory() mismatch (-want +got):\n%s", diff)
	}
}

// Given a history of only negative earnings, verify that an empty slice is returned.
func TestFairValuePriceHistory_AllNegativeEarnings(t *testing.T) {
	fairValuePE := 20.0
	earnings := []types.AnnualEarningRecord{
		{Ticker: "TEST", FiscalDateEnding: "2024-12-31", ReportedEPS: -2.0},
		{Ticker: "TEST", FiscalDateEnding: "2023-12-31", ReportedEPS: -1.5},
	}

	history := FairValuePriceHistory(fairValuePE, earnings)

	if len(history) != 0 {
		t.Errorf("FairValuePriceHistory() expected an empty slice, but got %d elements", len(history))
	}
}

// Given an empty earnings history, verify that an empty slice is returned.
func TestFairValuePriceHistory_EmptyInput(t *testing.T) {
	fairValuePE := 20.0
	earnings := []types.AnnualEarningRecord{}

	history := FairValuePriceHistory(fairValuePE, earnings)

	if len(history) != 0 {
		t.Errorf("FairValuePriceHistory() expected an empty slice for empty input, but got %d elements", len(history))
	}
}
