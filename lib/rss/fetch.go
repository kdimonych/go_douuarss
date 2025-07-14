package rss

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kdimonych/go_douuarss/lib/common"
)

func Fetch(ctx context.Context, urlStr string) ([]byte, error) {
	err := common.ValidateURL(urlStr)
	if err != nil {
		return nil, &FetchError{Code: ErrorCodeInvalidUrl, Details: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, &FetchError{Code: ErrorCodeUnreachable, Details: err}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &FetchError{Code: ErrorCodeUnreachable, Details: err}
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, &FetchError{ErrorCodeHttpError, fmt.Errorf("HTTP Status: %s", res.Status)}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &FetchError{ErrorCodeNoData, err}
	}

	return body, nil
}

// ================ Private methods ===================
