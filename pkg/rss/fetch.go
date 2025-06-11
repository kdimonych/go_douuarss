package rss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Fetch(urlStr string) ([]byte, error) {
	if urlStr == "" {
		urlStr = "https://dou.ua/feed/"
	}

	parsed, err := url.Parse(urlStr)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, FetchError{Code: ErrorCodeUnreachable, Details: fmt.Errorf("invalid URL: %s", urlStr)}
	}

	ctx := context.Background() // or use a timeout: context.WithTimeout(...)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, FetchError{Code: ErrorCodeUnreachable, Details: err}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, FetchError{Code: ErrorCodeUnreachable, Details: err}
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, FetchError{ErrorCodeHttpError, fmt.Errorf("HTTP Status: %s", res.Status)}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, FetchError{ErrorCodeNoData, err}
	}

	return body, nil
}
