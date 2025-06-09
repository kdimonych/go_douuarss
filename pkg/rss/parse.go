package rss

import (
	"encoding/xml"
	"fmt"
	"time"
)

// UnmarshalText implements for RFC1123Z time format
func (ct *RFC1123ZDate) UnmarshalText(text []byte) error {
	t, err := time.Parse(time.RFC1123Z, string(text))
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func Parse(blob []byte) ([]Channel, error) {
	var rss Xml
	err := xml.Unmarshal(blob, &rss)
	if err != nil {
		return nil, fmt.Errorf("Parse error: %w", err)
	}

	return rss.Channel, nil
}

func FetchAndParse(url string) ([]Channel, error) {
	data, fetchErr := Fetch(url)
	if fetchErr != nil {
		return nil, fmt.Errorf("Fetch error: %w", fetchErr)
	}

	channels, parseErr := Parse(data)
	if parseErr != nil {
		return nil, parseErr
	}

	if len(channels) == 0 {
		return nil, &FetchError{ErrorCodeNoData, "No channels found in the RSS feed"}
	}

	return channels, nil
}
