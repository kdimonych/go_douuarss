package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/kdimonych/go_douuarss/lib/news_service"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Panic("DATABASE_URL is not set")
	}

	config := &news_service.Config{
		DatabaseURL:   dbURL,
		MigrationsDir: "./migrations",
	}

	service, err := news_service.NewNewsServiceBuilder().Build(config)
	if err != nil {
		log.Panicf("Unable to initialize news service: %v", err)
	}

	if err := service.Init(); err != nil {
		log.Panicf("Unable to initialize news service: %v", err)
	}

	waitGroup := &sync.WaitGroup{}

	if err := service.Start(waitGroup, context.Background()); err != nil {
		log.Panicf("Unable to start news service: %v", err)
	}

	waitGroup.Wait()

	service.Stop()
}
