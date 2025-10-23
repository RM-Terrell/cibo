// internal/io/write_file.go

package io

import (
	"cibo/internal/types"
	"fmt"
	"io"

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
