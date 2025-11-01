package pipelines

import (
	"cibo/internal/statistics/algos"
	"cibo/internal/statistics/parse"
	"cibo/internal/statistics/utils"
	"cibo/internal/types"
	"fmt"
	"path/filepath"

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
	Ticker    string
	StartDate string
	EndDate   string
}

type LynchFairValueOutputs struct {
	RecordCount       int
	FilePath          string
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
	dailyPricesJson, err := p.apiClient.FetchDailyPrice(input.Ticker)
	if err != nil {
		return nil, fmt.Errorf("daily prices API fetch failed: %w", err)
	}
	annualEarningsJson, err := p.apiClient.FetchEarnings(input.Ticker)
	if err != nil {
		return nil, fmt.Errorf("annual earnings API fetch failed: %w", err)
	}
	stockSplitsJson, err := p.apiClient.FetchStockSplits((input.Ticker))
	if err != nil {
		return nil, fmt.Errorf("stock splits API fetch failed: %w", err)
	}

	dailyPricesRecords, err := parse.ParseDailyPricesToFlat(dailyPricesJson, true)
	if err != nil {
		return nil, fmt.Errorf("daily prices parsing failed: %w", err)
	}
	annualEarningsRecords, err := parse.ParseAnnualEarningsToFlat(annualEarningsJson, true)
	if err != nil {
		return nil, fmt.Errorf("annual earnings parsing failed: %w", err)
	}
	stockSplitRecords, err := parse.ParseStockSplitsToFlat(stockSplitsJson)
	if err != nil {
		return nil, fmt.Errorf("stock splits parsing failed: %w", err)
	}

	adjustedDailyPrices := utils.AdjustForStockSplits(dailyPricesRecords, stockSplitRecords)
	filteredDailyPrices, err := utils.FilterDailyPricesWithinDateRange(adjustedDailyPrices, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to filter daily prices: %w", err)
	}

	filteredAnnualEarnings, err := utils.FilterAnnualEarningsWithinDateRange(annualEarningsRecords, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to filter annual earnings: %w", err)
	}

	fairValuePriceRecords, err := algos.CalculateFairValueHistory(filteredAnnualEarnings)
	if err != nil {
		return nil, fmt.Errorf("could not calculate fair value: %w", err)
	}

	combinedData := types.DailyAndFairPriceToCombined(filteredDailyPrices, fairValuePriceRecords)

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

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for '%s': %w", fileName, err)
	}

	output := &LynchFairValueOutputs{
		RecordCount:       len(filteredDailyPrices),
		FilePath:          absPath,
		CombinedPriceData: combinedData,
		Logs:              []string{writeLogMessage},
	}

	return output, nil
}
