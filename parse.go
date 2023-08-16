package main

import (
	"github.com/mmcdole/gofeed"
)

type Config struct {
	RSSLinks      []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

func ParseRSS(url string) ([]Post, error) {
	var posts []Post

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	for i, item := range feed.Items {
		if i >= 1 {
			break
		}

		post := Post{
			Title:   item.Title,
			Content: item.Description,
			PubTime: item.PublishedParsed.Unix(),
			Link:    item.Link,
		}
		posts = append(posts, post)
	}

	return posts, nil
}
