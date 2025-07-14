package rss

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRealRssProvider(t *testing.T) {
	t.Skip("This test requires a live RSS feed and is skipped by default")
	rssUrl := "https://dou.ua/feed/"
	testProviderId := RssProviderId(1)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	rssProviderFabric := NewRssProviderFabric()

	channel := make(chan RssMessage, 10)
	provider, err := rssProviderFabric.CreateRssProvider(testProviderId, rssUrl, channel)
	if err != nil {
		t.Fatalf("failed to create RealRssProvider: %v", err)
	}
	defer provider.Stop()
	err = provider.Start(wg, ctx)
	if err != nil {
		t.Fatalf("failed to start RealRssProvider: %v", err)
	}

	if provider.Url() != rssUrl {
		t.Fatalf("expected associated URL %s, got %s", rssUrl, provider.Url())
	}

	if provider.Id() != testProviderId {
		t.Fatalf("expected provider ID %v, got %v", testProviderId, provider.Id())
	}

	select {
	case msg := <-channel:
		fmt.Printf("Received channel: %s\n", msg.Channel.Title)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for channel message")
	}
}
