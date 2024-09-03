package ui

import (
	"math/rand"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/pterm/pterm"
)

func FetchEntries(client *feedbin.Client, cache *cache.Cache, showProgress bool) ([]feedbin.Entry, error) {
	var spinner *pterm.SpinnerPrinter
	if showProgress {
		spinner, _ = pterm.DefaultSpinner.Start("Fetching entries...")
	}

	if cachedEntries, ok := cache.Get(); ok {
		if showProgress {
			spinner.Success("Entries fetched from cache")
		}
		return cachedEntries, nil
	}

	entries, err := client.GetLatestEntries()
	if err != nil {
		if showProgress {
			spinner.Fail("Failed to fetch entries")
		}
		return nil, err
	}

	cache.Set(entries)
	if showProgress {
		spinner.Success("Entries fetched successfully")
	}
	return entries, nil
}

func DisplayEntries(entries []feedbin.Entry, displayLimit int, randomEntries bool) {
	if randomEntries {
		rand.Shuffle(len(entries), func(i, j int) {
			entries[i], entries[j] = entries[j], entries[i]
		})
	}

	pterm.DefaultHeader.WithFullWidth().Println("Feedbin Latest Feeds")

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
