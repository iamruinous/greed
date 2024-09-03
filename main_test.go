package main

import (
	"os"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/iamruinous/greed/internal/ui"
)

func TestFeedbinClient(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is authenticated
		username, password, ok := r.BasicAuth()
		if !ok || username != "testuser" || password != "testpass" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Mock response data
		mockEntries := []feedbin.Entry{
			{
				ID:          1,
				Title:       "Test Entry",
				Author:      "Test Author",
				Summary:     "Test Summary",
				URL:         "https://test.com",
				PublishedAt: time.Now(),
			},
		}

		// Encode and send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockEntries)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "greed-test-cache")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Use the mock server URL instead of the real Feedbin API
	client := feedbin.NewClient("testuser", "testpass", 5)
	client.BaseURL = server.URL

	cache, err := cache.New(tempDir)
	if err != nil {
		t.Fatalf("Error creating cache: %v", err)
	}

	entries, err := ui.FetchEntries(client, cache, false)
	if err != nil {
		t.Fatalf("Error getting latest feeds: %v", err)
	}

	if len(entries) == 0 {
		t.Error("Expected non-empty feeds slice")
	}

	// Additional assertions
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Title != "Test Entry" {
		t.Errorf("Expected entry title 'Test Entry', got '%s'", entries[0].Title)
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
