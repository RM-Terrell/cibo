package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type MockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

var errSimulatedNetwork = errors.New("simulated network failure")

// Given a valid symbol, verify that the response body is returned correctly.
func TestFetchDailyStockData_Success(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"key":"value"}`)),
			},
			Err: nil,
		},
	}
	apiClient := &Client{
		apiKey:     "test_api_key",
		httpClient: mockClient,
	}
	expectedBody := []byte(`{"key":"value"}`)

	body, err := apiClient.FetchDailyStockData("IBM")

	if err != nil {
		t.Errorf("Did not expect an error but got: %v", err)
	}

	if diff := cmp.Diff(string(expectedBody), string(body)); diff != "" {
		t.Errorf("FetchDailyStockData() mismatch (-want +got):\n%s", diff)
	}
}

// Given a network error occurs, verify that a request failure error is returned.
func TestFetchDailyStockData_NetworkError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: nil,
			Err:      errSimulatedNetwork,
		},
	}
	apiClient := &Client{
		apiKey:     "test_api_key",
		httpClient: mockClient,
	}

	body, err := apiClient.FetchDailyStockData("GOOG")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if !errors.Is(err, errSimulatedNetwork) {
		t.Errorf("Expected error to be '%v', but got '%v'", errSimulatedNetwork, err)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

// Given the API returns a non-200 status, verify that a status code error is returned.
func TestFetchDailyStockData_APIReturnsNon200(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":"symbol not found"}`)),
			},
			Err: nil,
		},
	}
	apiClient := &Client{
		apiKey:     "test_api_key",
		httpClient: mockClient,
	}
	expectedError := "API returned non-200 status code: 404, body: {\"error\":\"symbol not found\"}"

	body, err := apiClient.FetchDailyStockData("AAPL")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyStockData() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

// Given an error occurs while reading the response body, verify that a read failure error is returned.
func TestFetchDailyStockData_BodyReadError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(&errorReader{}),
			},
			Err: nil,
		},
	}
	apiClient := &Client{
		apiKey:     "test_api_key",
		httpClient: mockClient,
	}
	expectedError := "failed to read response body: simulated read error"

	body, err := apiClient.FetchDailyStockData("TSLA")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyStockData() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

// Given an a new client, verify that a new client is created with the correct default values.
func TestNewClient(t *testing.T) {
	apiKey := "my-secret-key"
	client := NewClient(apiKey)

	if diff := cmp.Diff(apiKey, client.apiKey); diff != "" {
		t.Errorf("NewClient() apiKey mismatch (-want +got):\n%s", diff)
	}

	if client.httpClient == nil {
		t.Fatal("Expected httpClient to be initialized, but it was nil")
	}

	expectedTimeout := 10 * time.Second
	if diff := cmp.Diff(expectedTimeout, client.httpClient.Timeout); diff != "" {
		t.Errorf("NewClient() httpClient.Timeout mismatch (-want +got):\n%s", diff)
	}

	numFields := reflect.TypeOf(*client).NumField()
	if diff := cmp.Diff(2, numFields); diff != "" {
		t.Errorf("NewClient() struct field count mismatch (-want +got):\n%s", diff)
	}
}
