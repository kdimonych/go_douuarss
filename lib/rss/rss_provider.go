package rss

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	chanelBufferSize = 100
	rssUrl           = "https://dou.ua/feed/"
	recheckPeriod    = 20 * time.Second
)

type RssProviderConnection struct {
	channel chan Channel
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (connection *RssProviderConnection) Close() {
	// Close the channel to signal that no more data will be sent
	if connection.cancel != nil {
		connection.cancel()
		connection.wg.Wait() // Wait for the worker goroutine to finish
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

func worker(ctx context.Context, wg *sync.WaitGroup, connection *RssProviderConnection) {
	defer wg.Done()
	for {
		// Check if the context is done before proceeding to fetch RSS feeds
		// If the context is done, close the channel and exit the worker
		// This allows the worker to gracefully exit when the context is canceled
		select {
		case <-ctx.Done():
			close(connection.channel)
			fmt.Printf("Terminated\n")
			return
		default:
			// Continue fetching RSS feeds until the context is done
		}

		channels, err := FetchAndParse(ctx, rssUrl)
		if err != nil {
			// Log the error and continue to the next iteration
			fmt.Printf("Error fetching RSS feed: %v\n", err)
			continue
		}

		for _, channel := range channels {
			connection.channel <- channel
		}

		fmt.Printf("====>> Sleep for: %v\n", recheckPeriod)
		_ = sleepWithContext(ctx, recheckPeriod)
	}
}

func StartRssProvider(parent_ctx context.Context) *RssProviderConnection {
	// Initialize the RssProviderConnection with a buffered channel
	ctx, cancel := context.WithCancel(parent_ctx)

	connection := &RssProviderConnection{
		channel: make(chan Channel, chanelBufferSize),
		cancel:  cancel,
		wg:      sync.WaitGroup{},
	}
	// Use a WaitGroup to manage the worker goroutine
	connection.wg.Add(1)
	// Start the worker goroutine to fetch RSS feeds
	go worker(ctx, &connection.wg, connection)

	return connection
}

func (connection *RssProviderConnection) GetChannel() <-chan Channel {
	return connection.channel
}
