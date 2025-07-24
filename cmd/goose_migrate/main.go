package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kdimonych/go_douuarss/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'up' or 'down' subcommands")
		os.Exit(1)
	}

	var migrationsDir string
	flag.StringVar(&migrationsDir, "migrations", "./migrations", "Path to migration files")

	cmd := os.Args[1]

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	migrate := func() int {
		m, err := storage.NewMigratorBuilder().Build(context.Background(), dbURL, migrationsDir)
		if err != nil {
			log.Fatalf("failed to create migrattor: %v", err)
		}
		defer m.Close()

		switch cmd {
		case "up":
			if err := m.Up(); err != nil {
				log.Printf("[Error] goose up failed: %v\n", err)
				return 1
			}
			log.Println("[Info] Migrations applied successfully!")
		case "down":
			if err := m.Down(); err != nil {
				log.Printf("[Error] goose down failed: %v\n", err)
				return 1
			}
			log.Println("[Info] Migration rolled back successfully!")
		default:
			fmt.Println("[Warning] expected 'up' or 'down' subcommands")
			return 1
		}
		return 0
	}

	os.Exit(migrate())
}
