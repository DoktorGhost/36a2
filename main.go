package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	configPath := "config.json"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	var wg sync.WaitGroup
	postsCh := make(chan []Post, len(config.RSSLinks))

	ticker := time.NewTicker(time.Duration(config.RequestPeriod) * time.Minute)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Regularly fetching RSS feeds...")
			for _, rssLink := range config.RSSLinks {
				posts, err := ParseRSS(rssLink)
				if err != nil {
					fmt.Printf("Error parsing RSS %s: %v\n", rssLink, err)
					continue
				}

				for i, post := range posts {
					fmt.Printf("Post %d:\nTitle: %s\nContent: %s\nPubTime: %s\nLink: %s\n\n",
						i+1, post.Title, post.Content, post.FormatPubTime(), post.Link)
				}
			}
		}
	}

	for _, rssLink := range config.RSSLinks {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			posts, err := ParseRSS(link)
			if err != nil {
				fmt.Printf("Error parsing RSS %s: %v\n", link, err)
				postsCh <- nil
				return
			}

			postsCh <- posts
		}(rssLink)
	}

	wg.Wait()
	close(postsCh)

	for posts := range postsCh {
		if posts != nil {
			for i, post := range posts {
				fmt.Printf("Post %d:\nTitle: %s\nContent: %s\nPubTime: %s\nLink: %s\n\n",
					i+1, post.Title, post.Content, post.FormatPubTime(), post.Link)
			}
		}
	}
}