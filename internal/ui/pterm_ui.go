package ui

import (
	"math/rand"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/pterm/pterm"
)

func FetchEntries(client *feedbin.Client, cache *cache.Cache) ([]feedbin.Entry, error) {
	if cachedEntries, ok := cache.Get(); ok {
		return cachedEntries, nil
	}

	entries, err := client.GetLatestEntries()
	if err != nil {
		return nil, err
	}

	cache.Set(entries)
	return entries, nil
}

func DisplayEntries(entries []feedbin.Entry, displayLimit int, randomEntries bool) {
	pterm.DefaultHeader.WithFullWidth().Println("Feedbin Latest Feeds")

	if randomEntries {
		rand.Shuffle(len(entries), func(i, j int) {
			entries[i], entries[j] = entries[j], entries[i]
		})
	}

	displayCount := min(displayLimit, len(entries))
	for _, entry := range entries[:displayCount] {
		pterm.DefaultSection.Println(entry.Title)
		pterm.Info.Printf("Author: %s\n", entry.Author)
		pterm.Info.Printf("Published: %s\n", entry.PublishedAt.Format("2006-01-02 15:04:05"))
		pterm.Info.Printf("URL: %s\n", entry.URL)
		pterm.Println(entry.Summary)
		pterm.Println()
	}

	pterm.DefaultBasicText.Println("Press Ctrl+C to quit")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
