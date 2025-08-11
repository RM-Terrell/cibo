package types

// FlatStockRecord is a clean, tidy format for final stock data.
// It represents a single day's closing price for a given stock.
type FlatStockRecord struct {
	Ticker       string  `parquet:"name=ticker,type=BYTE_ARRAY,convertedtype=UTF8"`
	Date         string  `parquet:"name=date,type=BYTE_ARRAY,convertedtype=UTF8"`
	ClosingPrice float64 `parquet:"name=close,type=DOUBLE"`
}
