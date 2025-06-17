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
	flag.CommandLine.Parse(os.Args[2:])

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	switch cmd {
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatalf("goose up failed: %v", err)
		}
		log.Println("Migrations applied successfully!")
	case "down":
		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatalf("goose down failed: %v", err)
		}
		log.Println("Migration rolled back successfully!")
	default:
		fmt.Println("expected 'up' or 'down' subcommands")
		os.Exit(1)
	}
}
