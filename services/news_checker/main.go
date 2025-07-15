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

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "/app/migrations"
		log.Printf("MIGRATIONS_DIR is not set. Use default value: %s", migrationsDir)
	}

	config := &news_service.Config{
		DatabaseURL:   dbURL,
		MigrationsDir: migrationsDir,
	}

	service, err := news_service.NewNewsServiceBuilder().Build(config)
	if err != nil {
		log.Panicf("Unable to initialize news service: %v", err)
	}

	if err := service.Init(); err != nil {
		log.Panicf("Unable to initialize news service: %v", err)
	}

	waitGroup := &sync.WaitGroup{}
	ctx := context.Background()

	if err := service.Start(waitGroup, ctx); err != nil {
		log.Panicf("Unable to start news service: %v", err)
	}

	waitGroup.Wait()

	service.Stop()
}
