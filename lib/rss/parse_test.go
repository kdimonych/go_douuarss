package rss

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func checkItem(t *testing.T, items []Item) {
	for _, item := range items {
		if item.Title != "Test Item" {
			t.Errorf("Expected item title 'Test Item', got '%s'", item.Title)
		}
		if item.Link != "http://example.com/item" {
			t.Errorf("Expected item link 'http://example.com/item', got '%s'", item.Link)
		}
		if item.Description != "Test Item Description" {
			t.Errorf("Expected item description 'Test Item Description', got '%s'", item.Description)
		}
		var expectedPubDate, err = time.Parse(time.RFC1123Z, "Mon, 09 Jun 2025 11:40:06 +0300")
		if err != nil {
			t.Fatalf("Failed to parse expected expectedPubDate: %v", err)
		}
		if item.PubDate.Time != expectedPubDate {
			t.Errorf("Expected item pubDate '%s', got '%s'", expectedPubDate.String(), item.PubDate.String())
		}
		if item.Creator != "Test Creator" {
			t.Errorf("Expected item creator 'Test Creator', got '%s'", item.Creator)
		}
	}
}

func checkChannel(t *testing.T, channels []Channel) {
	for _, channel := range channels {
		if channel.Title != "Test Channel" {
			t.Errorf("Expected channel title 'Test Channel', got '%s'", channel.Title)
		}
		if channel.Link != "http://example.com" {
			t.Errorf("Expected channel title 'http://example.com', got '%s'", channel.Link)
		}
		if channel.Description != "Test Description" {
			t.Errorf("Expected channel description 'Test Description', got '%s'", channel.Description)
		}
		if channel.Language != "en-us" {
			t.Errorf("Expected channel language 'en-us', got '%s'", channel.Language)
		}

		var expectedLastBuildDate, err = time.Parse(time.RFC1123Z, "Mon, 09 Jun 2025 10:30:06 +0300")
		if err != nil {
			t.Fatalf("Failed to parse expected lastBuildDate: %v", err)
		}
		if channel.LastBuildDate.Time != expectedLastBuildDate {
			t.Errorf("Expected channel lastBuildDate '%s', got '%s'", expectedLastBuildDate.String(), channel.LastBuildDate.String())
		}
		if len(channel.Items) == 0 {
			t.Fatal("Expected at least one item in the channel")
		}

		checkItem(t, channel.Items)
	}
}

func TestParse(t *testing.T) {
	blob := []byte(`
	<rss>
		<channel>
			<title>Test Channel</title>
			<link>http://example.com</link>
			<description>Test Description</description>
			<language>en-us</language>
			<lastBuildDate>Mon, 09 Jun 2025 10:30:06 +0300</lastBuildDate>
			<item>
				<title xmlns:dc="http://purl.org/dc/elements/1.1/">Test Item</title>
				<link>http://example.com/item</link>
				<description>Test Item Description</description>
				<pubDate>Mon, 09 Jun 2025 11:40:06 +0300</pubDate>
				<dc:creator>Test Creator</dc:creator>
			</item>
		</channel>
	</rss>`)

	channels, err := Parse(blob)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(channels) == 0 {
		t.Fatal("Expected at least one channel")
	}

	checkChannel(t, channels)
}

func TestParseRealFeed(t *testing.T) {
	badBlob := []byte(`Bad XML Data`)

	_, err := Parse(badBlob)
	if err == nil {
		t.Fatalf("Expected parse error, got nil")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeInvalidData {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeInvalidData, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}

func TestParse_BadData(t *testing.T) {
	data, err := os.ReadFile("testdata/test_real_rss_blob.xml")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	channels, parseErr := Parse(data)
	if parseErr != nil {
		t.Fatalf("Parse error: %v", parseErr)
	}

	if len(channels) == 0 {
		t.Fatal("Expected at least one channel")
	}
}

func TestFetchAndParse(t *testing.T) {
	ctx := context.Background()

	data, err := os.ReadFile("testdata/test_real_rss_blob.xml")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(string(data))),
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()

	channels, err := FetchAndParse(fetchClient, ctx, testUrl)
	if err != nil {
		t.Fatalf("Failed to fetch and parse RSS feed: %v", err)
	}

	if len(channels) == 0 {
		t.Fatal("No channels found in the RSS feed")
	}

	log.Printf("Fetched %d channels from %s", len(channels), testUrl)
	for _, channel := range channels {
		log.Printf("Channel: %s", channel.Title)
	}
}

func TestFetchAndParse_NoFetchClient(t *testing.T) {
	ctx := context.Background()

	_, err := FetchAndParse(nil, ctx, testUrl)
	if err == nil {
		t.Fatalf("Expected error when fetch client is nil, got nil")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeInternalError {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeInternalError, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}

func TestFetchAndParse_FetchError(t *testing.T) {
	ctx := context.Background()

	mockResp := &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("404 Not Found")),
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()

	_, err := FetchAndParse(fetchClient, ctx, testUrl)
	if err == nil {
		t.Fatalf("Expected error when fetch client is nil, got nil")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeHttpError {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeHttpError, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}

func TestFetchAndParse_ParseError(t *testing.T) {
	ctx := context.Background()

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("Bad XML Data")),
	}

	client := &http.Client{
		Transport: &mockRoundTripper{resp: mockResp, err: nil},
	}

	fetchClient := NewFetchClientBuilder().WithHttpClient(client).Build()

	_, err := FetchAndParse(fetchClient, ctx, testUrl)
	if err == nil {
		t.Fatalf("Expected error when fetch client is nil, got nil")
	}

	if fetchErr, ok := err.(*FetchError); ok {
		if fetchErr.Code != ErrorCodeInvalidData {
			t.Fatalf("Expected error code %d, got %d", ErrorCodeInvalidData, fetchErr.Code)
		}
	} else {
		t.Fatalf("Expected FetchError, got %T: %v", err, err)
	}
}
