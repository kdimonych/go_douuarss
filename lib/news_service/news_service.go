package news_service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kdimonych/go_douuarss/lib/rss"
	"github.com/kdimonych/go_douuarss/lib/storage"
)

type Config struct {
	DatabaseURL   string
	MigrationsDir string
}

type NewsService interface {
	Init() error
	Start(externalWg *sync.WaitGroup, parentCtx context.Context) error
	Stop()
}

type newsServiceImpl struct {
	storageSrv       storage.Storage
	rssClientService rss.RssClientService

	cancel     context.CancelFunc
	ctx        context.Context
	wg         *sync.WaitGroup
	externalWg *sync.WaitGroup
}

type NewsServiceBuilder interface {
	WithDbConnectionFabric(dbConnectionFabric storage.DbConectionFabric) NewsServiceBuilder
	WithStorageBuilder(storageBuilder storage.StorageBuilder) NewsServiceBuilder
	WithRssClientServiceBuilder(rssClientServiceBuilder rss.RssClientServiceBuilder) NewsServiceBuilder
	WithMigratorBuilder(migratorBuilder storage.MigratorBuilder) NewsServiceBuilder
	Build(config *Config) (NewsService, error)
}

type newsServiceBuilderImpl struct {
	dbConnectionFabric      storage.DbConectionFabric
	storageBuilder          storage.StorageBuilder
	migratorBuilder         storage.MigratorBuilder
	rssClientServiceBuilder rss.RssClientServiceBuilder
}

func NewNewsServiceBuilder() NewsServiceBuilder {
	return &newsServiceBuilderImpl{
		dbConnectionFabric:      storage.NewDbConnectionFabric(),
		storageBuilder:          storage.NewStorageBuilder(),
		migratorBuilder:         storage.NewMigratorBuilder(),
		rssClientServiceBuilder: rss.NewRssClientServiceBuilder(),
	}
}

func (b *newsServiceBuilderImpl) WithDbConnectionFabric(dbConnectionFabric storage.DbConectionFabric) NewsServiceBuilder {
	if dbConnectionFabric == nil {
		log.Println("Nil DbConnectionFabric provided. Using default DbConnectionFabric")
		dbConnectionFabric = storage.NewDbConnectionFabric()
	}
	b.dbConnectionFabric = dbConnectionFabric
	return b
}

func (b *newsServiceBuilderImpl) WithStorageBuilder(storageBuilder storage.StorageBuilder) NewsServiceBuilder {
	if storageBuilder == nil {
		log.Println("Nil StorageBuilder provided. Using default StorageBuilder")
		storageBuilder = storage.NewStorageBuilder()
	}
	b.storageBuilder = storageBuilder
	return b
}

func (b *newsServiceBuilderImpl) WithRssClientServiceBuilder(rssClientServiceBuilder rss.RssClientServiceBuilder) NewsServiceBuilder {
	if rssClientServiceBuilder == nil {
		log.Println("Nil RssClientServiceBuilder provided. Using default RssClientServiceBuilder")
		rssClientServiceBuilder = rss.NewRssClientServiceBuilder()
	}
	b.rssClientServiceBuilder = rssClientServiceBuilder
	return b
}

func (b *newsServiceBuilderImpl) WithMigratorBuilder(migratorBuilder storage.MigratorBuilder) NewsServiceBuilder {
	if migratorBuilder == nil {
		migratorBuilder = storage.NewMigratorBuilder()
		log.Println("Nil MigratorBuilder provided. Using default MigratorBuilder")
	}
	b.migratorBuilder = migratorBuilder
	return b
}

func (b *newsServiceBuilderImpl) runMigrations(config *Config) error {
	if config.MigrationsDir == "" {
		return fmt.Errorf("migrations directory cannot be empty")
	}

	migrator, err := b.migratorBuilder.
		WithDbConnectionFabric(b.dbConnectionFabric).
		Build(config.DatabaseURL, config.MigrationsDir)
	if err != nil {
		return fmt.Errorf("failed to build migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func validateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if config.DatabaseURL == "" {
		return fmt.Errorf("database URL cannot be empty")
	}
	if config.MigrationsDir == "" {
		return fmt.Errorf("migrations directory cannot be empty")
	}
	return nil
}

func (b *newsServiceBuilderImpl) Build(config *Config) (NewsService, error) {
	// Validate the configuration
	err := validateConfig(config)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Start the migrator to apply any pending migrations
	err = b.runMigrations(config)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	storageSrv, err := b.storageBuilder.
		WithDbConnectionFabric(b.dbConnectionFabric).
		Build(config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}

	rssClientService, err := b.rssClientServiceBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create RSS client service: %w", err)
	}

	return &newsServiceImpl{
		storageSrv:       storageSrv,
		rssClientService: rssClientService,
		wg:               &sync.WaitGroup{},
		externalWg:       nil,
	}, nil
}

// ===== NewsService Implementation =====
func (*newsServiceImpl) Init() error {
	return nil
}

func (srv *newsServiceImpl) Start(externalWg *sync.WaitGroup, parentCtx context.Context) error {
	log.Println("Starting news service...")

	if externalWg == nil {
		return fmt.Errorf("external wait group cannot be nil")
	}

	srv.externalWg = externalWg

	if err := srv.Init(); err != nil {
		return fmt.Errorf("failed to initialize news service: %w", err)
	}

	if parentCtx == nil {
		log.Println("Parent context is nil, using background context")
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)
	srv.ctx = ctx
	srv.cancel = cancel

	feeds, err := srv.storageSrv.GetFeedsOnly()
	if err != nil {
		return fmt.Errorf("unable to get feeds from storage: %w", err)
	}

	log.Printf("Found %d feeds in storage, starting to fetch news...\n", len(feeds))

	// Start fetching news from RSS feeds
	for _, feed := range feeds {
		log.Printf("Starting to fetch news for feed: %s\n", feed.Url)
		err = srv.rssClientService.AddRssFeedProvider(rss.RssProviderId(feed.Id), feed.Url)
		if err != nil {
			log.Printf("Unable to add RSS feed provider for %s (id: %v): %v\n", feed.Url, feed.Id, err)
			continue
		}
	}

	srv.wg.Add(1)
	srv.externalWg.Add(1)
	go srv.worker()

	log.Println("News service started successfully")
	return nil
}

func (srv *newsServiceImpl) Stop() {
	if srv.cancel != nil {
		// Stop the RSS client service before stopping the news service
		srv.rssClientService.Stop()

		srv.cancel()
		srv.wg.Wait()    // Wait for the worker goroutine to finish
		srv.cancel = nil // Clear the cancel function to avoid double cancellation

		srv.storageSrv.Close()
	}

	log.Println("News service stopped successfully")
}

func (srv *newsServiceImpl) processRssMessage(msg *rss.RssMessage) {
	rssChannel := msg.Channel

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("Received Channel: %s\n", rssChannel.Title)
	fmt.Printf("Link: %s\n", rssChannel.Link)
	fmt.Printf("Description: %s\n", rssChannel.Description)
	fmt.Printf("Language: %s\n", rssChannel.Language)
	fmt.Printf("Last Build Date: %v\n", rssChannel.LastBuildDate)

	fmt.Println("--------------------------------------------------")
	for _, item := range rssChannel.Items {
		fmt.Printf("Item Title: %s\n", item.Title)
		fmt.Printf("Item Link: %s\n", item.Link)
		fmt.Printf("Item PubDate: %v\n", item.PubDate)
		fmt.Printf("Item Creator: %s\n", item.Creator)
		fmt.Printf("Item Description:\n%s\n", item.Description)
		fmt.Println("--------------------------------------------------")
	}

	// Here you can insert the channel and items into the database
	channel := storage.ChannelFromRSS(&msg.Channel)

	if _, err := srv.storageSrv.InsertOrMergeChannel(storage.FeedId(msg.Id), &channel); err != nil {
		log.Printf("Unable to insert or merge channel %s: %v\n", channel.Title, err)
		return
	}
}

func (srv *newsServiceImpl) worker() {
	defer func() {
		srv.wg.Done()
		srv.externalWg.Done()
	}()

	for {
		select {
		case <-srv.ctx.Done():
			log.Printf("Worker context canceled, exiting..., %s\n", srv.ctx.Err().Error())
			return
		case msg := <-srv.rssClientService.GetRssMessage():
			srv.processRssMessage(&msg)
		}
	}
}
