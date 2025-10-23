package pipelines

import (
	"cibo/internal/types"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// --- Mock Implementations ---
type mockAPIClient struct {
	dailyPriceResponse   []byte
	earningsResponse     []byte
	shouldReturnFetchErr bool
}

func (m *mockAPIClient) FetchDailyPrice(ticker string) ([]byte, error) {
	if m.shouldReturnFetchErr {
		return nil, errors.New("mock API fetch error")
	}
	return m.dailyPriceResponse, nil
}

func (m *mockAPIClient) FetchEarnings(ticker string) ([]byte, error) {
	if m.shouldReturnFetchErr {
		return nil, errors.New("mock API fetch error")
	}
	return m.earningsResponse, nil
}

type mockParquetWriter struct {
	shouldReturnWriteErr bool
	wasCalled            bool
	receivedData         []types.CombinedPriceRecord
}

func (m *mockParquetWriter) WriteCombinedPriceDataToParquet(records []types.CombinedPriceRecord, writer io.WriteCloser) (string, error) {
	m.wasCalled = true
	m.receivedData = records
	if m.shouldReturnWriteErr {
		return "", errors.New("mock parquet write error")
	}
	// Return a mock log message on success
	return "mock write success log", nil
}

// --- Test Suite ---

// Given that all dependencies succeed, verify that the pipeline runs correctly
// and produces the expected combined data output.
func TestLynchFairValuePipeline_RunPipeline_Success(t *testing.T) {
	mockClient := &mockAPIClient{
		dailyPriceResponse: []byte(`{
			"Meta Data": {"2. Symbol": "TEST"},
			"Time Series (Daily)": {"2025-01-01": {"4. close": "150.00"}}
		}`),
		earningsResponse: []byte(`{
			"symbol": "TEST",
			"annualEarnings": [
				{"fiscalDateEnding": "2024-12-31", "reportedEPS": "10.00"},
				{"fiscalDateEnding": "2023-12-31", "reportedEPS": "5.00"}
			]
		}`),
		shouldReturnFetchErr: false,
	}
	mockWriter := &mockParquetWriter{}

	dummyFileName := "TEST.parquet"
	defer os.Remove(dummyFileName)

	pipeline := NewLynchFairValuePipeline(mockClient, mockWriter)
	input := LynchFairValueInputs{Ticker: "TEST"}

	output, err := pipeline.RunPipeline(input)

	if err != nil {
		t.Fatalf("RunPipeline() returned an unexpected error: %v", err)
	}
	if output == nil {
		t.Fatal("RunPipeline() returned a nil output on success")
	}

	expectedData := []types.CombinedPriceRecord{
		{Ticker: "TEST", Date: "2025-01-01", Price: 150.00, Series: "daily_price"},
		{Ticker: "TEST", Date: "2024-12-31", Price: 997.1612494011704, Series: "fair_value"},
		{Ticker: "TEST", Date: "2023-12-31", Price: 498.5806247005852, Series: "fair_value"},
	}

	// Use a sorter to make the test robust against the order of appends.
	sorter := cmpopts.SortSlices(func(a, b types.CombinedPriceRecord) bool {
		if a.Date != b.Date {
			return a.Date > b.Date // Sort by date descending to match typical time series
		}
		return a.Series < b.Series // Then by series for stability
	})

	if diff := cmp.Diff(expectedData, output.CombinedPriceData, sorter); diff != "" {
		t.Errorf("RunPipeline() mismatch in CombinedPriceData (-want +got):\n%s", diff)
	}
	if !mockWriter.wasCalled {
		t.Error("Expected WriteCombinedPriceDataToParquet to be called, but it was not")
	}
	if output.FileName != dummyFileName {
		t.Errorf("Expected FileName to be '%s', got '%s'", dummyFileName, output.FileName)
	}
}

// Given that the API client returns an error, verify that the pipeline
// stops and propagates the error correctly.
func TestLynchFairValuePipeline_RunPipeline_APIFetchError(t *testing.T) {
	mockClient := &mockAPIClient{shouldReturnFetchErr: true} // This is the failure case
	mockWriter := &mockParquetWriter{}

	pipeline := NewLynchFairValuePipeline(mockClient, mockWriter)
	input := LynchFairValueInputs{Ticker: "FAIL"}

	output, err := pipeline.RunPipeline(input)

	if err == nil {
		t.Fatal("RunPipeline() was expected to return an error, but it returned nil")
	}
	if output != nil {
		t.Errorf("RunPipeline() was expected to return a nil output on error, but it did not")
	}
	if mockWriter.wasCalled {
		t.Error("WriteCombinedPriceDataToParquet should not be called when the API fetch fails")
	}
}

// Given that the Parquet writer returns an error, verify that the pipeline
// stops and propagates the error correctly.
func TestLynchFairValuePipeline_RunPipeline_ParquetWriteError(t *testing.T) {
	// The mock JSON is also updated here to allow the pipeline to proceed to the writing step.
	mockClient := &mockAPIClient{
		dailyPriceResponse: []byte(`{
			"Meta Data": {"2. Symbol": "TEST"},
			"Time Series (Daily)": {"2025-01-01": {"4. close": "150.00"}}
		}`),
		earningsResponse: []byte(`{
			"symbol": "TEST",
			"annualEarnings": [
				{"fiscalDateEnding": "2024-12-31", "reportedEPS": "10.00"},
				{"fiscalDateEnding": "2023-12-31", "reportedEPS": "5.00"}
			]
		}`),
	}
	mockWriter := &mockParquetWriter{shouldReturnWriteErr: true} // This is the failure case

	dummyFileName := "TEST.parquet"
	defer os.Remove(dummyFileName)

	pipeline := NewLynchFairValuePipeline(mockClient, mockWriter)
	input := LynchFairValueInputs{Ticker: "TEST"}

	output, err := pipeline.RunPipeline(input)

	if err == nil {
		t.Fatal("RunPipeline() was expected to return an error for a failed write, but it returned nil")
	}
	if output != nil {
		t.Errorf("RunPipeline() was expected to return a nil output on error, but it did not")
	}
	if !mockWriter.wasCalled {
		t.Error("Expected WriteCombinedPriceDataToParquet to be called, even on failure")
	}
}
