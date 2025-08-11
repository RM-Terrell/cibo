package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"stats_engine/types"
)

// example expected json
/*
   {
   "Meta Data": {
       "1. Information": "Daily Prices (open, high, low, close) and Volumes",
       "2. Symbol": "IBM",
       "3. Last Refreshed": "2025-08-01",
       "4. Output Size": "Full size",
       "5. Time Zone": "US/Eastern"
   },
   "Time Series (Daily)": {
       "2025-08-01": {
           "1. open": "251.4050",
           "2. high": "251.4791",
           "3. low": "245.6100",
           "4. close": "250.0500",
           "5. volume": "9683404"
       },
       ...
*/

type AlphaVantageResponse struct {
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
Function to take jsonData from an API response and parse it into a collection
of individual stock prices.
*/
func ParseToFlat(jsonData []byte, skipErrors bool) ([]types.FlatStockRecord, error) {
	var response AlphaVantageResponse

	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %w", err)
	}

	ticker := response.MetaData.Symbol
	if ticker == "" {
		return nil, fmt.Errorf("ticker not found in JSON Meta Data when parsing")
	}

	records := make([]types.FlatStockRecord, 0, len(response.TimeSeries))

	for rawDate, rawDataPoint := range response.TimeSeries {
		closingPrice, err := strconv.ParseFloat(rawDataPoint.Close, 64)
		if err != nil {
			if skipErrors {
				log.Printf("Warning: could not parse close price for date %s, skipping record. Error: %v", rawDate, err)
				continue
			}
			return nil, fmt.Errorf("could not parse close price '%s' for date %s: %w", rawDataPoint.Close, rawDate, err)
		}

		records = append(records, types.FlatStockRecord{
			Ticker:       ticker,
			Date:         rawDate,
			ClosingPrice: closingPrice,
		})
	}

	return records, nil
}
