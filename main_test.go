package main

import (
	"os"
	"testing"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
)

func TestFeedbinClient(t *testing.T) {
	client := feedbin.NewClient(
		os.Getenv("FEEDBIN_USERNAME"),
		os.Getenv("FEEDBIN_PASSWORD"),
	)
	feeds, err := client.GetLatestFeeds()
	if err != nil {
		t.Fatalf("Error getting latest feeds: %v", err)
	}
	if len(feeds) == 0 {
		t.Error("Expected non-empty feeds slice")
	}
}

func TestCache(t *testing.T) {
	cache := cache.New()
	mockEntries := []feedbin.Entry{
		{ID: 1, Title: "Test Entry", URL: "https://test.com"},
	}

	cache.Set(mockEntries)
	cachedEntries, ok := cache.Get()
	if !ok {
		t.Error("Expected cache hit, got miss")
	}
	if len(cachedEntries) != len(mockEntries) {
		t.Errorf("Expected %d entries, got %d", len(mockEntries), len(cachedEntries))
	}
}
