package main

import (
	"fmt"
	"time"

	//"github.com/jackc/pgx/v5"
	"github.com/kdimonych/go_douuarss/pkg/rss"
)

func main() {
	// dbURL := os.Getenv("DATABASE_URL")
	// if dbURL == "" {
	// 	log.Fatal("DATABASE_URL is not set")
	// }

	// conn, err := pgx.Connect(context.Background(), dbURL)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(context.Background())

	c := make(chan rss.Channel)
	go rss_provider(c)

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

func rss_provider(c chan rss.Channel) {
	url := "https://dou.ua/feed/"
	// This function is a placeholder for the RSS provider logic.
	// It can be implemented to fetch and process RSS feeds.
	fmt.Println("RSS provider function called")

	for {
		channels, err := rss.FetchAndParse(url)
		if err != nil {
			fmt.Printf("Error fetching and parsing RSS feed: %v\n", err)
			continue
		}

		if len(channels) == 0 {
			fmt.Println("No channels found in the RSS feed")
			continue
		}

		for _, channel := range channels {
			c <- channel // Sending the RSS channel to the channel
		}

		fmt.Println("Sleep for 10 seconds before fetching again")
		time.Sleep(10000 * time.Millisecond)
	}

}
