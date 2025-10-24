package io

import (
	"cibo/internal/types"
	"fmt"
	"io"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

type ParquetClient struct{}

func NewParquetClient() *ParquetClient {
	return &ParquetClient{}
}

// Write price data to a parquet file
func (p *ParquetClient) WriteCombinedPriceDataToParquet(
	combinedData []types.CombinedPriceRecord,
	w io.WriteCloser,
) (string, error) {
	fw, ok := w.(source.ParquetFile)
	if !ok {
		return "", fmt.Errorf("writer is not a valid source.ParquetFile")
	}

	combinedDataParquet := types.CombinedPricesToParquet(combinedData)
	pw, err := writer.NewParquetWriter(fw, new(types.CombinedPriceRecordParquet), 4)
	if err != nil {
		return "", fmt.Errorf("failed to create parquet writer: %w", err)
	}

	for _, record := range combinedDataParquet {
		if err = pw.Write(record); err != nil {
			return "", fmt.Errorf("failed to write record: %w", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		return "", fmt.Errorf("failed to stop parquet writer: %w", err)
	}

	successMessage := fmt.Sprintf("Successfully wrote %d combined records to Parquet file", len(combinedDataParquet))
	return successMessage, nil
}

// Read price data from a parquet file.
func (p *ParquetClient) ReadCombinedPriceDataFromParquet(filePath string) ([]types.CombinedPriceRecordParquet, error) {
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(types.CombinedPriceRecordParquet), 4)
	if err != nil {
		return nil, fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer pr.ReadStop()

	numRecords := int(pr.GetNumRows())
	records := make([]types.CombinedPriceRecordParquet, numRecords)

	if numRecords == 0 {
		return records, nil // Return empty slice for empty file
	}

	if err := pr.Read(&records); err != nil {
		return nil, fmt.Errorf("failed to read records from parquet file: %w", err)
	}

	return records, nil
}
