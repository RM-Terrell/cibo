package utils

import (
	"cibo/internal/types"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// --- Test Data ---

// Base price data, pre-sorted from newest to oldest as the function expects.
var sortedDailyPrices = []types.DailyStockRecord{
	{Ticker: "TEST", Date: "2024-07-05", ClosingPrice: 150.0},
	{Ticker: "TEST", Date: "2024-07-04", ClosingPrice: 148.0},
	{Ticker: "TEST", Date: "2024-07-03", ClosingPrice: 146.0}, // A 2-for-1 split happens after this day's close
	{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 288.0}, // Pre-split
	{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 284.0}, // Pre-split
}

// Given a simple 2-for-1 stock split (SplitFactor: 2.0), verify that all prices after the split
// date are halved.
func TestAdjustForStockSplits_SingleSplit(t *testing.T) {
	splits := []types.StockSplitRecord{
		{Ticker: "TEST", EffectiveDate: "2024-07-03", SplitFactor: 2.0},
	}

	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-07-05", ClosingPrice: 150.0},
		{Ticker: "TEST", Date: "2024-07-04", ClosingPrice: 148.0},
		{Ticker: "TEST", Date: "2024-07-03", ClosingPrice: 146.0},
		{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 144.0}, // Adjusted by 2.0
		{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 142.0}, // Adjusted by 2.0
	}

	result := AdjustForStockSplits(sortedDailyPrices, splits)

	if diff := cmp.Diff(expected, result, cmpopts.EquateApprox(0.001, 0)); diff != "" {
		t.Errorf("Adjusted prices mismatch (-want +got):\n%s", diff)
	}
}

// Given multiple stock splits, verify that prices are adjusted by the cumulative factor.
func TestAdjustForStockSplits_MultipleSplits(t *testing.T) {
	prices := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-07-03", ClosingPrice: 150.0},
		{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 148.0}, // 2-for-1 split after this
		{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 290.0}, // Pre-split price
		{Ticker: "TEST", Date: "2023-12-31", ClosingPrice: 280.0}, // 3-for-1 split after this
		{Ticker: "TEST", Date: "2023-12-30", ClosingPrice: 810.0}, // Pre-split price
	}
	splits := []types.StockSplitRecord{
		{Ticker: "TEST", EffectiveDate: "2024-07-02", SplitFactor: 2.0},
		{Ticker: "TEST", EffectiveDate: "2023-12-31", SplitFactor: 3.0},
	}

	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-07-03", ClosingPrice: 150.0},
		{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 148.0},
		{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 145.0}, // Adjusted by 2.0
		{Ticker: "TEST", Date: "2023-12-31", ClosingPrice: 140.0}, // Adjusted by 2.0
		{Ticker: "TEST", Date: "2023-12-30", ClosingPrice: 135.0}, // Adjusted by 2.0 * 3.0 = 6.0
	}

	result := AdjustForStockSplits(prices, splits)

	if diff := cmp.Diff(expected, result, cmpopts.EquateApprox(0.001, 0)); diff != "" {
		t.Errorf("Adjusted prices mismatch (-want +got):\n%s", diff)
	}
}

// Given a 1-for-10 reverse split reverse stock split (SplitFactor: .1),
// verify that historical prices are multiplied accordingly.
func TestAdjustForStockSplits_ReverseSplit(t *testing.T) {
	splits := []types.StockSplitRecord{
		{Ticker: "TEST", EffectiveDate: "2024-07-03", SplitFactor: 0.1},
	}

	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-07-05", ClosingPrice: 150.0},
		{Ticker: "TEST", Date: "2024-07-04", ClosingPrice: 148.0},
		{Ticker: "TEST", Date: "2024-07-03", ClosingPrice: 146.0},
		{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 2880.0}, // Adjusted from 288.0 (288 / 0.1)
		{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 2840.0}, // Adjusted from 284.0 (284 / 0.1)
	}

	result := AdjustForStockSplits(sortedDailyPrices, splits)

	if diff := cmp.Diff(expected, result, cmpopts.EquateApprox(0.001, 0)); diff != "" {
		t.Errorf("Adjusted prices mismatch (-want +got):\n%s", diff)
	}
}

// Given no stock splits, verify that the function returns an identical slice of prices.
func TestAdjustForStockSplits_NoSplits(t *testing.T) {
	result := AdjustForStockSplits(sortedDailyPrices, []types.StockSplitRecord{})

	if diff := cmp.Diff(sortedDailyPrices, result); diff != "" {
		t.Errorf("Prices should be unchanged when there are no splits (-want +got):\n%s", diff)
	}

	// Verify it's a new slice, not a pointer to the original.
	if &sortedDailyPrices[0] == &result[0] {
		t.Error("Function should return a new slice, not modify the original")
	}
}

// Given an empty slice of daily prices, verify that an empty slice is returned.
func TestAdjustForStockSplits_EmptyPrices(t *testing.T) {
	splits := []types.StockSplitRecord{
		{Ticker: "TEST", EffectiveDate: "2024-07-03", SplitFactor: 2.0},
	}
	result := AdjustForStockSplits([]types.DailyStockRecord{}, splits)

	if len(result) != 0 {
		t.Errorf("Expected an empty slice for empty input, but got %d records", len(result))
	}
}

// Given a split that occurs on the same day as a price record, verify that the price on that day is not adjusted.
func TestAdjustForStockSplits_SplitOnTradingDay(t *testing.T) {
	prices := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 150.0},
		{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 300.0}, // This is the day of the split.
		{Ticker: "TEST", Date: "2024-06-30", ClosingPrice: 296.0},
	}
	splits := []types.StockSplitRecord{
		{Ticker: "TEST", EffectiveDate: "2024-07-01", SplitFactor: 2.0},
	}

	expected := []types.DailyStockRecord{
		{Ticker: "TEST", Date: "2024-07-02", ClosingPrice: 150.0},
		{Ticker: "TEST", Date: "2024-07-01", ClosingPrice: 300.0}, // Price on split day is NOT adjusted.
		{Ticker: "TEST", Date: "2024-06-30", ClosingPrice: 148.0}, // Price before split day IS adjusted.
	}

	result := AdjustForStockSplits(prices, splits)

	if diff := cmp.Diff(expected, result, cmpopts.EquateApprox(0.001, 0)); diff != "" {
		t.Errorf("Adjusted prices mismatch (-want +got):\n%s", diff)
	}
}
