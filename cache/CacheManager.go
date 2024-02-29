package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	timeToLive int
	entries    map[string]cacheEntry
	mu         sync.Mutex
}

func NewCache(ttl int) *Cache {
	cache := Cache{
		timeToLive: ttl,
		entries:    make(map[string]cacheEntry),
		mu:         sync.Mutex{},
	}
	go cache.clearing(time.Tick(time.Second * time.Duration(ttl)))
	return &cache
}

func (cache *Cache) clearing(channel <-chan time.Time) {
	ttl := float64(cache.timeToLive)
	for _ = range channel {
		for key, entry := range cache.entries {
			seconds := time.Now().Sub(entry.createdAt).Seconds()
			if seconds >= ttl {
				cache.delete(key)
			}
		}
	}
}

func (cache *Cache) delete(key string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	delete(cache.entries, key)
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.entries[key]
	return entry.val, ok

}
