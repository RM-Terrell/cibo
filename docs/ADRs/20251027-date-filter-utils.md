# Utils and date filtering

## Context

I forgot a feature when i first built this thing. The form input takes in a range of dates that will effect not just the range of daily stock prices shown, but also by narrowing the range of annual earnings data used for the Compound Annual Growth Rate the resulting fair value price will change potentially wildly. The stock of Apple is a great example wheres its previous very low PE ratio became irrelevant when the company went from making mostly desktop computer hardware to making iphones which contained an app market place that was basically a money printer. Suddenly a fairly stable company with a defined market had a whole new money printer to run and it was a high tech digital sales and services growth company.

## The problem

But the Alpha Vantage API doesn't take date ranges to filter results so Im doing it myself. One option if I cared about maximum performance would be to implement the date filtering into the various `ParseXToFlat` methods, as they handle date data while iterating each record and could narrow results during that linear scan. That would return the desired data in one go, which then gets sent into the fair value pipeline and eventually charted. I however care more about readability, portability, adaptability, etc than the performance loss of another linear scan of the data and I feel that having all that parsing logic in one method will bite me one day when I want date filtering of my standard Records datatypes anyways as utility functions.

## The solution

So I will be implementing these as separate functions, parsing the data to flat, then scanning it again and narrowing down its date range. If you ever wonder why Im doing multiple traverses of the same data, thats why.
