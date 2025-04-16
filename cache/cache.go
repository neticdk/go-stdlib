package cache

import (
	"runtime"
	"sync"
	"time"
)

// Cache defines the interface for a key-value cache with time-based expiration.
// Implementations are expected to be safe for concurrent use.
type Cache interface {
	// AddInterval adds a key-value pair to the cache with a specific time-to-live (TTL).
	// If the key already exists and overwrites are disallowed by the implementation's
	// configuration (e.g., allowOverwrites=false), the operation should be ignored
	// if the existing item has not expired.
	// If ttl is negative, it should be treated as 0.
	AddInterval(key string, value any, ttl time.Duration)

	// Add adds a key-value pair to the cache using the implementation's default TTL.
	// It should behave like AddInterval regarding overwrites.
	Add(key string, value any)

	// SetInterval adds or updates a key-value pair in the cache with a specific TTL.
	// Unlike AddInterval, SetInterval should always overwrite the key if it exists,
	// regardless of implementation configuration or expiration status.
	// If ttl is negative, it should be treated as 0.
	SetInterval(key string, value any, ttl time.Duration)

	// Set adds or updates a key-value pair using the implementation's default TTL.
	// It should always overwrite the key if it exists.
	Set(key string, value any)

	// RenewInterval updates the expiration time of an existing key with a new TTL,
	// calculated from the current time.
	// If the key does not exist, the operation should have no effect.
	// If ttl is negative, the operation should have no effect.
	RenewInterval(key string, ttl time.Duration)

	// Renew updates the expiration time of an existing key using the implementation's
	// default TTL as the new TTL, calculated from the current time.
	// If the key does not exist, the operation should have no effect.
	Renew(key string)

	// GetItems retrieves all unexpired items in the cache.
	// It returns a map of string keys to Item values
	// representing the current state of the cache.
	GetItems() map[string]Item

	// GetAllItems retrieves all items in the cache regardless of expiration.
	// It returns a map of string keys to Item values
	// representing the current state of the cache.
	GetAllItems() map[string]Item

	// Get retrieves the value associated with the given key.
	// It returns the value and true if the key exists and has not expired according
	// to the implementation's clock and expiration logic.
	// Otherwise, it returns nil and false.
	Get(key string) (any, bool)

	// GetTTL retrieves the remaining time-to-live (TTL) for the given key.
	// It returns the remaining duration and true if the key exists and has not expired.
	// Otherwise, it returns 0 and false.
	GetTTL(key string) (time.Duration, bool)

	// Delete removes a key-value pair from the cache.
	// If the key does not exist, the operation should have no effect.
	Delete(key string)

	// Clear removes all key-value pairs from the cache.
	Clear()

	// DeleteExpired manually triggers the removal of all expired items from the cache
	// based on the implementation's clock and expiration logic.
	DeleteExpired()

	// Stop stops the garbage collector if it is running.
	// This is typically called when the cache is no longer needed.
	// It is important to call this function to avoid resource leaks
	// and ensure that the goroutine is properly cleaned up.
	Stop()

	// GetGarbageCollector returns the current garbage collector.
	GetGarbageCollector() *garbageCollector

	// SetGarbageCollector sets the garbage collector for the cache.
	SetGarbageCollector(gc *garbageCollector)

	// GetClock returns the clock used by the cache for time-related operations.
	GetClock() Clock

	// SetClock sets the clock for the cache.
	SetClock(clock Clock)
}

// cache represents an in-memory key-value store with time-based expiration.
// It is safe for concurrent use.
type cache struct {
	// items holds the cache data.
	items map[string]Item
	// mu provides concurrency control for safe access to items.
	mu sync.RWMutex
	// garbageCollector manages the periodic removal of expired items.
	garbageCollector *garbageCollector
	// cacheInterval specifies the default duration for items added without a specific TTL.
	cacheInterval time.Duration
	// garbageCollectorInterval defines how often the garbage collector runs.
	garbageCollectorInterval time.Duration
	// allowOverwrites determines if Add operations can overwrite existing, non-expired keys.
	allowOverwrites bool
	// clock provides an abstraction for time, useful for testing.
	clock Clock
}

// Option defines a function signature for configuring a Cache instance.
type Option func(*cache)

// WithCacheInterval sets the default time-to-live (TTL) for items added to the cache
// without a specific interval. The default is 1 second.
func WithCacheInterval(interval time.Duration) Option {
	return func(c *cache) {
		c.cacheInterval = interval
	}
}

// WithGarbageCollectorInterval sets the interval at which the garbage collector
// checks for and removes expired items. If the interval is zero or negative,
// garbage collection is disabled.
func WithGarbageCollectorInterval(interval time.Duration) Option {
	return func(c *cache) {
		c.garbageCollectorInterval = interval
	}
}

// WithAllowOverwrites configures the cache to allow the Add method to overwrite
// an existing key even if it hasn't expired. By default, Add does not overwrite
// existing, non-expired keys.
func WithAllowOverwrites() Option {
	return func(c *cache) {
		c.allowOverwrites = true
	}
}

// WithItems initializes the cache with a pre-defined map of string keys to Item values.
// Note that the expiration times within the provided Items must be set correctly
// relative to the desired start time, as this function does not modify them.
// It replaces any existing items in the cache.
func WithItems(i map[string]Item) Option {
	return func(c *cache) {
		c.items = i
	}
}

// WithMap initializes the cache with a map of string keys to arbitrary values.
// Each value will be wrapped in an Item and assigned an expiration time based
// on the cache's default cacheInterval, calculated from the time this option is applied.
// It replaces any existing items in the cache.
func WithMap(m map[string]any) Option {
	return func(c *cache) {
		items := make(map[string]Item, len(m))
		now := c.clock.Now()
		expTime := now.Add(c.cacheInterval)
		for k, v := range m {
			items[k] = Item{
				Value:      v,
				Expiration: expTime, // Use pre-calculated expiration
			}
		}
		c.items = items
	}
}

// newDefaultCache creates a Cache instance with default settings:
// - cacheInterval: 1 second
// - garbageCollectorInterval: 2 seconds
// - allowOverwrites: false
// - clock: realClock{}
// - items: initialized empty map
func newDefaultCache() *cache {
	c := &cache{
		items:                    make(map[string]Item),
		cacheInterval:            1 * time.Second,
		garbageCollectorInterval: 2 * time.Second,
		allowOverwrites:          false,
		clock:                    &realClock{},
	}
	return c
}

// NewCache creates and initializes a new Cache.
// It accepts optional configuration functions (Options) to customize cache behavior.
// If a positive cacheInterval is set (either by default or via WithCacheInterval),
// it starts a background garbage collector goroutine to periodically remove expired items.
// The garbage collector's run interval can be configured with WithGarbageCollectorInterval.
// A finalizer is set to stop the garbage collector
// when the cache is garbage collected by Go's runtime.
func NewCache(opts ...Option) Cache {
	c := newDefaultCache()

	// Apply user-provided options
	for _, opt := range opts {
		opt(c)
	}

	// Start garbage collector only if items have a default expiration > 0
	if c.cacheInterval > 0 {
		startGarbageCollector(c, c.garbageCollectorInterval)
		runtime.SetFinalizer(c, stopGarbageCollector)
	}
	return c
}

// AddInterval adds a key-value pair to the cache with a specific time-to-live (TTL).
// If the key already exists and allowOverwrites is false (the default),
// the operation is ignored if the existing item has not expired.
// If ttl is negative, it is treated as 0 (item expires immediately or
// is considered expired).
// This method is safe for concurrent use.
func (c *cache) AddInterval(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if item exists, overwrites are disallowed, and item is not expired
	if item, ok := c.items[key]; ok && !c.allowOverwrites && !item.isExpired(c.clock) {
		return // Do not overwrite
	}

	// Ensure TTL is not negative
	ttl = max(0, ttl)
	c.items[key] = Item{
		Value:      value,
		Expiration: c.clock.Now().Add(ttl), // Calculate expiration time
	}
}

// Add adds a key-value pair to the cache using the default cacheInterval as the TTL.
// It behaves like AddInterval regarding overwrites: if the key exists, allowOverwrites
// is false, and the item hasn't expired, the operation is ignored.
// This method is safe for concurrent use.
func (c *cache) Add(key string, value any) {
	c.AddInterval(key, value, c.cacheInterval)
}

// SetInterval adds or updates a key-value pair in the cache with a specific TTL.
// Unlike AddInterval, SetInterval always overwrites the key if it exists, regardless
// of the allowOverwrites setting or expiration status.
// If ttl is negative, it is treated as 0.
// This method is safe for concurrent use.
func (c *cache) SetInterval(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure TTL is not negative
	ttl = max(0, ttl)
	c.items[key] = Item{
		Value:      value,
		Expiration: c.clock.Now().Add(ttl), // Calculate expiration time
	}
}

// Set adds or updates a key-value pair using the default cacheInterval as the TTL.
// It always overwrites the key if it exists.
// This method is safe for concurrent use.
func (c *cache) Set(key string, value any) {
	c.SetInterval(key, value, c.cacheInterval)
}

// RenewInterval updates the expiration time of an existing key with a new TTL,
// calculated from the current time.
// If the key does not exist, the operation has no effect.
// If ttl is negative, the operation has no effect (item retains original expiration).
// This method is safe for concurrent use.
func (c *cache) RenewInterval(key string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return // Key not found
	}

	// Do not allow negative TTL for renewal
	if ttl < 0 {
		return
	}

	// Update expiration based on current time and new TTL
	item.Expiration = c.clock.Now().Add(ttl)
	c.items[key] = item
}

// Renew updates the expiration time of an existing key using the default cacheInterval
// as the new TTL, calculated from the current time.
// If the key does not exist, the operation has no effect.
// This method is safe for concurrent use.
func (c *cache) Renew(key string) {
	c.RenewInterval(key, c.cacheInterval)
}

// GetAllItems retrieves all items in the cache regardless of expiration.
// It returns a map of string keys to Item values representing the current state of the cache.
// This method is safe for concurrent use.
// Note: The returned map is a copy of the internal items map to prevent
// concurrent modification issues. The caller should not modify the returned map.
func (c *cache) GetAllItems() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy of the items map to avoid concurrent modification
	itemsCopy := make(map[string]Item, len(c.items))
	for k, v := range c.items {
		itemsCopy[k] = v
	}
	return itemsCopy
}

// GetItems retrieves all unexpired items in the cache.
// It returns a map of string keys to Item values representing the current state of the cache.
// This method is safe for concurrent use.
// Note: The returned map is a copy of the internal items map to prevent
// concurrent modification issues. The caller should not modify the returned map.
func (c *cache) GetItems() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy of the items map to avoid concurrent modification
	itemsCopy := make(map[string]Item, len(c.items))
	for k, v := range c.items {
		if v.isExpired(c.clock) {
			continue // Skip expired items
		}
		itemsCopy[k] = v
	}
	return itemsCopy
}

// Get retrieves the value associated with the given key.
// It returns the value and true if the key exists and has not expired.
// Otherwise, it returns nil and false.
// This method is safe for concurrent use.
func (c *cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	// Check existence and expiration
	if !exists || item.isExpired(c.clock) {
		return nil, false
	}
	return item.Value, true
}

// GetTTL retrieves the remaining time-to-live (TTL) for the given key.
// It returns the remaining duration and true if the key exists and has not expired.
// Otherwise, it returns 0 and false.
// This method is safe for concurrent use.
func (c *cache) GetTTL(key string) (time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	// Check existence and expiration
	if !exists || item.isExpired(c.clock) {
		return 0, false
	}
	// Calculate remaining duration
	return item.Expiration.Sub(c.clock.Now()), true
}

// Delete removes a key-value pair from the cache.
// If the key does not exist, the operation has no effect.
// This method is safe for concurrent use.
func (c *cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all key-value pairs from the cache, effectively resetting it.
// This method is safe for concurrent use.
func (c *cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Item)
}

// DeleteExpired manually scans the cache and removes all items that have expired
// based on the current time provided by the cache's clock.
// This is typically handled automatically by the background garbage collector if enabled,
// but can be called manually if needed (e.g., if the GC interval is long or GC is disabled).
// Note: This method acquires a write lock for each deletion, which might impact performance
// during the scan if the cache is heavily contended.
func (c *cache) DeleteExpired() {
	c.mu.RLock()
	keysToDelete := []string{}
	now := c.clock.Now()
	for key, item := range c.items {
		if item.Expiration.Before(now) { // Direct comparison instead of isExpired for efficiency here
			keysToDelete = append(keysToDelete, key)
		}
	}
	c.mu.RUnlock()

	// Now delete the collected keys under write lock(s) via Delete method
	if len(keysToDelete) > 0 {
		// Although Delete locks/unlocks individually, we can call it repeatedly.
		// A single lock around the loop might be marginally faster for many deletions,
		// but increases lock contention time. Individual deletes keep lock duration short.
		for _, key := range keysToDelete {
			c.Delete(key) // Delete already handles locking
		}
	}
}

// Stop stops the garbage collector if it is running.
// This is typically called when the cache is no longer needed.
// It is important to call this function to avoid resource leaks
// and ensure that the goroutine is properly cleaned up.
func (c *cache) Stop() {
	stopGarbageCollector(c)
}

// GetGarbageCollector returns the current garbage collector.
// This is typically used for testing or monitoring purposes.
func (c *cache) GetGarbageCollector() *garbageCollector {
	return c.garbageCollector
}

// SetGarbageCollector sets the garbage collector for the cache.
// This is typically used for testing or monitoring purposes.
func (c *cache) SetGarbageCollector(gc *garbageCollector) {
	c.garbageCollector = gc
}

// GetClock returns the clock used by the cache for time-related operations.
func (c *cache) GetClock() Clock {
	return c.clock
}

// SetClock sets the clock for the cache.
// This is typically used for testing purposes to mock time behavior.
func (c *cache) SetClock(clock Clock) {
	c.clock = clock
}
