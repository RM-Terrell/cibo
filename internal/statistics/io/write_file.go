package io

import (
	"fmt"
	"log"

	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"

	"cibo/internal/types"
)

// Writes combined price data to a Parquet file.
func WriteCombinedPriceDataToParquet(
	combinedData []types.CombinedPriceRecord,
	fw source.ParquetFile,
) error {
	combinedDataParquet := types.CombinedPricesToParquet(combinedData)
	pw, err := writer.NewParquetWriter(fw, new(types.CombinedPriceRecordParquet), 4)
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %w", err)
	}

	for _, record := range combinedDataParquet {
		if err = pw.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		return fmt.Errorf("failed to stop parquet writer: %w", err)
	}

	log.Printf("Successfully wrote %d combined records to Parquet file", len(combinedDataParquet))
	return nil
}
