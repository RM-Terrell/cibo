package algos

import (
	"fmt"
	"math"
	"sort"
	"time"

	"cibo/internal/types"
)

/*
The goal of this module is to calculate a "fair market value" of a given stock
using traditional fair value equations from Benjamin Graham and Peter Lynch,
to act as a model for the valuation of a company based on past and projected
future earnings.

Fair Value Price = Earnings Per Share (EPS) × Fair Value P/E Ratio

The number of years chosen for calculating the growth rate will wildly effect
the Compound Annual Growth Rate, it is likely this should be setup as a user
input in the system to run different simulations. Current (the last few years)
growth compared to multi decade history might or might not be desired depending
on context.

https://www.investopedia.com/terms/p/pegyratio.asp

https://www.investopedia.com/terms/b/benjamin-method.asp

General calculation flow process:
	1. Get annual earnings data collection
	2. Calculate long term growth rate (single number)
		Compound Annual Growth Rate (CAGR) =(Ending EPS/Beginning EPS)^(1/Number of Years) − 1
	3. Calculate long term growth PE ratio based on CAGR (single number)
		Fair Value P/E = (CAGR) × 100
	4. Calculate fair value prices for time period to create a curve (collection of data)
		FairValuePrice(t) = reportedEPS(t) * FairValueP/E

The point of all of this is try to find situations where a stocks current price may have either
outran the companies underlying business dynamics and thus may yield poor returns in the future,
or have possibly over sold to a value under what the business is really "worth" and thus
represent a good buying opportunity.

Remember. All models are wrong but some are useful.
*/

/*
Function to find the starting and ending positive EPS values from a collection of annual earnings data.
CAGR only works on stable, profitable growth companies, so this function removes early unprofitable years
of data and then boldly assumes all other years will be positive. Handling of intermittent negative years
needs to be handled by another function.
*/
func ProfitableEarningsStartingAndEnding(earnings []types.AnnualEarningRecord) (start types.AnnualEarningRecord, end types.AnnualEarningRecord, err error) {
	minimumYearsRequired := 2
	if len(earnings) < minimumYearsRequired {
		err = fmt.Errorf("not enough data points. Minimum: %d", minimumYearsRequired)
		return
	}

	// Creating a copy here to avoid unintended slice sort side effects
	sortedEarnings := make([]types.AnnualEarningRecord, len(earnings))
	copy(sortedEarnings, earnings)

	/*
		Data will likely already be sorted but sort it anyways to be defensive. When the day comes that annual
		earnings data sets are so large to be a performance issue when sorting, this program wont matter anyways.
	*/
	sort.Slice(sortedEarnings, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", sortedEarnings[i].FiscalDateEnding)
		dateJ, _ := time.Parse("2006-01-02", sortedEarnings[j].FiscalDateEnding)
		return dateI.Before(dateJ)
	})

	/*
		Find the first earning with a positive EPS to use as the starting point in cases where a company
		had negative earnings in its early years.
	*/
	foundStart := false
	for _, earning := range sortedEarnings {
		if earning.ReportedEPS > 0 {
			start = earning
			foundStart = true
			break
		}
	}

	/*
		TODO this might be a good target for logic to signal cases with lots of negative earnings with
		better user feedback and suggestion for analysis of early / unprofitable companies
	*/
	if !foundStart {
		err = fmt.Errorf("no valid starting point with positive EPS found. Company just burns money")
		return
	}

	end = sortedEarnings[len(sortedEarnings)-1]

	if end.ReportedEPS <= 0 {
		err = fmt.Errorf("ending EPS must be positive for a current CAGR calculation")
		return
	}

	return start, end, nil
}

/*
Calculate the Compound Annual Growth Rate for a collection of annual earnings.
*/
func CAGR(earnings []types.AnnualEarningRecord) (float64, error) {
	startEarning, endEarning, err := ProfitableEarningsStartingAndEnding(earnings)
	if err != nil {
		return 0, fmt.Errorf("could not determine calculation endpoints: %w", err)
	}

	startDate, _ := time.Parse("2006-01-02", startEarning.FiscalDateEnding)
	endDate, _ := time.Parse("2006-01-02", endEarning.FiscalDateEnding)

	/*
		Using 365.25 accounts for leap years.
		Possible cause of calculation discrepancies with other systems
	*/
	years := endDate.Sub(startDate).Hours() / 24 / 365.25
	if years < 1.0 {
		return 0, fmt.Errorf("period between start and end date must be at least one year")
	}

	growthRatio := endEarning.ReportedEPS / startEarning.ReportedEPS
	cagr := math.Pow(growthRatio, 1.0/years) - 1.0

	return cagr, nil
}

/*
Calculate the Fair Value PE ratio given a Compound Annual Growth Rate.
*/
func FairValuePE(cagr float64) float64 {
	return cagr * 100
}

/*
Generate a history of estimated fair value prices for a given stock based on its earnings history
and a fair value PE ratio.
WARNING: This function assumes sorted data.
*/
func FairValuePriceHistory(fairValuePE float64, historicalEarnings []types.AnnualEarningRecord) []types.FairValuePriceRecord {
	var fairValuePriceHistory []types.FairValuePriceRecord
	for _, earning := range historicalEarnings {
		// You cant really calculate EPS on a negative earnings or you get a negative fair value.
		// Again, lots of negative values might need to be a case to bubble up to the user
		// and suggest using a different analysis. For now, skip and leave a hole in the plot
		if earning.ReportedEPS > 0 {
			fairValuePrice := earning.ReportedEPS * fairValuePE
			fairValuePriceHistory = append(fairValuePriceHistory, types.FairValuePriceRecord{
				Ticker:         earning.Ticker,
				FairValuePrice: fairValuePrice,
				Date:           earning.FiscalDateEnding,
			})
		}
	}

	return fairValuePriceHistory
}

/*
Pipeline function that orchestrates the full fair value calculation process.
*/
func CalculateFairValueHistory(earnings []types.AnnualEarningRecord) ([]types.FairValuePriceRecord, error) {
	cagr, err := CAGR(earnings)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate CAGR: %w", err)
	}

	fairValuePE := FairValuePE(cagr)

	fairValueHistory := FairValuePriceHistory(fairValuePE, earnings)

	return fairValueHistory, nil
}
