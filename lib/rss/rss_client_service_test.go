package rss

import (
	"fmt"
	"testing"
	"time"
)

func TestRealRssClientService(t *testing.T) {
	t.Skip("This test requires a live RSS feed and is skipped by default")

	rssUrl := "https://dou.ua/feed/"
	testProviderId := RssProviderId(1)

	serviceBuilder := NewRssClientServiceBuilder()
	rssClientService, err := serviceBuilder.Build()
	if err != nil {
		t.Fatalf("failed to create RssClientService: %v", err)
	}
	err = rssClientService.AddRssFeedProvider(testProviderId, rssUrl)
	if err != nil {
		t.Fatalf("failed to add RSS feed provider: %v", err)
	}

	defer rssClientService.Stop()

	providerList, err := rssClientService.ListProviders()
	if err != nil {
		t.Fatalf("failed to list providers: %v", err)
	}

	if len(providerList) == 0 {
		t.Fatal("expected at least one provider in the list")
	}
	if providerList[0].Id != testProviderId {
		t.Fatalf("expected provider ID %v, got %v", testProviderId, providerList[0].Id)
	}

	channel := rssClientService.GetRssMessage()
	if channel == nil {
		t.Fatal("expected non-nil message channel")
	}

	select {
	case msg := <-channel:
		fmt.Printf("Received channel: %s\n", msg.Channel.Title)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for channel message")
	}

	err = rssClientService.RemoveRssFeedProvider(testProviderId)
	if err != nil {
		t.Fatalf("failed to remove RSS feed provider: %v", err)
	}
}
