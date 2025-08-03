package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kdimonych/go_douuarss/internal/common"
	"github.com/kdimonych/go_douuarss/internal/news_service"
	"github.com/spf13/cobra"
)

func mainServiceRun(dbURL, migrationsDir string) {
	if dbURL == "" {
		log.Panic("DATABASE_URL is not set")
	}
	if migrationsDir == "" {
		migrationsDir = "/app/internal/storage/migrations"
		log.Printf("[Error] MIGRATIONS_DIR is not set. Use default value: %s", migrationsDir)
	}

	config := &news_service.Config{
		DatabaseURL:   dbURL,
		MigrationsDir: migrationsDir,
	}

	ctx, cancel := context.WithCancel(context.Background())

	service, err := news_service.NewNewsServiceBuilder().Build(ctx, config)
	if err != nil {
		log.Panicf("[Panic] Unable to initialize news service: %v", err)
	}

	if err := service.Init(); err != nil {
		log.Panicf("[Panic] Unable to initialize news service: %v", err)
	}

	waitGroup := &sync.WaitGroup{}

	if err := service.Start(waitGroup, ctx); err != nil {
		log.Panicf("[Panic] Unable to start news service: %v", err)
	}

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("[Info] Shutdown signal received")
		cancel() // cancel context to stop service
	}()

	waitGroup.Wait()
	service.Stop()
}

func registerFeedRun(dbURL, migrationsDir, feedURL string) {
	if dbURL == "" {
		log.Panic("DATABASE_URL is not set")
	}

	if migrationsDir == "" {
		migrationsDir = "/app/migrations"
		log.Printf("[Error] MIGRATIONS_DIR is not set. Use default value: %s", migrationsDir)
	}

	config := &news_service.Config{
		DatabaseURL:   dbURL,
		MigrationsDir: migrationsDir,
	}

	service, err := news_service.NewNewsServiceBuilder().Build(context.Background(), config)
	if err != nil {
		log.Panicf("[Panic] Unable to initialize news service: %v", err)
	}

	if err = service.Init(); err != nil {
		log.Panicf("[Panic] Unable to initialize news service: %v", err)
	}

	feedId, err := service.RegisterRssFeed(context.Background(), feedURL)
	if err != nil {
		log.Panicf("[Panic] Unable to register RSS feed: %v", err)
	}
	log.Printf("[Info] Registered RSS feed with ID: %d", feedId)

	service.Stop()
}

func main() {
	var dbURL string
	var migrationsDir string
	var feedURL string

	dbURL = os.Getenv("DATABASE_URL")
	migrationsDir = os.Getenv("MIGRATIONS_DIR")

	rootCmd := &cobra.Command{
		Use:   "news-checker",
		Short: "News Checker Service",
		Run: func(_ *cobra.Command, _ []string) {
			mainServiceRun(dbURL, migrationsDir)
		},
	}

	rootCmd.PersistentFlags().StringVar(&dbURL, "db", os.Getenv("DATABASE_URL"), "Database URL")
	rootCmd.PersistentFlags().StringVar(&migrationsDir, "migrations", os.Getenv("MIGRATIONS_DIR"), "Migrations directory")

	addCmd := &cobra.Command{
		Use:   "add-feed",
		Short: "Add a new RSS feed",
		Run: func(_ *cobra.Command, _ []string) {
			if feedURL == "" {
				log.Panic("Feed URL is required")
			}

			registerFeedRun(dbURL, migrationsDir, feedURL)
		},
	}
	addCmd.Flags().StringVar(&feedURL, "url", "", "RSS Feed URL")
	if err := addCmd.MarkFlagRequired("url"); err != nil {
		log.Panicf("Failed to mark flag as required: %v", err)
	}

	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Panicf("Command failed: %v", common.UnwrapAll(err))
	}
}
