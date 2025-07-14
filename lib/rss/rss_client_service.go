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
	Id  RssProviderId
	Url string
}

type RssClientService interface {
	AddRssFeedProvider(id RssProviderId, rssUrl string) error
	RemoveRssFeedProvider(id RssProviderId) error
	ListProviders() ([]RssProviderInfo, error)
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
	if _, exists := service.rssProviders[id]; exists {
		return fmt.Errorf("RSS provider with id %v already exists", id)
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
	return err
}

func (service *rssClientServiceImpl) RemoveRssFeedProvider(id RssProviderId) error {
	provider, exists := service.rssProviders[id]
	if !exists {
		return fmt.Errorf("RSS provider with id %v does not exist", id)
	}

	provider.Stop()
	delete(service.rssProviders, id)
	return nil
}

func (service *rssClientServiceImpl) ListProviders() ([]RssProviderInfo, error) {
	var providers []RssProviderInfo
	for id, provider := range service.rssProviders {
		providers = append(providers, RssProviderInfo{
			Id:  id,
			Url: provider.Url(),
		})
	}
	return providers, nil
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
