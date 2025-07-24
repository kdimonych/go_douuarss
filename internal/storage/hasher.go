package storage

import "strings"

type ChannelHasher interface {
	Hash(channel *Channel) string
}

type ItemHasher interface {
	// Hash returns the hash of the given data.
	Hash(item *Item) string
}

func NewChannelHasher() ChannelHasher {
	return &defaultChannelHasherImpl{}
}

func NewItemHasher() ItemHasher {
	return &defaultChannelItemImpl{}
}

// ==================== Private ====================

type defaultChannelHasherImpl struct{}

func (*defaultChannelHasherImpl) Hash(channel *Channel) string {
	hash := strings.ToLower(channel.Title) + "_" + channel.LastBuildDate.String()
	return hash
}

type defaultChannelItemImpl struct{}

func (*defaultChannelItemImpl) Hash(item *Item) string {
	hash := strings.ToLower(item.Title) + "_" + item.PubDate.String()
	return hash
}
