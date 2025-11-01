package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"

	"cibo/internal/types"
)

/*
Parsing module to handle data structures coming back from AlphaVantage.
Hard tied to their data structures, any changes to their API will
require changes here too.
*/

type DailyPricesResponse struct {
	MetaData   MetaDataContainer         `json:"Meta Data"`
	TimeSeries map[string]DailyDataPoint `json:"Time Series (Daily)"`
}

type MetaDataContainer struct {
	Symbol string `json:"2. Symbol"`
}

type DailyDataPoint struct {
	Close string `json:"4. close"`
}

/*
Function to take json data of daily prices and parse it into a collection
of individual stock prices.
*/
func ParseDailyPricesToFlat(jsonData []byte, skipErrors bool) ([]types.DailyStockRecord, error) {
	var response DailyPricesResponse

	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %w", err)
	}

	ticker := response.MetaData.Symbol
	if ticker == "" {
		return nil, fmt.Errorf("ticker not found in JSON Meta Data when parsing")
	}

	records := make([]types.DailyStockRecord, 0, len(response.TimeSeries))

	for rawDate, rawDataPoint := range response.TimeSeries {
		closingPrice, err := strconv.ParseFloat(rawDataPoint.Close, 64)
		if err != nil {
			if skipErrors {
				log.Printf("Warning: could not parse close price for date %s, skipping record. Error: %v", rawDate, err)
				continue
			}
			return nil, fmt.Errorf("could not parse close price '%s' for date %s: %w", rawDataPoint.Close, rawDate, err)
		}

		records = append(records, types.DailyStockRecord{
			Ticker:       ticker,
			Date:         rawDate,
			ClosingPrice: closingPrice,
		})
	}

	// Sort the records by date, descending (newest to oldest), to restore the order
	// that was lost when unmarshaling into a map.
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date > records[j].Date
	})

	return records, nil
}

type AnnualEarningResponse struct {
	Symbol         string          `json:"symbol"`
	AnnualEarnings []AnnualEarning `json:"annualEarnings"`
}

type AnnualEarning struct {
	FiscalDateEnding string `json:"fiscalDateEnding"`
	ReportedEPS      string `json:"reportedEPS"`
	EstimatedEPS     string `json:"estimatedEPS"`
	Surprise         string `json:"surprise"`
}

/*
Function to take json data of annual earnings and parse it into a collection
of individual annual earnings data points.
*/
func ParseAnnualEarningsToFlat(jsonData []byte, skipErrors bool) ([]types.AnnualEarningRecord, error) {
	var response AnnualEarningResponse

	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %w", err)
	}

	ticker := response.Symbol
	if ticker == "" {
		return nil, fmt.Errorf("ticker not found in JSON when parsing")
	}

	records := make([]types.AnnualEarningRecord, 0, len(response.AnnualEarnings))

	for _, earnings := range response.AnnualEarnings {
		fiscalDateEnding := earnings.FiscalDateEnding
		reportedEPS, epsParseError := strconv.ParseFloat(earnings.ReportedEPS, 64)
		if epsParseError != nil {
			if skipErrors {
				log.Printf("Warning: could not parse reported EPS for date %s, skipping record. Error: %v",
					fiscalDateEnding, epsParseError)
				continue
			}
			return nil, fmt.Errorf("could not parse reported EPS for date %s: %w",
				fiscalDateEnding, epsParseError)
		}

		records = append(records, types.AnnualEarningRecord{
			Ticker:           ticker,
			FiscalDateEnding: fiscalDateEnding,
			ReportedEPS:      reportedEPS,
		})
	}

	return records, nil
}

type StockSplitResponse struct {
	Symbol string       `json:"symbol"`
	Data   []SplitEvent `json:"data"`
}

type SplitEvent struct {
	EffectiveDate string `json:"effective_date"`
	SplitFactor   string `json:"split_factor"`
}

// Takes json data of stock splits and parses it into a collection
// of individual stock split events.
func ParseStockSplitsToFlat(jsonData []byte) ([]types.StockSplitRecord, error) {
	var response StockSplitResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling stock splits json: %w", err)
	}

	ticker := response.Symbol
	if ticker == "" {
		return nil, fmt.Errorf("ticker not found in JSON when parsing stock splits")
	}

	records := make([]types.StockSplitRecord, 0, len(response.Data))
	for _, split := range response.Data {
		splitFactor, err := strconv.ParseFloat(split.SplitFactor, 64)
		if err != nil {
			// Skip records with unparseable split factors
			log.Printf("Warning: could not parse split factor for date %s, skipping record. Error: %v", split.EffectiveDate, err)
			continue
		}

		records = append(records, types.StockSplitRecord{
			Ticker:        ticker,
			EffectiveDate: split.EffectiveDate,
			SplitFactor:   splitFactor,
		})
	}

	return records, nil
}
