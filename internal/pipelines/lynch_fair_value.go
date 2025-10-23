package pipelines

import (
	"cibo/internal/statistics/algos"
	"cibo/internal/statistics/parse"
	"cibo/internal/types"
	"fmt"

	"github.com/xitongsys/parquet-go-source/local"
)

// LynchFairValuePipeline orchestrates the business logic for generating fair value reports via the Lynch method
// This package is intended to wrap the full data pipeline that fetches data via an api client, parses it,
// transforms it, etc and returns a final data set.

type LynchFairValuePipeline struct {
	apiClient     APIClient
	parquetWriter ParquetWriter
}

type LynchFairValueInputs struct {
	Ticker string
}

type LynchFairValueOutputs struct {
	RecordCount       int
	FileName          string
	CombinedPriceData []types.CombinedPriceRecord
	Logs              []string
}

func NewLynchFairValuePipeline(client APIClient, writer ParquetWriter) *LynchFairValuePipeline {
	return &LynchFairValuePipeline{
		apiClient:     client,
		parquetWriter: writer,
	}
}

func (p *LynchFairValuePipeline) RunPipeline(input LynchFairValueInputs) (*LynchFairValueOutputs, error) {
	// todo pass in date ranges here
	dailyPricesJson, err := p.apiClient.FetchDailyPrice(input.Ticker)
	if err != nil {
		return nil, fmt.Errorf("daily prices API fetch failed: %w", err)
	}

	// todo pass in date ranges here
	annualEarningsJson, err := p.apiClient.FetchEarnings(input.Ticker)
	if err != nil {
		return nil, fmt.Errorf("annual earnings API fetch failed: %w", err)
	}

	dailyPricesRecords, err := parse.ParseDailyPricesToFlat(dailyPricesJson, true)
	if err != nil {
		return nil, fmt.Errorf("daily prices parsing failed: %w", err)
	}

	annualEarningsRecords, err := parse.ParseAnnualEarningsToFlat(annualEarningsJson, true)
	if err != nil {
		return nil, fmt.Errorf("annual earnings parsing failed: %w", err)
	}

	fairValuePriceRecords, err := algos.CalculateFairValueHistory(annualEarningsRecords)
	if err != nil {
		return nil, fmt.Errorf("could not calculate fair value: %w", err)
	}

	combinedData := types.DailyAndFairPriceToCombined(dailyPricesRecords, fairValuePriceRecords)

	fileName := fmt.Sprintf("%s.parquet", input.Ticker)
	fw, err := local.NewLocalFileWriter(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file '%s': %w", fileName, err)
	}
	defer fw.Close()

	writeLogMessage, err := p.parquetWriter.WriteCombinedPriceDataToParquet(combinedData, fw)
	if err != nil {
		return nil, fmt.Errorf("failed to write parquet data: %w", err)
	}

	output := &LynchFairValueOutputs{
		RecordCount:       len(dailyPricesRecords),
		FileName:          fileName,
		CombinedPriceData: combinedData,
		Logs:              []string{writeLogMessage},
	}

	return output, nil
}
