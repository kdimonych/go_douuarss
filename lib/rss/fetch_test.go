package rss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func TestFetchOk(t *testing.T) {
	ctx := context.Background()

	expectedData, err := os.ReadFile("testdata/test_real_rss_blob.xml")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(string(expectedData))),
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()
	data, err := fetchClient.Fetch(ctx, testUrl)
	if err != nil {
		t.Fatalf("Failed to fetch and parse RSS feed: %v", err)
	}

	if !bytes.Equal(expectedData, data) {
		t.Fatalf("Fetched data does not match expected data. Expected: %s, Got: %s", string(expectedData), string(data))
	}
}

func TestFetchErrorCodeUnreachable(t *testing.T) {
	ctx := context.Background()

	client := &http.Client{
		Transport: &mockRoundTripper{resp: nil, err: fmt.Errorf("network error")},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()
	_, err := fetchClient.Fetch(ctx, testUrl)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeUnreachable {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeUnreachable, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}

func TestFetch_ErrorCodeHttpError(t *testing.T) {
	ctx := context.Background()

	mockResp := &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("404 Not Found")),
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()
	_, err := fetchClient.Fetch(ctx, testUrl)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeHttpError {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeHttpError, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}

type errorReadCloser struct {
	err error
}

func (e *errorReadCloser) Read([]byte) (int, error) {
	return 0, e.err
}

func (*errorReadCloser) Close() error {
	return nil
}

func TestFetch_ErrorCodeNoData(t *testing.T) {
	ctx := context.Background()

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       &errorReadCloser{err: fmt.Errorf("forced read error")},
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()
	_, err := fetchClient.Fetch(ctx, testUrl)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeNoData {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeNoData, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}

func TestFetch_ErrorCodeInvalidUrl(t *testing.T) {
	ctx := context.Background()
	testUrl := "bla bla"

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()
	_, err := fetchClient.Fetch(ctx, testUrl)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeInvalidUrl {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeInvalidUrl, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}
