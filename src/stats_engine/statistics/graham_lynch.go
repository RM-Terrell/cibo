package statistics

/*
The goal of this module to calculate a "fair market value" of a given stock
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

1. Annual earnings data
2. Calculate long term growth rate
	Compound Annual Growth Rate (CAGR) =(Ending EPS/Beginning EPS)^(1/Number of Years) − 1

3. Calculate fair value PE ratio based on CAGR
	Fair Value P/E = Earnings Growth Rate (CAGR) × 100
4. Calculate fair value prices for time period to create a curve
	FairValuePrice(t) = reportedEPS(t) * FairValueP/E
*/
