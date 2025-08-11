package io

import (
	"log"

	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"

	"stats_engine/types"
)

// Function to write stock data to a parquet file
func WriteToParquet(records []types.FlatStockRecord, fw source.ParquetFile) error {
	pw, err := writer.NewParquetWriter(fw, new(types.FlatStockRecord), 4)
	if err != nil {
		return err
	}

	for _, record := range records {
		if err = pw.Write(record); err != nil {
			return err
		}
	}

	if err = pw.WriteStop(); err != nil {
		return err
	}

	log.Printf("Successfully wrote %d records", len(records))
	return nil
}
