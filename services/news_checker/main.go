package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/kdimonych/go_douuarss/lib/common"
	"github.com/kdimonych/go_douuarss/lib/news_service"
	"github.com/spf13/cobra"
)

func main_service(dbURL, migrationsDir string) {
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

	service, err := news_service.NewNewsServiceBuilder().Build(config)
	if err != nil {
		log.Panicf("[Panic] Unable to initialize news service: %v", err)
	}

	if err := service.Init(); err != nil {
		log.Panicf("[Panic] Unable to initialize news service: %v", err)
	}

	waitGroup := &sync.WaitGroup{}
	ctx := context.Background()

	if err := service.Start(waitGroup, ctx); err != nil {
		log.Panicf("[Panic] Unable to start news service: %v", err)
	}

	waitGroup.Wait()
	service.Stop()
}

func initCmd() *cobra.Command {
	var dbURL string
	var migrationsDir string

	dbURL = os.Getenv("DATABASE_URL")
	migrationsDir = os.Getenv("MIGRATIONS_DIR")

	rootCmd := &cobra.Command{
		Use:   "news-checker",
		Short: "News Checker Service",
		Run: func(_ *cobra.Command, _ []string) {
			main_service(dbURL, migrationsDir)
		},
	}

	rootCmd.Flags().StringVar(&dbURL, "db", os.Getenv("DATABASE_URL"), "Database URL")
	rootCmd.Flags().StringVar(&migrationsDir, "migrations", os.Getenv("MIGRATIONS_DIR"), "Migrations directory")

	return rootCmd
}

func main() {
	rootCmd := initCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Panicf("Command failed: %v", common.UnwrapAll(err))
	}
}
