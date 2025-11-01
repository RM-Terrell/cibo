package main

import (
	"log"
	"net/http"
	"os"
)

var (
	aaplDailyPriceJSON []byte
	aaplEarningsJSON   []byte
	aaplSplitsJSON     []byte
)

// A helper function to read a file and exit fatally if it fails, as the mock server that can't run without its data.
func loadJSONData(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read mock data file %s: %v", filePath, err)
	}
	return data
}

// Mimics the Alpha Vantage /query endpoint.
func queryHandler(w http.ResponseWriter, r *http.Request) {
	function := r.URL.Query().Get("function")
	symbol := r.URL.Query().Get("symbol")

	log.Printf("Received request for function: %s, symbol: %s", function, symbol)

	w.Header().Set("Content-Type", "application/json")

	if symbol != "AAPL" {
		log.Printf("Unsupported symbol requested: %s", symbol)
		http.Error(w, `{"error": "Unsupported mock symbol. Only AAPL is available."}`, http.StatusBadRequest)
		return
	}

	var responseData []byte

	switch function {
	case "TIME_SERIES_DAILY":
		responseData = aaplDailyPriceJSON
	case "EARNINGS":
		responseData = aaplEarningsJSON
	case "SPLITS":
		responseData = aaplSplitsJSON
	default:
		log.Printf("Unknown function requested: %s", function)
		http.Error(w, `{"error": "Unknown function"}`, http.StatusBadRequest)
		return
	}

	_, err := w.Write(responseData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	log.Println("Loading mock data from files...")
	aaplDailyPriceJSON = loadJSONData("tools/mock_alpha_vantage_api/data/AAPL_TIME_SERIES_DAILY.json")
	aaplEarningsJSON = loadJSONData("tools/mock_alpha_vantage_api/data/AAPL_EARNINGS.json")
	aaplSplitsJSON = loadJSONData("tools/mock_alpha_vantage_api/data/AAPL_SPLITS.json")
	log.Println("Mock data loaded successfully.")

	http.HandleFunc("/query", queryHandler)

	port := "8080"
	log.Printf("Starting mock API server on http://localhost:%s", port)
	log.Println("Press Ctrl+C to shut down.")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
