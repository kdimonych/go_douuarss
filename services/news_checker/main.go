package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kdimonych/go_douuarss/lib/rss"
	"github.com/kdimonych/go_douuarss/lib/storage"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Panic("DATABASE_URL is not set")
	}

	provider := rss.StartRssProvider(context.Background())
	defer provider.Close()

	s, err := storage.NewStorage(dbURL)
	if err != nil {
		log.Panic("Unable to initialize storage: %w", err)
	}
	defer s.Close()

	c := provider.GetChannel()

	for channel := range c {
		// Here you can process the channel received from the RSS provider
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Printf("Received Channel: %s\n", channel.Title)
		fmt.Printf("Link: %s\n", channel.Link)
		fmt.Printf("Description: %s\n", channel.Description)
		fmt.Printf("Language: %s\n", channel.Language)
		fmt.Printf("Last Build Date: %v\n", channel.LastBuildDate)

		fmt.Println("--------------------------------------------------")
		for _, item := range channel.Items {
			fmt.Printf("Item Title: %s\n", item.Title)
			fmt.Printf("Item Link: %s\n", item.Link)
			fmt.Printf("Item PubDate: %v\n", item.PubDate)
			fmt.Printf("Item Creator: %s\n", item.Creator)
			fmt.Printf("Item Description:\n%s\n", item.Description)
			fmt.Println("--------------------------------------------------")
		}

		if _, err := s.InsertOrMergeChannel(&channel); err != nil {
			log.Printf("Unable to insert or merge channel %s: %v", channel.Title, err)
			continue
		}
		// Here you can insert the channel and items into the database
	}
}
