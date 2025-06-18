package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}
		defer db.Close()

		switch cmd {
		case "up":
			if err := goose.Up(db, migrationsDir); err != nil {
				log.Printf("goose up failed: %v", err)
				return 1
			}
			log.Println("Migrations applied successfully!")
		case "down":
			if err := goose.Down(db, migrationsDir); err != nil {
				log.Printf("goose down failed: %v", err)
				return 1
			}
			log.Println("Migration rolled back successfully!")
		default:
			fmt.Println("expected 'up' or 'down' subcommands")
			return 1
		}
		return 0
	}

	os.Exit(migrate())
}
