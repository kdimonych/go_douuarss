package storage

import "time"

type FeedId int64
type ChannelId int64
type ItemId int64

type Feed struct {
	Id       FeedId    // Unique identifier for the feed
	Name     string    // The feed name
	Url      string    // The feed URL
	Channels []Channel // List of channels in the RSS feed
	Hash     string
}

type Channel struct {
	Id            ChannelId // Unique identifier for the channel
	Title         string    // Title of the RSS channel
	Link          string    // Link to the RSS channel
	Description   string    // Description of the RSS channel
	Language      string    // Language of the RSS channel
	LastBuildDate time.Time // Last build date of the RSS channel
	Items         []Item    // List of items in the RSS channel
	Hash          string    // Hash of the channel for deduplication
}

type Item struct {
	Id          ItemId    // Unique identifier for the item
	Title       string    // Title of the article
	Link        string    // Link to the article
	Description string    // Description of the article
	PubDate     time.Time // Publication date of the article
	Creator     string    // List of creators of the article
	Hash        string    // Hash of the item for deduplication
}
