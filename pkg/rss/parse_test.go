package rss

import (
	"os"
	"testing"
	"time"
)

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

		var expectedLastBuildDate time.Time
		expectedLastBuildDate, err = time.Parse(time.RFC1123Z, "Mon, 09 Jun 2025 10:30:06 +0300")
		if err != nil {
			t.Fatalf("Failed to parse expected lastBuildDate: %v", err)
		}
		if channel.LastBuildDate.Time != expectedLastBuildDate {
			t.Errorf("Expected channel lastBuildDate '%s', got '%s'", expectedLastBuildDate.String(), channel.LastBuildDate.Time.String())
		}
		if len(channel.Items) == 0 {
			t.Fatal("Expected at least one item in the channel")
		}

		for _, item := range channel.Items {
			if item.Title != "Test Item" {
				t.Errorf("Expected item title 'Test Item', got '%s'", item.Title)
			}
			if item.Link != "http://example.com/item" {
				t.Errorf("Expected item link 'http://example.com/item', got '%s'", item.Link)
			}
			if item.Description != "Test Item Description" {
				t.Errorf("Expected item description 'Test Item Description', got '%s'", item.Description)
			}
			var expectedPubDate time.Time
			expectedPubDate, err = time.Parse(time.RFC1123Z, "Mon, 09 Jun 2025 11:40:06 +0300")
			if err != nil {
				t.Fatalf("Failed to parse expected expectedPubDate: %v", err)
			}
			if item.PubDate.Time != expectedPubDate {
				t.Errorf("Expected item pubDate '%s', got '%s'", expectedPubDate.String(), item.PubDate.Time.String())
			}
			if item.Creator != "Test Creator" {
				t.Errorf("Expected item creator 'Test Creator', got '%s'", item.Creator)
			}
		}
	}
}

func TestParseRealFeed(t *testing.T) {
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
