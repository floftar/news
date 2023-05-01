package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/floftar/news/internal"
)

type Source struct {
	Name string
	Url  string
}

type FetchResult struct {
	Body  *string
	Error error
}

const numberOfItems = 5

var sources = [...]Source{
	{Name: "Yle", Url: "https://feeds.yle.fi/uutiset/v1/majorHeadlines/YLE_UUTISET.rss"},
	{Name: "Hacker News", Url: "https://news.ycombinator.com/rss"},
	{Name: "Lobsters", Url: "https://lobste.rs/rss"},
}

func main() {
	var wg sync.WaitGroup
	wg.Add(len(sources))

	channels := createChannels(len(sources))
	res := make([]FetchResult, len(sources))

	for i, s := range sources {
		go fetch(s.Url, &wg, channels[i])
	}

	for i := range sources {
		res[i] = <-channels[i]
	}

	wg.Wait()

	for i, s := range sources {
		if res[i].Error != nil {
			printError(s, res[i].Error)
		} else {
			printItems(s, res[i].Body)
		}
	}
}

func createChannels(n int) []chan FetchResult {
	channels := make([]chan FetchResult, n)

	for i := range channels {
		channels[i] = make(chan FetchResult)
	}

	return channels
}

func printError(source Source, err error) {
	fmt.Println(source.Name)
	fmt.Println(err)
	fmt.Println()
}

func printItems(source Source, body *string) {
	feed, err := internal.ParseFeed(*body)
	if err != nil {
		printError(source, err)
		return
	}

	fmt.Println(source.Name)

	maxDigits := int(math.Log10(float64(len(feed.Items)))) + 1
	padded := fmt.Sprintf("%%%dd %%s\n", maxDigits)
	indent := strings.Repeat(" ", maxDigits+1)

	for i, item := range feed.Items {
		fmt.Printf(padded, i+1, item.Title)
		fmt.Printf("%s%s - %s\n", indent, item.Link, internal.GetAge(time.Now(), item.Published))

		if i >= numberOfItems-1 {
			break
		}
	}

	fmt.Println()
}

func fetch(url string, wg *sync.WaitGroup, ch chan<- FetchResult) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		ch <- FetchResult{Body: nil, Error: error(err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ch <- FetchResult{Body: nil, Error: errors.New(resp.Status)}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- FetchResult{Body: nil, Error: fmt.Errorf("while reading %s: %v", url, err)}
		return
	}

	s := string(body[:])
	ch <- FetchResult{Body: &s, Error: nil}
}
