package cache

import (
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func New() *Cache {
	return &Cache{
		items: make(map[string][]byte),
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.items[key]
	return val, ok
}
