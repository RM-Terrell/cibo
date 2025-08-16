package io

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
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

		expectedRows := int64(len(sampleRecords))
		if diff := cmp.Diff(expectedRows, pr.GetNumRows()); diff != "" {
			t.Fatalf("WriteToParquet() row count mismatch (-want +got):\n%s", diff)
		}

		readRecords := make([]types.FlatStockRecord, len(sampleRecords))
		if err := pr.Read(&readRecords); err != nil {
			t.Fatalf("Failed to read records from ParquetReader: %v", err)
		}

		if diff := cmp.Diff(sampleRecords, readRecords); diff != "" {
			t.Errorf("WriteToParquet() record mismatch (-want +got):\n%s", diff)
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

		fr, err := local.NewLocalFileReader(filePath)
		if err != nil {
			t.Fatalf("Failed to create local file reader for empty file: %v", err)
		}
		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(types.FlatStockRecord), 4)
		if err != nil {
			t.Fatalf("Failed to create ParquetReader for empty file: %v", err)
		}

		expectedRows := int64(0)
		if diff := cmp.Diff(expectedRows, pr.GetNumRows()); diff != "" {
			t.Errorf("WriteToParquet() empty slice row count mismatch (-want +got):\n%s", diff)
		}
	})
}
