package utils

import (
	"cibo/internal/types"
)

// AdjustForStockSplits adjusts historical stock prices for stock splits. Assumes data is pre sorted by date
func AdjustForStockSplits(dailyPrices []types.DailyStockRecord, splits []types.StockSplitRecord) []types.DailyStockRecord {
	adjustedPrices := make([]types.DailyStockRecord, len(dailyPrices))
	copy(adjustedPrices, dailyPrices)

	if len(splits) == 0 {
		return adjustedPrices
	}

	cumulativeFactor := 1.0
	splitIndex := 0

	for i := range adjustedPrices {
		for splitIndex < len(splits) && adjustedPrices[i].Date < splits[splitIndex].EffectiveDate {
			cumulativeFactor *= splits[splitIndex].SplitFactor
			splitIndex++
		}

		if cumulativeFactor != 1.0 {
			adjustedPrices[i].ClosingPrice /= cumulativeFactor
		}
	}

	return adjustedPrices
}
