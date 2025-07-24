package storage

import (
	"github.com/kdimonych/go_douuarss/internal/rss"
)

func ItemFromRSS(rssItem *rss.Item) Item {
	return Item{
		// Id will be set by the storage
		Title:       rssItem.Title,
		Link:        rssItem.Link,
		Description: rssItem.Description,
		PubDate:     rssItem.PubDate.Time,
		Creator:     rssItem.Creator,
		// Hash will be set by the storage
	}
}

func ChannelFromRSS(rssChannel *rss.Channel) Channel {
	items := make([]Item, len(rssChannel.Items))
	for i, item := range rssChannel.Items {
		items[i] = ItemFromRSS(&item)
	}

	return Channel{
		// Id will be set by the storage
		Title:         rssChannel.Title,
		Link:          rssChannel.Link,
		Description:   rssChannel.Description,
		Language:      rssChannel.Language,
		LastBuildDate: rssChannel.LastBuildDate.Time,
		Items:         items,
		// Hash will be set by the storage
	}
}

func ChannelsFromRSS(rssChannels []rss.Channel) []Channel {
	channels := make([]Channel, len(rssChannels))
	for i, rssChannel := range rssChannels {
		channels[i] = ChannelFromRSS(&rssChannel)
	}
	return channels
}
