package parse

import (
	"GoNews/pcg/typeStruct"

	"github.com/mmcdole/gofeed"
)

type Config struct {
	RSSLinks      []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

func ParseRSS(url string) ([]typeStruct.Post, error) {
	var posts []typeStruct.Post

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	for _, item := range feed.Items {

		post := typeStruct.Post{
			Title:   item.Title,
			Content: item.Description,
			PubTime: item.Published,
			Link:    item.Link,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// для тестирования с фикстурой rss
func ParseRSSFixture(fixtureXML string) ([]typeStruct.Post, error) {
	var posts []typeStruct.Post

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(fixtureXML)
	if err != nil {
		return nil, err
	}

	for _, item := range feed.Items {
		post := typeStruct.Post{
			Title:   item.Title,
			Content: item.Description,
			PubTime: item.Published,
			Link:    item.Link,
		}
		posts = append(posts, post)
	}

	return posts, nil
}
