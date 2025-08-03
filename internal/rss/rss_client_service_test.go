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

	// Try to add the same provider again
	err = rssClientService.AddRssFeedProvider(testProviderId, rssUrl)
	if err != nil {
		t.Fatalf("failed to add existing RSS feed provider: %v", err)
	}

	providerList := rssClientService.ListProviders()

	if len(providerList) == 0 {
		t.Fatal("expected at least one provider in the list")
	}
	if providerList[0].Id != testProviderId {
		t.Fatalf("expected provider ID %v, got %v", testProviderId, providerList[0].Id)
	}

	if !providerList[0].Active {
		t.Fatal("expected provider to be active after adding")
	}

	providerInfo := rssClientService.ProviderInfo(providerList[0].Id)
	if providerInfo == nil {
		t.Fatal("expected non-nil provider info")
	}

	if providerInfo.Url != providerList[0].Url {
		t.Fatalf("expected provider URL %s, got %s", providerList[0].Url, providerInfo.Url)
	}
	if providerInfo.Id != providerList[0].Id {
		t.Fatalf("expected provider ID %v, got %v", providerList[0].Id, providerInfo.Id)
	}
	if providerInfo.Active != providerList[0].Active {
		t.Fatalf("expected provider active status %v, got %v", providerList[0].Active, providerInfo.Active)
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

	rssClientService.RemoveRssFeedProvider(testProviderId)

	if len(rssClientService.ListProviders()) != 0 {
		t.Fatal("expected no providers after removal")
	}
}
