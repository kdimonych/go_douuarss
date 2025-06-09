package rss

import (
	"encoding/xml"
	"time"
)

type Xml struct {
	XMLName xml.Name  `xml:"rss"`
	Channel []Channel `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"channel"`       // Title of the RSS channel
	Link          string `xml:"link"`          // Link to the RSS channel
	Description   string `xml:"description"`   // Description of the RSS channel
	Language      string `xml:"language"`      // Language of the RSS channel
	LastBuildDate string `xml:"lastBuildDate"` // Last build date of the RSS channel
	Items         []Item `xml:"item"`          // List of items in the RSS channel
}

type Item struct {
	Title       string    `xml:"title"`       // Title of the article
	Link        string    `xml:"link"`        // Link to the article
	Description string    `xml:"description"` // Description of the article
	PubDate     PubDate   `xml:"pubDate"`     // Publication date of the article
	Creators    []Creator `xml:"creator"`     // List of creators of the article
}

type PubDate struct {
	time.Time
}

type Creator struct {
	Name string `xml:"__text"` // Name of the creator
}
