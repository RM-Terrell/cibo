package io

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"

	"cibo/internal/types"
)

// Given daily and fair value records, verify they are correctly written
func TestWriteCombinedDataHappyPath(t *testing.T) {
	combinedData := []types.CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-01-02", Price: 102.5, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-12-31", Price: 150.0, Series: "fair_value"},
	}

	expectedOutput := []types.CombinedPriceRecordParquet{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-01-02", Price: 102.5, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-12-31", Price: 150.0, Series: "fair_value"},
	}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "success.parquet")
	fw, err := local.NewLocalFileWriter(filePath)
	if err != nil {
		t.Fatalf("Failed to create file writer: %v", err)
	}

	client := NewParquetClient()
	_, err = client.WriteCombinedPriceDataToParquet(combinedData, fw)
	if err != nil {
		t.Fatalf("WriteCombinedPriceData returned an unexpected error: %v", err)
	}

	if closeErr := fw.Close(); closeErr != nil {
		t.Fatalf("Failed to close file writer: %v", closeErr)
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
}

// Given only daily prices, verify only daily are written. Defensive paranoia test in case of future data series logic mishandling
func TestOnlyDailyDefensiveWriting(t *testing.T) {
	combinedData := []types.CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"}, // note only daily_price series
	}

	expectedOutput := []types.CombinedPriceRecordParquet{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
	}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "daily_only.parquet")
	fw, _ := local.NewLocalFileWriter(filePath)

	client := NewParquetClient()
	_, err := client.WriteCombinedPriceDataToParquet(combinedData, fw)
	if err != nil {
		t.Fatalf("WriteCombinedPriceData returned an unexpected error: %v", err)
	}

	if closeErr := fw.Close(); closeErr != nil {
		t.Fatalf("Failed to close file writer: %v", closeErr)
	}

	fr, _ := local.NewLocalFileReader(filePath)
	defer fr.Close()
	pr, _ := reader.NewParquetReader(fr, new(types.CombinedPriceRecordParquet), 4)

	readRecords := make([]types.CombinedPriceRecordParquet, 1)
	pr.Read(&readRecords)

	if diff := cmp.Diff(expectedOutput, readRecords); diff != "" {
		t.Errorf("Record mismatch (-want +got):\n%s", diff)
	}
}

// Given an empty value in the daily prices, verify an empty file is created without error
func TestWriteCombineEmptyDaily(t *testing.T) {
	combinedDataEmpty := []types.CombinedPriceRecord{}

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty.parquet")
	fw, _ := local.NewLocalFileWriter(filePath)

	client := NewParquetClient()
	_, err := client.WriteCombinedPriceDataToParquet(combinedDataEmpty, fw)
	if closeErr := fw.Close(); closeErr != nil {
		t.Fatalf("Failed to close file writer: %v", closeErr)
	}

	if err != nil {
		t.Fatalf("WriteCombinedPriceData with empty input returned an error: %v", err)
	}

	fr, _ := local.NewLocalFileReader(filePath)
	defer fr.Close()
	pr, _ := reader.NewParquetReader(fr, new(types.CombinedPriceRecordParquet), 4)

	if pr.GetNumRows() != 0 {
		t.Errorf("Expected 0 rows for empty/nil input, but got %d", pr.GetNumRows())
	}
}

// todo update these
func TestReadCombinedDataHappyPath(t *testing.T) {
	// GIVEN: A set of records to write and then read back
	recordsToWrite := []types.CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-01-02", Price: 102.5, Series: "fair_value"},
	}
	expectedOutput := []types.CombinedPriceRecordParquet{
		{Ticker: "TEST", Date: "2025-01-01", Price: 100.0, Series: "daily_price"},
		{Ticker: "TEST", Date: "2025-01-02", Price: 102.5, Series: "fair_value"},
	}

	// SETUP: Write a temporary file using our writer method
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "read_test.parquet")
	fw, _ := local.NewLocalFileWriter(filePath)
	client := NewParquetClient()
	_, err := client.WriteCombinedPriceDataToParquet(recordsToWrite, fw)
	if err != nil {
		t.Fatalf("Setup failed: could not write parquet file for reading: %v", err)
	}
	fw.Close()

	// ACTION: Call the reader method
	readRecords, err := client.ReadCombinedPriceDataFromParquet(filePath)

	// VERIFY: The data is read correctly without error
	if err != nil {
		t.Fatalf("ReadCombinedPriceDataFromParquet returned an unexpected error: %v", err)
	}

	sorter := cmpopts.SortSlices(func(a, b types.CombinedPriceRecordParquet) bool { return a.Date < b.Date })
	if diff := cmp.Diff(expectedOutput, readRecords, sorter); diff != "" {
		t.Errorf("Record mismatch (-want +got):\n%s", diff)
	}
}

// Given a non-existent file path, verify an error is returned.
func TestReadDataFileNotFound(t *testing.T) {
	client := NewParquetClient()
	_, err := client.ReadCombinedPriceDataFromParquet("non_existent_file.parquet")

	if err == nil {
		t.Fatal("Expected an error when reading a non-existent file, but got nil")
	}
}

// Given an empty Parquet file, verify an empty slice is returned without error.
func TestReadEmptyFile(t *testing.T) {
	// SETUP: Write an empty file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty_read_test.parquet")
	fw, _ := local.NewLocalFileWriter(filePath)
	client := NewParquetClient()
	_, err := client.WriteCombinedPriceDataToParquet([]types.CombinedPriceRecord{}, fw)
	if err != nil {
		t.Fatalf("Setup failed: could not write empty parquet file: %v", err)
	}
	fw.Close()

	// ACTION: Read the empty file
	records, err := client.ReadCombinedPriceDataFromParquet(filePath)

	// VERIFY: No error and an empty slice
	if err != nil {
		t.Fatalf("Expected no error when reading an empty file, but got: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("Expected an empty slice when reading an empty file, but got %d records", len(records))
	}
}
