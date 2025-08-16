package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchDailyStockData retrieves the raw daily time series data for a given stock symbol.
// It returns the raw response body as a byte slice
func (c *Client) FetchDailyStockData(symbol string) ([]byte, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&outputsize=full",
		symbol,
		c.apiKey,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}
