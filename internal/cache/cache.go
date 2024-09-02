package cache

import (
	"sync"
	"time"

	"github.com/iamruinous/greed/internal/feedbin"
)

type Cache struct {
	entries []feedbin.Entry
	mu      sync.RWMutex
	exp     time.Time
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(entries []feedbin.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = entries
	c.exp = time.Now().Add(15 * time.Minute)
}

func (c *Cache) Get() ([]feedbin.Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if time.Now().Before(c.exp) {
		return c.entries, true
	}
	return nil, false
}
