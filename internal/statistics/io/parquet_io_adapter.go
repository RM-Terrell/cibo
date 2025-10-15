package io

import (
	"cibo/internal/types"
	"fmt"
	"io"

	"github.com/xitongsys/parquet-go/source"
)

type ParquetIOAdapter struct{}

func NewParquetIOAdapter() *ParquetIOAdapter {
	return &ParquetIOAdapter{}
}

// TODO some day this adapter file, and the original should be merged to follow
// the method-on-struct format seen in the API client concrete implementation
// instead of this adapter thing. That will require unit test updates too.
func (a *ParquetIOAdapter) WriteCombinedPriceDataToParquet(
	records []types.CombinedPriceRecord,
	writer io.WriteCloser,
) error {
	parquetFile, ok := writer.(source.ParquetFile)
	if !ok {
		return fmt.Errorf("writer is not a valid source.ParquetFile")
	}

	return WriteCombinedPriceDataToParquet(records, parquetFile)
}
