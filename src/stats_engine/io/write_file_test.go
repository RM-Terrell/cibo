package io

import (
	"path/filepath"
	"stats_engine/types"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func TestWriteCombinedPriceData(t *testing.T) {

	t.Run("Given daily and fair value records, verify they are correctly merged and written", func(t *testing.T) {
		dailyPrices := []types.DailyStockRecord{
			{Ticker: "TEST", Date: "2025-01-01", ClosingPrice: 100.0},
			{Ticker: "TEST", Date: "2025-01-02", ClosingPrice: 102.5},
		}
		fairValuePrices := []types.FairValuePriceRecord{
			{Ticker: "TEST", Date: "2025-12-31", FairValuePrice: 150.0},
		}

		expectedOutput := []types.CombinedPriceRecordParquet{
			{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "actual_price"},
			{Ticker: "TEST", Date: "2025-01-02", Price: 102.5, Series: "actual_price"},
			{Ticker: "TEST", Date: "2025-12-31", Price: 150.0, Series: "fair_value"},
		}

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "success.parquet")
		fw, err := local.NewLocalFileWriter(filePath)
		if err != nil {
			t.Fatalf("Failed to create file writer: %v", err)
		}
		err = WriteCombinedPriceData(dailyPrices, fairValuePrices, fw)
		if closeErr := fw.Close(); closeErr != nil {
			t.Fatalf("Failed to close file writer: %v", closeErr)
		}
		if err != nil {
			t.Fatalf("WriteCombinedPriceData returned an unexpected error: %v", err)
		}

		fr, _ := local.NewLocalFileReader(filePath)
		defer fr.Close()
		pr, _ := reader.NewParquetReader(fr, new(types.CombinedPriceRecordParquet), 4)

		if pr.GetNumRows() != int64(len(expectedOutput)) {
			t.Fatalf("Row count mismatch: want %d, got %d", len(expectedOutput), pr.GetNumRows())
		}

		readRecords := make([]types.CombinedPriceRecordParquet, len(expectedOutput))
		if err := pr.Read(&readRecords); err != nil {
			t.Fatalf("Failed to read records: %v", err)
		}

		sorter := cmpopts.SortSlices(func(a, b types.CombinedPriceRecordParquet) bool {
			if a.Date != b.Date {
				return a.Date < b.Date
			}
			return a.Series < b.Series // Differentiate if dates are the same
		})
		if diff := cmp.Diff(expectedOutput, readRecords, sorter); diff != "" {
			t.Errorf("Record mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Given only daily prices, verify only they are written", func(t *testing.T) {
		dailyPrices := []types.DailyStockRecord{
			{Ticker: "TEST", Date: "2025-01-01", ClosingPrice: 100.0},
		}
		// Distinguishing line, empty fair value data
		fairValuePrices := []types.FairValuePriceRecord{}

		expectedOutput := []types.CombinedPriceRecordParquet{
			{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "actual_price"},
		}

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "daily_only.parquet")
		fw, _ := local.NewLocalFileWriter(filePath)
		err := WriteCombinedPriceData(dailyPrices, fairValuePrices, fw)
		fw.Close()
		if err != nil {
			t.Fatalf("WriteCombinedPriceData returned an unexpected error: %v", err)
		}

		fr, _ := local.NewLocalFileReader(filePath)
		defer fr.Close()
		pr, _ := reader.NewParquetReader(fr, new(types.CombinedPriceRecordParquet), 4)

		readRecords := make([]types.CombinedPriceRecordParquet, 1)
		pr.Read(&readRecords)

		if diff := cmp.Diff(expectedOutput, readRecords); diff != "" {
			t.Errorf("Record mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Given empty or nil slices, verify an empty file is created without error", func(t *testing.T) {
		testCases := map[string]struct {
			daily     []types.DailyStockRecord
			fairValue []types.FairValuePriceRecord
		}{
			"BothEmpty": {daily: []types.DailyStockRecord{}, fairValue: []types.FairValuePriceRecord{}},
			"BothNil":   {daily: nil, fairValue: nil},
			"OneNil":    {daily: []types.DailyStockRecord{}, fairValue: nil},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tempDir := t.TempDir()
				filePath := filepath.Join(tempDir, "empty.parquet")
				fw, _ := local.NewLocalFileWriter(filePath)
				err := WriteCombinedPriceData(tc.daily, tc.fairValue, fw)
				fw.Close()

				if err != nil {
					t.Fatalf("WriteCombinedPriceData with empty/nil input returned an error: %v", err)
				}

				fr, _ := local.NewLocalFileReader(filePath)
				defer fr.Close()
				pr, _ := reader.NewParquetReader(fr, new(types.CombinedPriceRecordParquet), 4)

				if pr.GetNumRows() != 0 {
					t.Errorf("Expected 0 rows for empty/nil input, but got %d", pr.GetNumRows())
				}
			})
		}
	})
}
