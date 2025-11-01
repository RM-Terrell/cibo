package parse

import (
	"fmt"
	"strings"
	"testing"

	"cibo/internal/types"

	"github.com/google/go-cmp/cmp"
)

// Given valid json data, verify that the parse function returns a correctly sorted collection of prices.
func TestParseDailyToFloatHappyPath(t *testing.T) {
	const (
		ticker      = "IBM"
		date1       = "2025-08-01"
		price1Str   = "250.05"
		price1Float = 250.05
		date2       = "2025-07-31"
		price2Str   = "251.99"
		price2Float = 251.99
		date3       = "2025-08-02"
		price3Str   = "252.50"
		price3Float = 252.50
	)

	// The dates in the JSON are intentionally out of order to simulate the unordered nature of a map.
	jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" }
			}
		}`, ticker, date1, price1Str, date2, price2Str, date3, price3Str)

	// The expected result must be sorted by date, descending (newest first).
	expected := []types.DailyStockRecord{
		{Ticker: ticker, Date: date3, ClosingPrice: price3Float},
		{Ticker: ticker, Date: date1, ClosingPrice: price1Float},
		{Ticker: ticker, Date: date2, ClosingPrice: price2Float},
	}

	records, err := ParseDailyPricesToFlat(jsonData, false)
	if err != nil {
		t.Fatalf(`ParseDailyPricesToFlat() returned an unexpected error: %v`, err)
	}

	if diff := cmp.Diff(expected, records); diff != "" {
		t.Errorf("ParseDailyPricesToFlat() mismatch (-want +got):\n%s", diff)
	}
}

// Given malformed json data, verify that the parser returns an error
func TestMalformedJsonParseDailyToFloat(t *testing.T) {
	jsonData := []byte(`{ "Meta Data": "invalid }`)
	_, err := ParseDailyPricesToFlat(jsonData, false)
	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}
	if !strings.Contains(err.Error(), "unmarshaling") {
		t.Errorf("Expected error message to contain 'unmarshaling', but got: %v", err)
	}
}

// Given valid json data but with a missing ticker and strict error mode, verify that the parser returns an error
func TestParseDailyMissingTickerStrict(t *testing.T) {
	jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "1. Information": "Some info" },
			"Time Series (Daily)": {}
		}`)
	_, err := ParseDailyPricesToFlat(jsonData, false)
	if err == nil {
		t.Fatal(`Expected an error for missing ticker, but got nil`)
	}
	if !strings.Contains(err.Error(), "ticker not found") {
		t.Errorf(`
				Expected error message to contain 'ticker not found'
				Got: %v`,
			err)
	}
}

// Given data with a non-numeric closing price and strict error mode, verify that the parser returns an error
func TestParseDailyNonNumericPriceStrict(t *testing.T) {
	const (
		ticker   = "IBM"
		date     = "2025-08-01"
		badPrice = "N/A"
	)
	jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {
				"%s": { "4. close": "%s" }
			}
		}`, ticker, date, badPrice)

	_, err := ParseDailyPricesToFlat(jsonData, false)
	if err == nil {
		t.Fatal(`Expected an error for invalid closing price, but got nil`)
	}
	if !strings.Contains(err.Error(), "could not parse close price") {
		t.Errorf(`
				Expected error message to contain 'could not parse close price'
				Got: %v`,
			err)
	}
}

// Given data with a non-numeric closing price and permissive errors mode, verify that the parser returns the valid records, correctly sorted.
func TestParseDailyNonNumericPermissive(t *testing.T) {
	const (
		ticker      = "GOOD"
		goodDate1   = "2025-08-01"
		goodPrice1  = "100.50"
		price1Float = 100.50
		badDate     = "2025-07-31"
		badPrice    = "not-a-number"
		goodDate2   = "2025-07-30"
		goodPrice2  = "99.25"
		price2Float = 99.25
	)
	jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" }
			}
		}`, ticker, goodDate1, goodPrice1, badDate, badPrice, goodDate2, goodPrice2)

	expected := []types.DailyStockRecord{
		{Ticker: ticker, Date: goodDate1, ClosingPrice: price1Float},
		{Ticker: ticker, Date: goodDate2, ClosingPrice: price2Float},
	}

	records, err := ParseDailyPricesToFlat(jsonData, true)
	if err != nil {
		t.Fatalf(`Expected no error when skipping, but got: %v`, err)
	}

	if diff := cmp.Diff(expected, records); diff != "" {
		t.Errorf("ParseDailyPricesToFlat() mismatch (-want +got):\n%s", diff)
	}
}

// Given valid json data but with an empty time series, verify that the parser returns an empty collection
func TestParseDailyValidJsonEmptyTimeSeries(t *testing.T) {
	const ticker = "IBM"
	jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {}
		}`, ticker)

	records, err := ParseDailyPricesToFlat(jsonData, false)
	if err != nil {
		t.Fatalf(`Expected no error for empty time series, but got: %v`, err)
	}
	if len(records) != 0 {
		t.Errorf(`Expected 0 records for empty time series, but got %d`, len(records))
	}
}

// Given valid json, should parse all records correctly
func TestParseAnnualHappyPath(t *testing.T) {
	const (
		ticker    = "IBM"
		date1     = "2025-06-30"
		eps1Str   = "4.40"
		eps1Float = 4.40
		date2     = "2024-06-30"
		eps2Str   = "4.15"
		eps2Float = 4.15
	)

	jsonData := []byte(fmt.Sprintf(`{
			"symbol": "%s",
			"annualEarnings": [
				{
					"fiscalDateEnding": "%s",
					"reportedEPS": "%s"
				},
				{
					"fiscalDateEnding": "%s",
					"reportedEPS": "%s"
				}
			]
		}`, ticker, date1, eps1Str, date2, eps2Str))

	expected := []types.AnnualEarningRecord{
		{Ticker: ticker, FiscalDateEnding: date1, ReportedEPS: eps1Float},
		{Ticker: ticker, FiscalDateEnding: date2, ReportedEPS: eps2Float},
	}

	records, err := ParseAnnualEarningsToFlat(jsonData, false)
	if err != nil {
		t.Fatalf("ParseAnnualEarningsToFlat() returned an unexpected error: %v", err)
	}

	if diff := cmp.Diff(expected, records); diff != "" {
		t.Errorf("ParseAnnualEarningsToFlat() mismatch (-want +got):\n%s", diff)
	}
}

// Given malformed json, should return an unmarshaling error
func TestParseAnnualMalformedJson(t *testing.T) {
	jsonData := []byte(`{"symbol": "IBM", "annualEarnings": [ { "fiscalDateEnding": "2025-06-30" ]}`) // Missing closing brace

	_, err := ParseAnnualEarningsToFlat(jsonData, false)
	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}
	if !strings.Contains(err.Error(), "unmarshaling") {
		t.Errorf("Expected error message to contain 'unmarshaling', but got: %v", err)
	}
}

// Given json with a missing ticker, should return an error
func TestParseAnnualMissingTicker(t *testing.T) {
	jsonData := []byte(`{
			"annualEarnings": [
				{ "fiscalDateEnding": "2025-06-30", "reportedEPS": "4.40" }
			]
		}`)

	_, err := ParseAnnualEarningsToFlat(jsonData, false)
	if err == nil {
		t.Fatal("Expected an error for missing ticker, but got nil")
	}
	if !strings.Contains(err.Error(), "ticker not found") {
		t.Errorf("Expected error message to contain 'ticker not found', but got: %v", err)
	}
}

// Given a non-numeric EPS with strict error mode, should return a parsing error
func TestParseAnnualNonNumericStrict(t *testing.T) {
	const ticker = "IBM"
	jsonData := []byte(fmt.Sprintf(`{
			"symbol": "%s",
			"annualEarnings": [
				{ "fiscalDateEnding": "2025-06-30", "reportedEPS": "N/A" }
			]
		}`, ticker))

	_, err := ParseAnnualEarningsToFlat(jsonData, false)
	if err == nil {
		t.Fatal("Expected an error for invalid reported EPS, but got nil")
	}
	if !strings.Contains(err.Error(), "could not parse reported EPS") {
		t.Errorf("Expected error message to contain 'could not parse reported EPS', but got: %v", err)
	}
}

// Given json with an empty annualEarnings array, should return an empty slice
func TestParseAnnualEmpty(t *testing.T) {
	const ticker = "GOOG"
	jsonData := []byte(fmt.Sprintf(`{
			"symbol": "%s",
			"annualEarnings": []
		}`, ticker))

	records, err := ParseAnnualEarningsToFlat(jsonData, false)
	if err != nil {
		t.Fatalf("Expected no error for an empty earnings array, but got: %v", err)
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records for an empty earnings array, but got %d", len(records))
	}
}

// Given a non-numeric EPS with permissive error mode, should skip the bad record
func TestParseAnnualNonNumericPermissive(t *testing.T) {
	const (
		ticker    = "MSFT"
		goodDate1 = "2025-06-30"
		goodEps1  = "12.50"
		badDate   = "2024-06-30"
		goodDate2 = "2023-06-30"
		goodEps2  = "10.25"
	)
	jsonData := []byte(fmt.Sprintf(`{
			"symbol": "%s",
			"annualEarnings": [
				{ "fiscalDateEnding": "%s", "reportedEPS": "%s" },
				{ "fiscalDateEnding": "%s", "reportedEPS": "None" },
				{ "fiscalDateEnding": "%s", "reportedEPS": "%s" }
			]
		}`, ticker, goodDate1, goodEps1, badDate, goodDate2, goodEps2))

	expected := []types.AnnualEarningRecord{
		{Ticker: ticker, FiscalDateEnding: goodDate1, ReportedEPS: 12.50},
		{Ticker: ticker, FiscalDateEnding: goodDate2, ReportedEPS: 10.25},
	}

	records, err := ParseAnnualEarningsToFlat(jsonData, true)
	if err != nil {
		t.Fatalf("Expected no error when skipping bad records, but got: %v", err)
	}

	if diff := cmp.Diff(expected, records); diff != "" {
		t.Errorf("ParseAnnualEarningsToFlat() mismatch (-want +got):\n%s", diff)
	}
}

// Given valid json data, verify that it is parsed into the correct collection of split records.
func TestParseStockSplitsHappyPath(t *testing.T) {
	const (
		ticker     = "NVDA"
		date1      = "2024-06-10"
		factor1Str = "10.0"
		factor1Num = 10.0
		date2      = "2021-07-20"
		factor2Str = "4.0"
		factor2Num = 4.0
		date3      = "2025-09-02"
		factor3Str = "0.1"
		factor3Num = 0.1
	)

	jsonData := []byte(fmt.Sprintf(`{
		"symbol": "%s",
		"data": [
			{ "effective_date": "%s", "split_factor": "%s" },
			{ "effective_date": "%s", "split_factor": "%s" },
			{ "effective_date": "%s", "split_factor": "%s" }
		]
	}`, ticker, date1, factor1Str, date2, factor2Str, date3, factor3Str))

	expected := []types.StockSplitRecord{
		{Ticker: ticker, EffectiveDate: date1, SplitFactor: factor1Num},
		{Ticker: ticker, EffectiveDate: date2, SplitFactor: factor2Num},
		{Ticker: ticker, EffectiveDate: date3, SplitFactor: factor3Num},
	}

	records, err := ParseStockSplitsToFlat(jsonData)
	if err != nil {
		t.Fatalf("ParseStockSplitsToFlat() returned an unexpected error: %v", err)
	}

	if diff := cmp.Diff(expected, records); diff != "" {
		t.Errorf("ParseStockSplitsToFlat() mismatch (-want +got):\n%s", diff)
	}
}

// Given malformed json, verify that an unmarshaling error is returned.
func TestParseStockSplitsMalformedJson(t *testing.T) {
	jsonData := []byte(`{"symbol": "NVDA", "data": [ { "effective_date": "2024-06-10" ]}`) // Missing closing brace

	_, err := ParseStockSplitsToFlat(jsonData)
	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}
	if !strings.Contains(err.Error(), "unmarshaling") {
		t.Errorf("Expected error message to contain 'unmarshaling', but got: %v", err)
	}
}

// Given valid json with a missing symbol, verify that an error is returned.
func TestParseStockSplitsMissingSymbol(t *testing.T) {
	jsonData := []byte(`{
		"data": [
			{ "effective_date": "2024-06-10", "split_factor": "10.0" }
		]
	}`)

	_, err := ParseStockSplitsToFlat(jsonData)
	if err == nil {
		t.Fatal("Expected an error for missing symbol, but got nil")
	}
	if !strings.Contains(err.Error(), "ticker not found") {
		t.Errorf("Expected error message to contain 'ticker not found', but got: %v", err)
	}
}

// Given valid json with an empty data array, verify that an empty slice is returned.
func TestParseStockSplitsEmptyData(t *testing.T) {
	const ticker = "RKLB"
	jsonData := []byte(fmt.Sprintf(`{
		"symbol": "%s",
		"data": []
	}`, ticker))

	records, err := ParseStockSplitsToFlat(jsonData)
	if err != nil {
		t.Fatalf("Expected no error for an empty data array, but got: %v", err)
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records for an empty data array, but got %d", len(records))
	}
}

// Given a record with a non-numeric split factor, verify that the record is skipped and no error is returned.
func TestParseStockSplitsSkipsBadFactor(t *testing.T) {
	const (
		ticker      = "TEST"
		goodDate1   = "2024-01-01"
		goodFactor  = "2.0"
		badDate     = "2023-01-01"
		goodDate2   = "2022-01-01"
		goodFactor2 = "3.0"
	)
	jsonData := []byte(fmt.Sprintf(`{
		"symbol": "%s",
		"data": [
			{ "effective_date": "%s", "split_factor": "%s" },
			{ "effective_date": "%s", "split_factor": "two-for-one" },
			{ "effective_date": "%s", "split_factor": "%s" }
		]
	}`, ticker, goodDate1, goodFactor, badDate, goodDate2, goodFactor2))

	expected := []types.StockSplitRecord{
		{Ticker: ticker, EffectiveDate: goodDate1, SplitFactor: 2.0},
		{Ticker: ticker, EffectiveDate: goodDate2, SplitFactor: 3.0},
	}

	records, err := ParseStockSplitsToFlat(jsonData)
	if err != nil {
		t.Fatalf("Expected no error when skipping bad records, but got: %v", err)
	}

	if diff := cmp.Diff(expected, records); diff != "" {
		t.Errorf("ParseStockSplitsToFlat() mismatch (-want +got):\n%s", diff)
	}
}
