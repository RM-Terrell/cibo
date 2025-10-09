package types

// ! Package note: converting to structs meant to write to files contains
// ! the assumption of the data structures being identical, which is
// ! enshrined in the use of type conversion syntax that WILL explode if the
// ! two structs ever differ in any way. This is intentional to help enshrine
// ! data integrity between structs.

// Converts a slice of annual earnings to a form ready for a Parquet file
func AnnualEarningsToParquet(
	records []AnnualEarningRecord) []AnnualEarningRecordParquet {
	parquetRecords := make([]AnnualEarningRecordParquet, len(records))

	for index, record := range records {
		parquetRecords[index] = AnnualEarningRecordParquet(record)
	}

	return parquetRecords
}

// Converts a slice of combined price records for Parquet writing.
func CombinedPricesToParquet(
	records []CombinedPriceRecord) []CombinedPriceRecordParquet {
	parquetRecords := make([]CombinedPriceRecordParquet, len(records))
	for i, record := range records {
		parquetRecords[i] = CombinedPriceRecordParquet(record)
	}
	return parquetRecords
}

func DailyAndFairPriceToCombined(
	dailyPrices []DailyStockRecord,
	fairValuePrices []FairValuePriceRecord) []CombinedPriceRecord {

	combinedData := make([]CombinedPriceRecord, 0, len(dailyPrices)+len(fairValuePrices))

	for _, record := range dailyPrices {
		combinedData = append(combinedData, CombinedPriceRecord{
			Ticker: record.Ticker,
			Date:   record.Date,
			Price:  record.ClosingPrice,
			Series: "daily_price",
		})
	}

	for _, record := range fairValuePrices {
		combinedData = append(combinedData, CombinedPriceRecord{
			Ticker: record.Ticker,
			Date:   record.Date,
			Price:  record.FairValuePrice,
			Series: "fair_value",
		})
	}

	return combinedData
}
