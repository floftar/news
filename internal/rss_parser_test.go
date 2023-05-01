package internal

import (
	"testing"
	"time"
)

// See: https://en.wikipedia.org/wiki/RSS
func TestWikiSample(t *testing.T) {
	feed := parseFeed(t, `
	<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
	  <channel>
		<title>RSS Title</title>
		<item>
		  <title>Example entry</title>
		  <link>http://www.example.com/blog/post/1</link>
		  <pubDate>Sun, 6 Sep 2009 16:20:00 +0000</pubDate>
		</item>
	  </channel>
	</rss>
	`)

	equal(t, feed.Title, "RSS Title")
	equal(t, 1, len(feed.Items))

	item := feed.Items[0]
	equal(t, item.Title, "Example entry")
	equal(t, item.Link, "http://www.example.com/blog/post/1")
	equalTime(t, item.Published, time.Date(2009, 9, 6, 16, 20, 0, 0, time.UTC))
}

func TestEncoding(t *testing.T) {
	feed := parseFeed(t, `
	<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
	  <channel>
		<title>RSS Title</title>
		<item>
		  <title>Moody&#x27;s downgrades US banking 🏦</title>
		  <link>http://www.example.com/blog/post/1</link>
		  <pubDate>Sun, 6 Sep 2009 16:20:00 +0000</pubDate>
		</item>
	  </channel>
	</rss>
	`)

	equal(t, feed.Items[0].Title, "Moody's downgrades US banking 🏦")

}

func TestMultipeAndSortByPublished(t *testing.T) {
	feed := parseFeed(t, `
	<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
	  <channel>
		<title>RSS Title</title>
		<item>
		  <title>Old</title>
		  <link>http://www.example.com</link>
		  <pubDate>Thu, 7 Sep 2000 16:20:00 +0000</pubDate>
		</item>
		<item>
		  <title>New</title>
		  <link>http://www.example.com?id=1</link>
		  <pubDate>Sun, 6 Sep 2009 18:21:00 +0000</pubDate>
		</item>
	  </channel>
	</rss>
	`)

	newItem := feed.Items[0]
	equal(t, newItem.Title, "New")
	equal(t, newItem.Link, "http://www.example.com?id=1")
	equalTime(t, newItem.Published, time.Date(2009, 9, 6, 18, 21, 0, 0, time.UTC))

	oldItem := feed.Items[1]
	equal(t, oldItem.Title, "Old")
	equal(t, oldItem.Link, "http://www.example.com")
	equalTime(t, oldItem.Published, time.Date(2000, 9, 7, 16, 20, 0, 0, time.UTC))
}

func TestDateParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{
			"Sun, 6 Sep 2009 16:20:00 +0000",
			time.Date(2009, 9, 6, 16, 20, 0, 0, time.UTC),
		},
		{
			"Sat, 28 Jan 2023 12:42:57 -0600",
			time.Date(2023, 1, 28, 12, 42, 57, 0, toTz("America/Chicago")),
		},
		{
			"Sat, 28 Jan 2023 16:38:53 +0200",
			time.Date(2023, 1, 28, 16, 38, 53, 0, toTz("Europe/Helsinki")),
		},
	}

	for _, test := range tests {
		if actual, err := parseDate(test.input); !actual.Equal(test.expected) {
			t.Errorf("parseDate(%q) = %v, error = %v", test.input, actual, err)
		}

	}
}

func parseFeed(t *testing.T, body string) Feed {
	feed, err := ParseFeed(body)

	if err != nil {
		t.Fatalf("%v", err)
	}

	return feed
}

func equal[T string | int](t *testing.T, actual T, expected T) {
	if expected != actual {
		t.Errorf("Expected %v; actual %v", expected, actual)
	}
}

func equalTime(t *testing.T, actual time.Time, expected time.Time) {
	if !expected.Equal(actual) {
		t.Errorf("Expected %v; actual %v", expected, actual)
	}
}

func toTz(loc string) *time.Location {
	tz, err := time.LoadLocation(loc)
	if err != nil {
		panic(err)
	}

	return tz
}
