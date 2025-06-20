package main

import (
	"context"
	"fmt"

	"github.com/kdimonych/go_douuarss/lib/rss"
)

func main() {
	provider := rss.StartRssProvider(context.Background())
	defer provider.Close()

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

		// Here you can insert the channel and items into the database
	}
}
