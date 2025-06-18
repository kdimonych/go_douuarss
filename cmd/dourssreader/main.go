package main

import (
	"context"
	"fmt"

	"github.com/kdimonych/go_douuarss/lib/rss"
)

func main() {
	url := "https://dou.ua/feed/"
	ctx := context.Background() // or use a timeout: context.WithTimeout(...)
	channels, err := rss.FetchAndParse(ctx, url)
	if err != nil {
		fmt.Printf("Error fetching and parsing RSS feed: %v\n", err)
		return
	}

	for _, channel := range channels {
		fmt.Printf("Channel Title: %s\n", channel.Title)
		fmt.Printf("Channel Link: %s\n", channel.Link)
		fmt.Printf("Channel Description: %s\n", channel.Description)
		fmt.Printf("Channel Language: %s\n", channel.Language)
		fmt.Printf("Last Build Date: %v\n", channel.LastBuildDate)

		fmt.Println("--------------------------------------------------")
		for _, item := range channel.Items {
			fmt.Printf("Item Title: %s\n", item.Title)
			fmt.Printf("Item Link: %s\n", item.Link)
			fmt.Printf("Item PubDate: %v\n", item.PubDate)
			fmt.Printf("Item Creator: %s\n", item.Creator)
			fmt.Printf("Item Description: \n%s\n", item.Description)
			fmt.Println("--------------------------------------------------")
		}
	}
}
