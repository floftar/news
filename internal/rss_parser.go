package internal

import (
	"encoding/xml"
	"fmt"
	"html"
	"sort"
	"strings"
	"time"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

type Feed struct {
	Title string
	Items []FeedItem
}

type FeedItem struct {
	Title     string
	Link      string
	Published time.Time
}

func ParseFeed(body string) (Feed, error) {
	rss := Rss{}

	decoder := xml.NewDecoder(strings.NewReader(body))
	err := decoder.Decode(&rss)
	if err != nil {
		return Feed{}, fmt.Errorf("parseFeed error: %v", err)
	}

	return Feed{
		Title: rss.Channel.Title,
		Items: sortItemsByPublished(parseItems(rss.Channel.Items)),
	}, nil
}

func parseItems(items []Item) []FeedItem {
	result := make([]FeedItem, 0, len(items))

	for _, item := range items {
		pubDate, err := parseDate(item.PubDate)

		if err != nil {
			fmt.Printf("parseItems error: %v", err.Error())
		}

		result = append(result, FeedItem{
			Title:     html.UnescapeString(item.Title),
			Link:      item.Link,
			Published: pubDate,
		})
	}

	return result
}

func sortItemsByPublished(items []FeedItem) []FeedItem {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Published.After(items[j].Published)
	})

	return items
}

var layouts = [...]string{
	time.RFC1123Z,
	"Mon, 2 Jan 2006 15:04:05 -0700",
}

func parseDate(date string) (time.Time, error) {
	for _, f := range layouts {
		t, err := time.Parse(f, date)

		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("parseDate: could not parse '%v'", date)
}
