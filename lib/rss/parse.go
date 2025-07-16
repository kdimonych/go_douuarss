package rss

import (
	"context"
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
		return nil, &FetchError{ErrorCodeInvalidData, fmt.Errorf("parse error: %w", err)}
	}

	return rss.Channel, nil
}

func FetchAndParse(client FetchClient, ctx context.Context, url string) ([]Channel, error) {
	if client == nil {
		return nil, &FetchError{ErrorCodeInternalError, fmt.Errorf("fetch client is nil")}
	}
	data, fetchErr := client.Fetch(ctx, url)
	if fetchErr != nil {
		return nil, fetchErr
	}

	channels, parseErr := Parse(data)
	if parseErr != nil {
		return nil, parseErr
	}

	return channels, nil
}
