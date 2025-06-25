package storage

import (
	"os"
	"testing"

	"github.com/kdimonych/go_douuarss/lib/rss"
)

func TestSaveCannel(t *testing.T) {
	t.Skip("This test requires a live database connection and is skipped by default")

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

	channels, err := rss.Parse(blob)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(channels) == 0 {
		t.Fatal("Expected at least one channel")
	}

	dbURL := "postgres://go_douuarss:go_douuarss@localhost:5432/go_douuarss?sslmode=disable"
	s, err := NewStorage(dbURL)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer s.Close()

	id, err := s.InsertOrMergeChannel(&channels[0])
	if err != nil {
		t.Fatalf("Failed to insert or merge channel: %v", err)
	}
	if id == 0 {
		t.Fatal("Expected a valid channel ID, got 0")
	}
	t.Logf("Channel inserted or merged with ID: %d", id)
}

func TestSaveCannelWithRealData(t *testing.T) {
	t.Skip("This test requires a live database connection and is skipped by default")

	data, err := os.ReadFile("../rss/testdata/test_real_rss_blob.xml")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	channels, err := rss.Parse(data)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(channels) == 0 {
		t.Fatal("Expected at least one channel")
	}

	dbURL := "postgres://go_douuarss:go_douuarss@localhost:5432/go_douuarss?sslmode=disable"
	s, err := NewStorage(dbURL)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer s.Close()

	id, err := s.InsertOrMergeChannel(&channels[0])
	if err != nil {
		t.Fatalf("Failed to insert or merge channel: %v", err)
	}
	if id == 0 {
		t.Fatal("Expected a valid channel ID, got 0")
	}
	t.Logf("Channel inserted or merged with ID: %d", id)
}
