package main

import (
	"os"
	"testing"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/iamruinous/greed/internal/ui"
)

func TestFeedbinClient(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "greed-test-cache")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	client := feedbin.NewClient(
		os.Getenv("FEEDBIN_USERNAME"),
		os.Getenv("FEEDBIN_PASSWORD"),
		5,
	)
	cache, err := cache.New(tempDir)
	if err != nil {
		t.Fatalf("Error creating cache: %v", err)
	}
	entries, err := ui.FetchEntries(client, cache)
	if err != nil {
		t.Fatalf("Error getting latest feeds: %v", err)
	}
	if len(entries) == 0 {
		t.Error("Expected non-empty feeds slice")
	}
}

func TestCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "greed-test-cache")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cache, err := cache.New(tempDir)
	if err != nil {
		t.Fatalf("Error creating cache: %v", err)
	}
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
