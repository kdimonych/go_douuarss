package rss

import (
	"fmt"
	"io"
	"net/http"
)

func Fetch(url string) ([]byte, *FetchError) {
	if url == "" {
		url = "https://dou.ua/feed/"
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, &FetchError{ErrorCodeUnreachable, err.Error()}
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, &FetchError{ErrorCodeHttpError, fmt.Sprintf("HTTP error: %s", res.Status)}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &FetchError{ErrorCodeNoData, err.Error()}
	}

	return body, nil
}
