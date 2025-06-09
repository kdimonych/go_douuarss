package rss

import (
	"os"
	"testing"
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
				<title>Test Item</title>
				<link>http://example.com/item</link>
				<description>Test Item Description</description>
				<pubDate>Mon, 09 Jun 2025 11:40:06 +0300</pubDate>
				<creator>Test Creator</creator>
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
