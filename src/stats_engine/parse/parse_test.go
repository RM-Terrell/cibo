package parse

import (
	"fmt"
	"strings"
	"testing"

	"stats_engine/types"
)

func TestParseToFlat(t *testing.T) {
	t.Run("Given valid json data with all required keys and values, verify that the parse function returns an expected collection of prices", func(t *testing.T) {
		const (
			ticker      = "IBM"
			date1       = "2025-08-01"
			price1Str   = "250.05"
			price1Float = 250.05
			date2       = "2025-07-31"
			price2Str   = "251.99"
			price2Float = 251.99
		)

		jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" }
			}
		}`, ticker, date1, price1Str, date2, price2Str)

		records, err := ParseToFlat(jsonData, false)
		if err != nil {
			t.Fatalf(`ParseToFlat() returned an unexpected error: %v`, err)
		}

		expectedLen := 2
		if len(records) != expectedLen {
			t.Fatalf(`
				Expected %d records
				Got %d`,
				expectedLen, len(records))
		}

		resultsMap := make(map[string]types.FlatStockRecord)
		for _, r := range records {
			resultsMap[r.Date] = r
		}

		if rec, ok := resultsMap[date1]; !ok {
			t.Errorf(`Expected record for date %s not found`, date1)
		} else {
			if rec.Ticker != ticker {
				t.Errorf(`
					Expected ticker '%s'
					Got '%s'`,
					ticker, rec.Ticker)
			}
			if rec.ClosingPrice != price1Float {
				t.Errorf(`
					Expected closing price %f
					Got %f`,
					price1Float, rec.ClosingPrice)
			}
		}

		if rec, ok := resultsMap[date2]; !ok {
			t.Errorf(`Expected record for date %s not found`, date2)
		} else {
			if rec.ClosingPrice != price2Float {
				t.Errorf(`
					Expected closing price %f
					Got %f`,
					price2Float, rec.ClosingPrice)
			}
		}
	})

	t.Run("Given malformed json data, verify that the parser returns an error", func(t *testing.T) {
		jsonData := []byte(`{ "Meta Data": "invalid }`)
		_, err := ParseToFlat(jsonData, false)
		if err == nil {
			t.Fatal(`Expected an error for malformed JSON, but got nil`)
		}
	})

	t.Run("Given valid json data but with a missing ticker and strict error mode, verify that the parser returns an error", func(t *testing.T) {
		jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "1. Information": "Some info" },
			"Time Series (Daily)": {}
		}`)
		_, err := ParseToFlat(jsonData, false)
		if err == nil {
			t.Fatal(`Expected an error for missing ticker, but got nil`)
		}
		if !strings.Contains(err.Error(), "ticker not found") {
			t.Errorf(`
				Expected error message to contain 'ticker not found'
				Got: %v`,
				err)
		}
	})

	t.Run("Given data with a non-numeric closing price and strict error mode, verify that the parser returns an error", func(t *testing.T) {
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

		// skipErrors being false is an important distinction in this test, and should result in error return
		_, err := ParseToFlat(jsonData, false)
		if err == nil {
			t.Fatal(`Expected an error for invalid closing price, but got nil`)
		}
		if !strings.Contains(err.Error(), "could not parse close price") {
			t.Errorf(`
				Expected error message to contain 'could not parse close price'
				Got: %v`,
				err)
		}
	})

	t.Run("Given data with a non-numeric closing price and permissive errors mode, verify that the parser returns the valid records only", func(t *testing.T) {
		const (
			ticker     = "GOOD"
			goodDate1  = "2025-08-01"
			goodPrice1 = "100.50"
			badDate    = "2025-07-31"
			badPrice   = "not-a-number"
			goodDate2  = "2025-07-30"
			goodPrice2 = "99.25"
		)
		jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" },
				"%s": { "4. close": "%s" }
			}
		}`, ticker, goodDate1, goodPrice1, badDate, badPrice, goodDate2, goodPrice2)

		// skipErrors being true is an important distinction in this test, and should NOT result in error return
		records, err := ParseToFlat(jsonData, true)
		if err != nil {
			t.Fatalf(`
				Expected no error when skipping
				Got: %v`,
				err)
		}

		expectedLen := 2
		if len(records) != expectedLen {
			t.Fatalf(`
				Expected %d records after skipping
				Got %d`,
				expectedLen, len(records))
		}

		resultsMap := make(map[string]types.FlatStockRecord)
		for _, r := range records {
			resultsMap[r.Date] = r
		}

		if _, ok := resultsMap[badDate]; ok {
			t.Errorf(`The invalid record for date %s should have been skipped`, badDate)
		}
		if _, ok := resultsMap[goodDate1]; !ok {
			t.Errorf(`The valid record for date %s is missing`, goodDate1)
		}
		if _, ok := resultsMap[goodDate2]; !ok {
			t.Errorf(`The valid record for date %s is missing`, goodDate2)
		}
	})

	t.Run("Given valid json data but with an empty time series, verify that the parser returns an empty collection", func(t *testing.T) {
		const ticker = "IBM"
		jsonData := fmt.Appendf(nil, `{
			"Meta Data": { "2. Symbol": "%s" },
			"Time Series (Daily)": {}
		}`, ticker)

		records, err := ParseToFlat(jsonData, false)
		if err != nil {
			t.Fatalf(`
				Expected no error for empty time series
				Got: %v`,
				err)
		}
		if len(records) != 0 {
			t.Errorf(`
				Expected 0 records for empty time series
				Got %d`,
				len(records))
		}
	})
}
