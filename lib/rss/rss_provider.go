package rss

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kdimonych/go_douuarss/lib/common"
)

const (
	chanelBufferSize = 100
	recheckPeriod    = 20 * time.Second
)

type RssProviderId int64

type RssMessage struct {
	Channel Channel
	Id      RssProviderId
}

type RssProvider interface {
	Start(externalWg *sync.WaitGroup, ctx context.Context) error // Start initializes the provider and begins fetching messages.
	Stop()                                                       // Stop stops the provider and cleans up resources.
	Url() string
	Id() RssProviderId
}

type rssProviderImpl struct {
	id         RssProviderId
	messageOut chan<- RssMessage
	cancel     context.CancelFunc
	ctx        context.Context
	wg         *sync.WaitGroup
	externalWg *sync.WaitGroup // Optional external WaitGroup to manage the lifecycle of the connection
	rssUrl     string
}

func (provider *rssProviderImpl) Start(externalWg *sync.WaitGroup, ctx context.Context) error {
	if externalWg == nil {
		return fmt.Errorf("[%v] external wait group is not provided", provider.id)
	}

	if ctx == nil {
		log.Printf("[%v] Parent context is nil, using background context\n", provider.id)
		ctx = context.Background()
	}

	if provider.messageOut == nil {
		return fmt.Errorf("[%v] messageOut channel cannot be nil", provider.id)
	}

	derivedCtx, cancel := context.WithCancel(ctx)
	provider.ctx = derivedCtx
	provider.cancel = cancel
	provider.externalWg = externalWg

	return provider.startRssProvider()
}

func (provider *rssProviderImpl) Stop() {
	if provider.cancel != nil {
		provider.cancel()
		provider.wg.Wait()    // Wait for the worker goroutine to finish
		provider.cancel = nil // Clear the cancel function to avoid double cancellation
	}
}

func (provider *rssProviderImpl) Url() string {
	return provider.rssUrl
}

func (provider *rssProviderImpl) Id() RssProviderId {
	return provider.id
}

type RssProviderFabric interface {
	CreateRssProvider(
		id RssProviderId,
		rssUrl string,
		messageOut chan<- RssMessage) (RssProvider, error)
}

type rssProviderFabric struct{}

func (*rssProviderFabric) CreateRssProvider(
	id RssProviderId,
	rssUrl string,
	messageOut chan<- RssMessage,
) (RssProvider, error) {
	// Validate the URL before creating the provider
	if err := common.ValidateURL(rssUrl); err != nil {
		return nil, err
	}
	if messageOut == nil {
		return nil, fmt.Errorf("messageOut channel cannot be nil")
	}
	rssProviderHandle := &rssProviderImpl{
		id:         id,
		messageOut: messageOut,
		cancel:     nil,
		ctx:        nil,
		wg:         &sync.WaitGroup{},
		externalWg: nil,
		rssUrl:     rssUrl,
	}

	return rssProviderHandle, nil
}

// NewRssProviderFabric creates a new instance of RssProviderFabric.
func NewRssProviderFabric() RssProviderFabric {
	return &rssProviderFabric{}
}

// ================== Private methods ===================

func (provider *rssProviderImpl) tryFetchAndParse() ([]Channel, error) {
	channels, err := FetchAndParse(provider.ctx, provider.rssUrl)
	if err != nil {
		// Handle unrecoverable errors
		if fetchErr, ok := err.(*FetchError); ok {
			if fetchErr.Code == ErrorCodeInvalidUrl {
				log.Printf("[%v] Unrecoverable error. Invalid URL: %s\n", provider.id, provider.rssUrl)
				return nil, fetchErr
			}
		}

		// Log the error and return empty channels
		log.Printf("[%v]  Error during fetching RSS feed: %s\n", provider.id, err.Error())
		return []Channel{}, nil
	}

	return channels, nil
}

func (provider *rssProviderImpl) worker() {
	// Ensure the WaitGroup is decremented when the worker exits
	defer func() {
		provider.wg.Done()
		provider.externalWg.Done()
	}()

	for {
		channels, err := provider.tryFetchAndParse()
		if err != nil {
			// Handle unrecoverable errors
			return
		}

		// Process the fetched channels
		for _, channel := range channels {
			provider.messageOut <- RssMessage{Channel: channel, Id: provider.id}
		}

		log.Printf("[%v] Sleep for: %v", provider.id, recheckPeriod)
		err = sleepWithContext(provider.ctx, recheckPeriod)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("[%v] Context canceled, exiting worker", provider.id)
			}

			if errors.Is(err, context.DeadlineExceeded) {
				log.Printf("[%v] Error during sleep: %s", provider.id, err.Error())
			}
			return
		}
	}
}

// SleepWithContext waits for the given duration or returns early if the context is canceled.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil // Slept the full duration
	case <-ctx.Done():
		return ctx.Err() // Woken up by context cancellation
	}
}

func (provider *rssProviderImpl) startRssProvider() error {
	// Use a WaitGroup to manage the worker goroutine
	provider.wg.Add(1)
	provider.externalWg.Add(1)

	// Start the worker goroutine to fetch RSS feeds
	go provider.worker()
	return nil
}
