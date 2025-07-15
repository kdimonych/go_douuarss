package news_service

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewsServiceBuilder(t *testing.T) {
	t.Skip("This test requires a live database and is skipped by default")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testFeedUrl := "https://dou.ua/feed/"

	config := &Config{
		DatabaseURL:   testDBURL,
		MigrationsDir: testMigrationsDir,
	}

	service, err := NewNewsServiceBuilder().Build(config)
	if err != nil {
		log.Panicf("Unable to initialize news service: %v", err)
	}

	if err := service.Init(); err != nil {
		log.Panicf("Unable to initialize news service: %v", err)
	}

	feedId, Err := service.AddRssFeed(testFeedUrl)
	if Err != nil {
		log.Panicf("Unable to add RSS feed: %v", Err)
	}
	log.Printf("Added RSS feed with ID: %d", feedId)

	waitGroup := &sync.WaitGroup{}

	if err := service.Start(waitGroup, ctx); err != nil {
		log.Panicf("Unable to start news service: %v", err)
	}

	waitGroup.Wait()
	service.Stop()
}
