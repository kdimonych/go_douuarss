package rss

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kdimonych/go_douuarss/internal/common"
)

type FetchClient interface {
	Fetch(ctx context.Context, urlStr string) ([]byte, error)
}

type fetchClientImpl struct {
	httpClient *http.Client
}

type FetchClientBuilder interface {
	WithHttpClient(client *http.Client) FetchClientBuilder
	Build() FetchClient
}

type fetchClientBuilderImpl struct {
	httpClient *http.Client
}

func NewFetchClientBuilder() FetchClientBuilder {
	return &fetchClientBuilderImpl{
		httpClient: http.DefaultClient,
	}
}

func (b *fetchClientBuilderImpl) WithHttpClient(client *http.Client) FetchClientBuilder {
	if client == nil {
		log.Printf("[Info] FetchClientBuilder: provided http client is nil, using default http client")
		client = http.DefaultClient
	}
	b.httpClient = client
	return b
}

func (b *fetchClientBuilderImpl) Build() FetchClient {
	return &fetchClientImpl{
		httpClient: b.httpClient,
	}
}

func (c *fetchClientImpl) Fetch(ctx context.Context, urlStr string) ([]byte, error) {
	err := common.ValidateURL(urlStr)
	if err != nil {
		return nil, &FetchError{Code: ErrorCodeInvalidUrl, Details: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, &FetchError{Code: ErrorCodeUnreachable, Details: fmt.Errorf("url: %s, err: %w", urlStr, err)}
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &FetchError{Code: ErrorCodeUnreachable, Details: fmt.Errorf("url: %s, err: %w", urlStr, err)}
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, &FetchError{ErrorCodeHttpError, fmt.Errorf("HTTP Status: %s, url: %s", res.Status, urlStr)}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &FetchError{ErrorCodeNoData, err}
	}

	return body, nil
}

// ================ Private methods ===================
