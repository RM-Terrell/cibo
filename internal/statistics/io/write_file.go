package io

import (
	"fmt"
	"log"

	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"

	"cibo/internal/types"
)

// Merges daily stock prices and annual fair value prices
// into a single long-format structure and writes them to a Parquet file.
func WriteCombinedPriceData(
	dailyPrices []types.DailyStockRecord,
	fairValuePrices []types.FairValuePriceRecord,
	fw source.ParquetFile,
) error {

	combinedData := make([]types.CombinedPriceRecordParquet, 0, len(dailyPrices)+len(fairValuePrices))

	for _, record := range dailyPrices {
		combinedData = append(combinedData, types.CombinedPriceRecordParquet{
			Ticker: record.Ticker,
			Date:   record.Date,
			Price:  record.ClosingPrice,
			Series: "actual_price",
		})
	}

	for _, record := range fairValuePrices {
		combinedData = append(combinedData, types.CombinedPriceRecordParquet{
			Ticker: record.Ticker,
			Date:   record.Date,
			Price:  record.FairValuePrice,
			Series: "fair_value",
		})
	}

	pw, err := writer.NewParquetWriter(fw, new(types.CombinedPriceRecordParquet), 4)
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %w", err)
	}

	for _, record := range combinedData {
		if err = pw.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		return fmt.Errorf("failed to stop parquet writer: %w", err)
	}

	log.Printf("Successfully wrote %d combined records to Parquet file", len(combinedData))
	return nil
}
