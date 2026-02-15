package main

import (
	"sync"
	"time"
)

// item holds the actual data and its "death certificate" timestamp
type item struct {
	value      any
	expiryTime int64 // Unix nanoseconds
}

type Cache struct {
	data sync.Map
}

// NewCache initializes the cache and starts the ONE and ONLY janitor
func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{}
	go c.janitor(cleanupInterval)
	return c
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	// Calculate the exact moment this key should die
	expiry := time.Now().Add(ttl).UnixNano()
	c.data.Store(key, item{
		value:      value,
		expiryTime: expiry,
	})
}

func (c *Cache) Get(key string) (any, bool) {
	val, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	it := val.(item)
	// Check if the item has already expired but hasn't been cleaned yet
	if time.Now().UnixNano() > it.expiryTime {
		return nil, false
	}

	return it.value, true
}

// janitor is a single background loop that cleans up expired items
func (c *Cache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		now := time.Now().UnixNano()

		// Range over all keys in the sync.Map
		c.data.Range(func(key, value any) bool {
			it := value.(item)
			if now > it.expiryTime {
				c.data.Delete(key)
			}
			return true // continue iteration
		})
	}
}

func main() {
	// Initialize with a janitor that sweeps every 2 seconds
	cache := NewCache(2 * time.Second)

	cache.Set("user_1", "Alice", 5*time.Second)

	// Immediate access
	if val, ok := cache.Get("user_1"); ok {
		println("Found:", val.(string))
	}

	// Wait for TTL and Janitor sweep
	time.Sleep(6 * time.Second)

	if _, ok := cache.Get("user_1"); !ok {
		println("Key cleaned up by Janitor")
	}
}
