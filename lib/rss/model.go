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
	Title         string       `xml:"title"`         // Title of the RSS channel
	Link          string       `xml:"link"`          // Link to the RSS channel
	Description   string       `xml:"description"`   // Description of the RSS channel
	Language      string       `xml:"language"`      // Language of the RSS channel
	LastBuildDate RFC1123ZDate `xml:"lastBuildDate"` // Last build date of the RSS channel
	Items         []Item       `xml:"item"`          // List of items in the RSS channel
}

type Item struct {
	Title       string       `xml:"title"`       // Title of the article
	Link        string       `xml:"link"`        // Link to the article
	Description string       `xml:"description"` // Description of the article
	PubDate     RFC1123ZDate `xml:"pubDate"`     // Publication date of the article
	Creator     string       `xml:"creator"`     // List of creators of the article
}

type RFC1123ZDate struct {
	time.Time
}
