package transliterate

import (
	"sync"
	"sync/atomic"
)

var defaultCache = newTranslitCache(1000)

// translitCache implements a concurrent-safe cache for character
// transliterations
type translitCache struct {
	sync.RWMutex
	entries map[rune]string
	maxSize int
	hits    uint64 // for monitoring cache effectiveness
}

// newTranslitCache creates a new cache with specified maximum size
func newTranslitCache(size int) *translitCache {
	return &translitCache{
		entries: make(map[rune]string),
		maxSize: size,
	}
}

// get retrieves a cached transliteration if available
func (c *translitCache) get(r rune) (string, bool) {
	c.RLock()
	defer c.RUnlock()

	val, ok := c.entries[r]
	if ok {
		atomic.AddUint64(&c.hits, 1)
	}
	return val, ok
}

// set adds or updates a cache entry
func (c *translitCache) set(r rune, s string) {
	c.Lock()
	defer c.Unlock()

	// Simple eviction strategy: clear cache if full
	if len(c.entries) >= c.maxSize {
		c.entries = make(map[rune]string)
	}

	c.entries[r] = s
}

// GetCacheStats returns the number of cache hits since last reset.
// This can be used to monitor cache effectiveness.
func GetCacheStats() (hits uint64) {
	return atomic.LoadUint64(&defaultCache.hits)
}

// ResetCacheStats zeros out the cache statistics counter.
// Useful for beginning a new monitoring period.
func ResetCacheStats() {
	atomic.StoreUint64(&defaultCache.hits, 0)
}

// ClearCache empties the transliteration cache.
// This can be useful when memory pressure is high or when
// preparing for a new batch of translations.
func ClearCache() {
	defaultCache.Lock()
	defer defaultCache.Unlock()
	defaultCache.entries = make(map[rune]string)
}

// GetCacheSize returns the current size of the transliteration cache.
// This can be useful for monitoring memory usage.
func GetCacheSize() int {
	defaultCache.RLock()
	defer defaultCache.RUnlock()
	return len(defaultCache.entries)
}
