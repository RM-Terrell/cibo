package io

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"

	"stats_engine/types"
)

func TestWriteToParquet(t *testing.T) {
	t.Run("Given the parquet writer is given multiple records, verify it handles them correctly.", func(t *testing.T) {
		sampleRecords := []types.FlatStockRecord{
			{Ticker: "GOOGL", ClosingPrice: 150.75},
			{Ticker: "AAPL", ClosingPrice: 175.25},
		}

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_stocks.parquet")

		fw, err := local.NewLocalFileWriter(filePath)
		if err != nil {
			t.Fatalf("Failed to create local file writer: %v", err)
		}

		writeErr := WriteToParquet(sampleRecords, fw)

		if closeErr := fw.Close(); closeErr != nil {
			t.Fatalf("Failed to close file writer: %v", closeErr)
		}

		if writeErr != nil {
			t.Fatalf("WriteToParquet returned an unexpected error: %v", writeErr)
		}

		fr, err := local.NewLocalFileReader(filePath)
		if err != nil {
			t.Fatalf("Failed to create local file reader: %v", err)
		}
		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(types.FlatStockRecord), 4)
		if err != nil {
			t.Fatalf("Failed to create ParquetReader: %v", err)
		}

		if pr.GetNumRows() != int64(len(sampleRecords)) {
			t.Fatalf("Expected %d records, but file has %d", len(sampleRecords), pr.GetNumRows())
		}

		readRecords := make([]types.FlatStockRecord, len(sampleRecords))
		if err := pr.Read(&readRecords); err != nil {
			t.Fatalf("Failed to read records from ParquetReader: %v", err)
		}

		if !reflect.DeepEqual(sampleRecords, readRecords) {
			t.Errorf("Written records do not match read records.\nWant: %+v\nGot:  %+v", sampleRecords, readRecords)
		}
	})

	t.Run("Given the parquet writer is given and empty slice, verify it handles it without error.", func(t *testing.T) {
		sampleRecords := []types.FlatStockRecord{}

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_empty.parquet")

		fw, err := local.NewLocalFileWriter(filePath)
		if err != nil {
			t.Fatalf("Failed to create local file writer: %v", err)
		}

		writeErr := WriteToParquet(sampleRecords, fw)

		if closeErr := fw.Close(); closeErr != nil {
			t.Fatalf("Failed to close file writer: %v", closeErr)
		}

		if writeErr != nil {
			t.Fatalf("WriteToParquet with empty slice returned an error: %v", writeErr)
		}

		// Expecting to read zero records
		fr, err := local.NewLocalFileReader(filePath)
		if err != nil {
			t.Fatalf("Failed to create local file reader for empty file: %v", err)
		}
		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(types.FlatStockRecord), 4)
		if err != nil {
			t.Fatalf("Failed to create ParquetReader for empty file: %v", err)
		}

		if pr.GetNumRows() != 0 {
			t.Errorf("Expected 0 rows for empty input, but got %d", pr.GetNumRows())
		}
	})
}
