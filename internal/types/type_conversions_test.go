package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestCombinedPricesToParquet(t *testing.T) {
	t.Run("Given a slice of records, verify it is converted correctly", func(t *testing.T) {
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
	})

	t.Run("Given an empty slice, verify an empty slice is returned", func(t *testing.T) {
		inputRecords := []CombinedPriceRecord{}
		expectedOutput := []CombinedPriceRecordParquet{}

		result := CombinedPricesToParquet(inputRecords)

		if diff := cmp.Diff(expectedOutput, result); diff != "" {
			t.Errorf("CombinedPricesToParquet() mismatch for empty slice (-want +got):\n%s", diff)
		}
	})

	t.Run("Given a nil slice, verify a non-nil empty slice is returned", func(t *testing.T) {
		var inputRecords []CombinedPriceRecord = nil

		result := CombinedPricesToParquet(inputRecords)

		if result == nil {
			t.Fatal("CombinedPricesToParquet() returned a nil slice for nil input, but expected an empty slice")
		}
		if len(result) != 0 {
			t.Errorf("CombinedPricesToParquet() expected a zero-length slice for nil input, but got %d elements", len(result))
		}
	})

	t.Run("Given a record with a negative price, verify the value is preserved", func(t *testing.T) {
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
	})
}
