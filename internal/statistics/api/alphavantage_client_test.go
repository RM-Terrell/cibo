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

// Given a valid symbol for daily prices, verify that the response body is returned correctly.
func TestFetchDailyPrice_Success(t *testing.T) {
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

	body, err := apiClient.FetchDailyPrice("IBM")

	if err != nil {
		t.Errorf("Did not expect an error but got: %v", err)
	}

	if diff := cmp.Diff(string(expectedBody), string(body)); diff != "" {
		t.Errorf("FetchDailyPrice() mismatch (-want +got):\n%s", diff)
	}
}

// Given a network error occurs in daily prices, verify that a request failure error is returned.
func TestFetchDailyPrice_NetworkError(t *testing.T) {
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

	body, err := apiClient.FetchDailyPrice("GOOG")

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

// Given the API returns a non-200 status for daily prices, verify that a status code error is returned.
func TestFetchDailyPrice_APIReturnsNon200(t *testing.T) {
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

	body, err := apiClient.FetchDailyPrice("AAPL")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyPrice() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

// Given an error occurs while reading the response body for daily prices, verify that a read failure error is returned.
func TestFetchDailyPrice_BodyReadError(t *testing.T) {
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

	body, err := apiClient.FetchDailyPrice("TSLA")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyPrice() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

//----------------------------------

// Given a valid symbol for earnings, verify that the response body is returned correctly.
func TestFetchEarnings_Success(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"annualEarnings":[]}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedBody := []byte(`{"annualEarnings":[]}`)

	body, err := apiClient.FetchEarnings("IBM")

	if err != nil {
		t.Errorf("Did not expect an error but got: %v", err)
	}
	if diff := cmp.Diff(string(expectedBody), string(body)); diff != "" {
		t.Errorf("FetchEarnings() mismatch (-want +got):\n%s", diff)
	}
}

// Given a network error occurs in earnings data, verify that a request failure error is returned.
func TestFetchEarnings_NetworkError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{Err: errSimulatedNetwork},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}

	_, err := apiClient.FetchEarnings("IBM")

	if !errors.Is(err, errSimulatedNetwork) {
		t.Errorf("Expected error to be '%v', but got '%v'", errSimulatedNetwork, err)
	}
}

// Given the API returns a non-200 status for daily prices, verify that a status code error is returned.
func TestFetchEarnings_APIReturnsNon200(t *testing.T) {
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

	body, err := apiClient.FetchEarnings("AAPL")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyPrice() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

// Given an error occurs while reading the response body for earnings, verify that a read failure error is returned.
func TestFetchEarnings_BodyReadError(t *testing.T) {
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

	body, err := apiClient.FetchEarnings("TSLA")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyPrice() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

//----------------------------------

// Given a valid symbol for overview data, verify that the response body is returned correctly.
func TestFetchOverview_Success(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"Symbol":"IBM"}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedBody := []byte(`{"Symbol":"IBM"}`)

	body, err := apiClient.FetchOverview("IBM")

	if err != nil {
		t.Errorf("Did not expect an error but got: %v", err)
	}
	if diff := cmp.Diff(string(expectedBody), string(body)); diff != "" {
		t.Errorf("FetchOverview() mismatch (-want +got):\n%s", diff)
	}
}

// Given the API returns a non-200 status for overview, verify that a status code error is returned
func TestFetchOverview_APIReturnsNon200(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":"bad request"}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedError := "API returned non-200 status code: 400, body: {\"error\":\"bad request\"}"

	_, err := apiClient.FetchOverview("IBM")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchOverview() error mismatch (-want +got):\n%s", diff)
	}
}

// Given a network error occurs in overview data, verify that a request failure error is returned.
func TestFetchOverview_NetworkError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{Err: errSimulatedNetwork},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}

	_, err := apiClient.FetchOverview("IBM")

	if !errors.Is(err, errSimulatedNetwork) {
		t.Errorf("Expected error to be '%v', but got '%v'", errSimulatedNetwork, err)
	}
}

// Given an error occurs while reading the response body for overview data, verify that a read failure error is returned.
func TestFetchOverview_BodyReadError(t *testing.T) {
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

	body, err := apiClient.FetchOverview("TSLA")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDailyPrice() error mismatch (-want +got):\n%s", diff)
	}
	if body != nil {
		t.Errorf("Expected body to be nil but got: %s", string(body))
	}
}

//----------------------------------

// Given a valid symbol for dividend data, verify that the response body is returned correctly.
func TestFetchDividends_Success(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedBody := []byte(`{"data":[]}`)

	body, err := apiClient.FetchDividends("IBM")

	if err != nil {
		t.Errorf("Did not expect an error but got: %v", err)
	}
	if diff := cmp.Diff(string(expectedBody), string(body)); diff != "" {
		t.Errorf("FetchDividends() mismatch (-want +got):\n%s", diff)
	}
}

// Given an error occurs while reading the response body for dividend data, verify that a read failure error is returned.
func TestFetchDividends_BodyReadError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(&errorReader{}),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedError := "failed to read response body: simulated read error"

	_, err := apiClient.FetchDividends("IBM")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDividends() error mismatch (-want +got):\n%s", diff)
	}
}

// Given the API returns a non-200 status for dividend data, verify that a status code error is returned.
func TestFetchDividend_APIReturnsNon200(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":"bad request"}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedError := "API returned non-200 status code: 400, body: {\"error\":\"bad request\"}"

	_, err := apiClient.FetchDividends("IBM")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchOverview() error mismatch (-want +got):\n%s", diff)
	}
}

// Given a network error occurs in dividend data, verify that a request failure error is returned.
func TestFetchDividend_NetworkError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{Err: errSimulatedNetwork},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}

	_, err := apiClient.FetchDividends("IBM")

	if !errors.Is(err, errSimulatedNetwork) {
		t.Errorf("Expected error to be '%v', but got '%v'", errSimulatedNetwork, err)
	}
}

// ----------------------------------
// Given a valid symbol for earnings estimates, verify that the response body is returned correctly.
func TestFetchEarningsEstimates_Success(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"estimates":[]}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedBody := []byte(`{"estimates":[]}`)

	body, err := apiClient.FetchEarningsEstimates("IBM")

	if err != nil {
		t.Errorf("Did not expect an error but got: %v", err)
	}
	if diff := cmp.Diff(string(expectedBody), string(body)); diff != "" {
		t.Errorf("FetchEarningsEstimates() mismatch (-want +got):\n%s", diff)
	}
}

// Given a network error occurs in earnings estimate data, verify that a request failure error is returned.
func TestFetchEarningsEstimates_NetworkError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{Err: errSimulatedNetwork},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}

	_, err := apiClient.FetchEarningsEstimates("IBM")

	if !errors.Is(err, errSimulatedNetwork) {
		t.Errorf("Expected error to be '%v', but got '%v'", errSimulatedNetwork, err)
	}
}

// Given the API returns a non-200 status for earnings estimate data, verify that a status code error is returned.
func TestFetchEarningsEstimate_APIReturnsNon200(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":"bad request"}`)),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedError := "API returned non-200 status code: 400, body: {\"error\":\"bad request\"}"

	_, err := apiClient.FetchEarningsEstimates("IBM")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchOverview() error mismatch (-want +got):\n%s", diff)
	}
}

// Given an error occurs while reading the response body for earnings estimate data, verify that a read failure error is returned.
func TestFetchEarningsEstimate_BodyReadError(t *testing.T) {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(&errorReader{}),
			},
		},
	}
	apiClient := &Client{apiKey: "test_api_key", httpClient: mockClient}
	expectedError := "failed to read response body: simulated read error"

	_, err := apiClient.FetchEarningsEstimates("IBM")

	if err == nil {
		t.Fatal("Expected an error, but got none.")
	}
	if diff := cmp.Diff(expectedError, err.Error()); diff != "" {
		t.Errorf("FetchDividends() error mismatch (-want +got):\n%s", diff)
	}
}

//----------------------------------

// Given a new client, verify that a new client is created with the correct default values.
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
