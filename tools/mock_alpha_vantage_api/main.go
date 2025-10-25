package main

import (
	"fmt"
	"log"
	"net/http"
)

const dailyPriceJSON = `
{
    "Meta Data": {
        "1. Information": "Mock Daily Prices",
        "2. Symbol": "MOCK",
        "3. Last Refreshed": "2025-10-25",
        "4. Output Size": "Full",
        "5. Time Zone": "US/Eastern"
    },
    "Time Series (Daily)": {
        "2025-10-25": {
            "1. open": "150.00",
            "2. high": "152.00",
            "3. low": "149.50",
            "4. close": "151.75",
            "5. volume": "1000000"
        },
        "2025-10-24": {
            "1. open": "148.00",
            "2. high": "150.50",
            "3. low": "147.50",
            "4. close": "149.90",
            "5. volume": "1200000"
        }
    }
}
`

const earningsJSON = `
{
    "symbol": "MOCK",
    "annualEarnings": [
        {
            "fiscalDateEnding": "2024-12-31",
            "reportedEPS": "12.50"
        },
        {
            "fiscalDateEnding": "2023-12-31",
            "reportedEPS": "10.00"
        }
    ],
    "quarterlyEarnings": []
}
`

// Mimics the Alpha Vantage /query endpoint and its various functions and symbols params
func queryHandler(w http.ResponseWriter, r *http.Request) {
	function := r.URL.Query().Get("function")
	symbol := r.URL.Query().Get("symbol")

	log.Printf("Received request for function: %s, symbol: %s", function, symbol)

	w.Header().Set("Content-Type", "application/json")

	switch function {
	case "TIME_SERIES_DAILY":
		fmt.Fprint(w, dailyPriceJSON)
	case "EARNINGS":
		fmt.Fprint(w, earningsJSON)
	default:
		log.Printf("Unknown function requested: %s", function)
		http.Error(w, `{"error": "Unknown function"}`, http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/query", queryHandler)

	port := "8080"
	log.Printf("Starting mock API server on http://localhost:%s", port)
	log.Println("Press Ctrl+C to shut down.")

	// This is a blocking call that starts the server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
