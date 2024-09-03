package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/iamruinous/greed/internal/feedbin"
)

type Cache struct {
	entries  []feedbin.Entry
	mu       sync.RWMutex
	exp      time.Time
	filePath string
}

func New(cacheDir string) (*Cache, error) {
	if cacheDir == "" {
		return nil, fmt.Errorf("cache directory is not set")
	}

	cacheFile := filepath.Join(cacheDir, "feedbin_cache.json")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{
		filePath: cacheFile,
	}, nil
}

func (c *Cache) Set(entries []feedbin.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = entries
	c.exp = time.Now().Add(5 * time.Minute)

	data, err := json.Marshal(struct {
		Entries []feedbin.Entry `json:"entries"`
		Exp     time.Time       `json:"exp"`
	}{
		Entries: c.entries,
		Exp:     c.exp,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := os.WriteFile(c.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func (c *Cache) Get() ([]feedbin.Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return nil, false
	}

	var cachedData struct {
		Entries []feedbin.Entry `json:"entries"`
		Exp     time.Time       `json:"exp"`
	}

	if err := json.Unmarshal(data, &cachedData); err != nil {
		return nil, false
	}

	if time.Now().Before(cachedData.Exp) {
		c.entries = cachedData.Entries
		c.exp = cachedData.Exp
		return c.entries, true
	}

	return nil, false
}
