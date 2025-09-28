package types

//! Package convention. use snake case for parquet column `name` fields.
//! SQL compat, readability, because i said so, etc

type DailyStockRecord struct {
	Ticker       string
	Date         string
	ClosingPrice float64
}

type AnnualEarningRecord struct {
	Ticker           string
	FiscalDateEnding string
	ReportedEPS      float64
}

type FairValuePriceRecord struct {
	Ticker         string
	FairValuePrice float64
	Date           string
}

/*
Intention of this type is to allow "long" writing of price data. Example:

ticker	date	series	        price
XYZ	2025-12-29	actual_price	150.25
XYZ	2025-12-30	actual_price	151.00
XYZ	2025-12-31	actual_price	150.80
XYZ	2025-12-31	fair_value	    175.50
*/
type CombinedPriceRecord struct {
	Ticker string
	Date   string
	Price  float64
	Series string // fair value estimate, actually daily, etc
}

// ---- Parquet types
//! New parquet types must be added to type_test.go for convention testing

type CombinedPriceRecordParquet struct {
	Ticker string  `parquet:"name=ticker,type=BYTE_ARRAY,convertedtype=UTF8"`
	Date   string  `parquet:"name=date,type=BYTE_ARRAY,convertedtype=UTF8"`
	Price  float64 `parquet:"name=price,type=DOUBLE"`
	Series string  `parquet:"name=series,type=BYTE_ARRAY,convertedtype=UTF8"`
}

type AnnualEarningRecordParquet struct {
	Ticker           string  `parquet:"name=ticker,type=BYTE_ARRAY,convertedtype=UTF8"`
	FiscalDateEnding string  `parquet:"name=fiscal_date_ending,type=BYTE_ARRAY,convertedtype=UTF8"`
	ReportedEPS      float64 `parquet:"name=reported_eps,type=DOUBLE"`
}
