package ui

import (
	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/pterm/pterm"
)

func FetchEntries(client *feedbin.Client, cache *cache.Cache) ([]feedbin.Entry, error) {
	if cachedEntries, ok := cache.Get(); ok {
		return cachedEntries, nil
	}

	entries, err := client.GetLatestFeeds()
	if err != nil {
		return nil, err
	}

	cache.Set(entries)
	return entries, nil
}

func DisplayEntries(entries []feedbin.Entry, displayLimit int) {
	pterm.DefaultHeader.WithFullWidth().Println("Feedbin Latest Feeds")

	for _, entry := range entries[:displayLimit] {
		pterm.DefaultSection.Println(entry.Title)
		pterm.Info.Printf("Author: %s\n", entry.Author)
		pterm.Info.Printf("Published: %s\n", entry.PublishedAt.Format("2006-01-02 15:04:05"))
		pterm.Info.Printf("URL: %s\n", entry.URL)
		pterm.Println(entry.Summary)
		pterm.Println()
	}

	pterm.DefaultBasicText.Println("Press Ctrl+C to quit")
}
