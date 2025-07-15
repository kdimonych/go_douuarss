package rss

import (
	"fmt"
	"log"
	"sync"
)

const (
	MessageQueueSize = 10
)

type RssProviderInfo struct {
	Id     RssProviderId
	Url    string
	Active bool
}

type RssClientService interface {
	AddRssFeedProvider(id RssProviderId, rssUrl string) error
	RemoveRssFeedProvider(id RssProviderId)
	ProviderInfo(id RssProviderId) *RssProviderInfo
	ListProviders() []RssProviderInfo
	GetRssMessage() <-chan RssMessage
	Stop() // Stop stops all RSS providers and cleans up resources.
}

type rssClientServiceImpl struct {
	rssProviderFabric RssProviderFabric
	messageChannel    chan RssMessage

	externalWg *sync.WaitGroup

	rssProviders map[RssProviderId]RssProvider
}

type RssClientServiceBuilder interface {
	WithRssProviderFabric(rssProviderFabric RssProviderFabric) RssClientServiceBuilder
	Build() (RssClientService, error)
}

type rssClientServiceBuilder struct {
	rssProviderFabric RssProviderFabric
}

func (b *rssClientServiceBuilder) WithRssProviderFabric(rssProviderFabric RssProviderFabric) RssClientServiceBuilder {
	if rssProviderFabric == nil {
		log.Println("Nil RssProviderFabric provided. Using default RssProviderFabric")
		rssProviderFabric = NewRssProviderFabric()
	}
	b.rssProviderFabric = rssProviderFabric
	return b
}

func (b *rssClientServiceBuilder) Build() (RssClientService, error) {
	return &rssClientServiceImpl{
		rssProviders:      make(map[RssProviderId]RssProvider),
		messageChannel:    make(chan RssMessage, MessageQueueSize),
		externalWg:        &sync.WaitGroup{},
		rssProviderFabric: b.rssProviderFabric,
	}, nil
}

func NewRssClientServiceBuilder() RssClientServiceBuilder {
	return &rssClientServiceBuilder{
		rssProviderFabric: NewRssProviderFabric(),
	}
}

func (service *rssClientServiceImpl) AddRssFeedProvider(id RssProviderId, rssUrl string) error {
	existedProvider, exists := service.rssProviders[id]
	if exists && existedProvider.IsActive() {
		return nil
	} else if exists {
		log.Printf("The provider %v exists but not active, remove it and try to add again", id)
		existedProvider.Stop()
		delete(service.rssProviders, id)
	}

	provider, err := service.rssProviderFabric.CreateRssProvider(id, rssUrl, service.messageChannel)
	if err != nil {
		return fmt.Errorf("failed to create RSS provider: %w", err)
	}

	err = provider.Start(service.externalWg, nil)
	if err != nil {
		log.Printf("Failed to start RSS provider %v: %v\n", provider.Id(), err)
		return fmt.Errorf("failed to start RSS provider: %w", err)
	}

	service.rssProviders[id] = provider
	log.Printf("Added new RSS provider with ID %v: %s\n", id, rssUrl)
	return err
}

func (service *rssClientServiceImpl) RemoveRssFeedProvider(id RssProviderId) {
	provider, exists := service.rssProviders[id]
	if exists {
		provider.Stop()
		delete(service.rssProviders, id)
	}
}

func (service *rssClientServiceImpl) ProviderInfo(id RssProviderId) *RssProviderInfo {
	provider, exists := service.rssProviders[id]
	if !exists {
		return nil
	}

	return &RssProviderInfo{
		Id:     provider.Id(),
		Url:    provider.Url(),
		Active: provider.IsActive(),
	}
}

func (service *rssClientServiceImpl) ListProviders() []RssProviderInfo {
	var providers []RssProviderInfo
	for id, provider := range service.rssProviders {
		providers = append(providers, RssProviderInfo{
			Id:     id,
			Url:    provider.Url(),
			Active: provider.IsActive(),
		})
	}
	return providers
}

func (service *rssClientServiceImpl) GetRssMessage() <-chan RssMessage {
	return service.messageChannel
}

func (*rssClientServiceImpl) Init() error {
	return nil
}

func (service *rssClientServiceImpl) Stop() {
	for _, provider := range service.rssProviders {
		provider.Stop()
		log.Printf("Stopped RSS provider %v\n", provider.Id())
	}

	service.externalWg.Wait() // Wait for all providers to finish

	close(service.messageChannel)
	service.rssProviders = make(map[RssProviderId]RssProvider)

	log.Println("RSS fetch service stopped successfully")
}
