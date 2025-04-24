package inmem

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/neticdk/go-stdlib/cache"
)

const (
	maxSize           = 10000
	defaultTTL        = 0
	defaultGCInterval = 5 * time.Minute
)

// item represents a single item in the cache.
// It contains the value and the expiration time.
// The expiration time is set to zero if the item does not expire.
type item[V any] struct {
	// value is the actual value stored in the cache.
	value V

	// expiresAt is the time when the item expires.
	expiresAt time.Time
}

func (i *item[V]) expire() {
	// Set the expiration time to the one value
	i.expiresAt = time.Unix(0, 1)
}

type SafeMapCacheOption[K comparable, V any] func(*safeMapCache[K, V])

// WithDefaultTTL sets the default time-to-live (TTL) for items in the cache.
// If the TTL is set to 0, items will not expire.
// The default TTL is 0, meaning items will not expire by default.
func WithDefaultTTL[K comparable, V any](ttl time.Duration) SafeMapCacheOption[K, V] {
	return func(c *safeMapCache[K, V]) {
		c.defaultTTL = ttl
	}
}

// WithMaxSize sets the maximum size of the cache.
// If the cache exceeds this size, it will not accept new items.
// The default size is 10000.
func WithMaxSize[K comparable, V any](size int) SafeMapCacheOption[K, V] {
	return func(c *safeMapCache[K, V]) {
		c.maxSize = size
	}
}

// WithGCInterval sets the interval for the garbage collector.
// The garbage collector will run at this interval to delete expired items.
// The default interval is 5 minutes.
// It cannot be lower than 5 seconds.
func WithGCInterval[K comparable, V any](interval time.Duration) SafeMapCacheOption[K, V] {
	return func(c *safeMapCache[K, V]) {
		if interval < 5*time.Second {
			interval = 5 * time.Second
		}
		c.defaultGCInterval = interval
	}
}

type SafeMapCache[K comparable, V any] struct {
	*safeMapCache[K, V]
}

// safeMapCache is a thread-safe in-memory cache implementation.
// It uses a map to store items and a mutex to ensure thread safety.
// Since there is only one mutex for the entire cache, it is not a good idea to use this
// implementation for high concurrency scenarios.
// It is designed for low to moderate concurrency.
type safeMapCache[K comparable, V any] struct {
	// items is the map that stores the cached items.
	items map[K]*item[V]

	// count is an atomic counter that keeps track of the number of items in the cache.
	// It is used for the Len() method to provide an approximate count.
	count atomic.Int64

	// defaultTTL is the default time-to-live (TTL) for items in the cache.
	// It is used when no TTL is specified for an item.
	// If set to 0, items will not expire.
	// By default, it is set to 0.
	defaultTTL time.Duration

	// defaultGCInterval is the default interval for the garbage collector.
	// It is used to determine how often expired items are deleted from the cache.
	// By default, it is set to 5 minutes.
	defaultGCInterval time.Duration

	// maxSize is the maximum number of items allowed in the cache.
	// If the cache exceeds this size, it will not accept new items.
	// By default, it is set to 10000.
	maxSize int

	// garbageCollector is the garbage collector that periodically
	// deletes expired items from the cache.
	garbageCollector GarbageCollector[K, V]

	// clock is an interface that provides the current time and
	// a ticker for scheduling garbage collection.
	// It allows for easy mocking and testing of time-dependent code.
	clock Clock

	// mu is a mutex that ensures thread safety for the cache operations.
	mu sync.RWMutex
}

// NewSafeMap creates a new instance of safeMapCache.
// It initializes the items map and sets up the garbage collector.
// The garbage collector is started in a separate goroutine.
// The clock is set to the real time package by default.
func NewSafeMap[K comparable, V any](opts ...SafeMapCacheOption[K, V]) *safeMapCache[K, V] {
	// Create the cache instance
	c := &safeMapCache[K, V]{
		items:             make(map[K]*item[V]),
		clock:             &realClock{},
		defaultTTL:        defaultTTL,
		defaultGCInterval: defaultGCInterval,
		maxSize:           maxSize,
	}

	// Apply options to the cache
	for _, opt := range opts {
		opt(c)
	}

	// Initialize the garbage collector
	gc := NewGarbageCollector[K, V](c.defaultGCInterval)
	if gc != nil {
		c.garbageCollector = gc
		go c.garbageCollector.Start(context.Background(), c)
	}

	// Add cleanup to stop the garbage collector when the cache is no longer needed
	runtime.AddCleanup(c, stopGarbageCollector[K, V], c.garbageCollector)

	return c
}

// Get retrieves an item from the cache.
// It uses a read lock to ensure thread safety.
// It checks if the item exists and if it has expired.
// If the item is found and not expired, it returns the value.
// If the item is not found or expired, it returns an error.
// The context is used to check for cancellation or timeout.
// If the context is canceled or times out, it returns an error.
func (c *safeMapCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		var zeroV V
		return zeroV, err
	}

	// Use a read lock to ensure thread safety
	c.mu.RLock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the read lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		var zeroV V
		c.mu.RUnlock() // Release Read Lock
		return zeroV, err
	}

	// Retrieve the item from the cache
	storedItem, ok := c.items[key]
	if !ok {
		var zeroV V
		c.mu.RUnlock() // Release Read Lock
		return zeroV, cache.NewErrCacheMiss()
	}

	// Check if the item has expired
	if !storedItem.expiresAt.IsZero() && c.clock.Now().After(storedItem.expiresAt) {
		var zeroV V
		c.mu.RUnlock() // Release Read Lock
		return zeroV, cache.NewErrExpired()
	}

	// Return the value of the item
	c.mu.RUnlock() // Release Read Lock
	return storedItem.value, nil
}

// Set stores an item in the cache with no expiration.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
func (c *safeMapCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return c.setInternal(ctx, key, value, c.defaultTTL)
}

// SetWithTTL stores an item in the cache with a specified expiration time.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
// It also checks if the TTL is valid (greater than 0).
func (c *safeMapCache[K, V]) SetWithTTL(ctx context.Context, key K, value V, ttl time.Duration) error {
	return c.setInternal(ctx, key, value, ttl)
}

// setInternal is a helper method that stores an item in the cache.
// It is used by both Set and SetWithTTL.
func (c *safeMapCache[K, V]) setInternal(ctx context.Context, key K, value V, ttl time.Duration) error {
	// Check if the context is canceled or timed out
	// before acquiring the lock
	if err := ctx.Err(); err != nil {
		return err
	}

	// Use a write lock to ensure thread safety
	c.mu.Lock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		c.mu.Unlock()
		return err
	}

	// Determine the expiration time based on the TTL
	now := c.clock.Now()
	expiry := now.Add(ttl)
	if ttl == 0 {
		expiry = time.Time{} // No expiration
	}

	// Determine existing and new item states.
	// This is used to calculate the delta for the count
	// which is used for Len().
	var delta int64
	existingItem, keyExists := c.items[key]
	wasActive := keyExists && (existingItem.expiresAt.IsZero() || now.Before(existingItem.expiresAt))
	willBeActive := expiry.IsZero() || now.Before(expiry)

	// Determine the delta based on the existing and new item states.
	if !keyExists && willBeActive {
		delta = 1 // Adding a new active item
	} else if keyExists && !wasActive && willBeActive {
		delta = 1 // Replacing inactive with active
	} else if keyExists && wasActive && !willBeActive {
		delta = -1 // Replacing active with inactive (or expired)
	} // Else no change in count (updating existing, non-expired item)

	if delta > 0 && c.count.Load() >= int64(c.maxSize) {
		c.mu.Unlock() // Release Write Lock
		return cache.NewErrCacheFull()
	}

	// Store the item in the cache
	c.items[key] = &item[V]{value: value, expiresAt: expiry}

	c.mu.Unlock() // Release Write Lock

	// Update the count atomically
	if delta != 0 {
		c.count.Add(delta)
	}

	return nil
}

// Delete sets an item from the cache to expire.
// It does not actually delete the item from the cache.
// The actual deletion is handled by the garbage collector.
// If the item does not exist, it does nothing.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
// It also checks if the item exists before deleting it.
func (c *safeMapCache[K, V]) Delete(ctx context.Context, key K) error {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		return err
	}

	// Use a write lock to ensure thread safety
	c.mu.Lock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		c.mu.Unlock() // Release Write Lock
		return err
	}

	// Expire the item from the cache
	if _, ok := c.items[key]; ok {
		c.items[key].expire()
		c.count.Add(-1) // Decrement the count atomically
	}

	c.mu.Unlock() // Release Write Lock
	return nil
}

// Len returns an approximate of the number of items in the cache.
// It relies on a counter that is updated when items are added or deleted.
// This is not a precise count, but it is efficient and fast.
//
// Example:
//
//	cache := NewSafeMap[string, any]()
//	// GC every minute
//	cache.SetWithTTL(context.Background(), "key", "value", 10*time.Second)
//	fmt.Print(cache.Len()) // 1
//	// Wait for 20 seconds
//	time.Sleep(20 * time.Second)
//	fmt.Print(cache.Len()) // 1 (item is expired but still counted)
//	// Wait for GC
//	time.Sleep(60 * time.Second)
//	fmt.Print(cache.Len()) // 0
func (c *safeMapCache[K, V]) Len() int {
	// Use a read lock to ensure thread safety
	c.mu.RLock()

	// Get the current count atomically
	l := int(c.count.Load())
	c.mu.RUnlock() // Release Read Lock

	return l // Return the count
}

// Clear removes all items from the cache.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
func (c *safeMapCache[K, V]) Clear(ctx context.Context) error {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		return err
	}

	// Use a write lock to ensure thread safety
	c.mu.Lock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		c.mu.Unlock() // Release Write Lock
		return err
	}

	// Clear the items map
	c.items = make(map[K]*item[V])
	c.count.Store(0) // Reset the count to 0

	c.mu.Unlock() // Release Write Lock
	return nil
}

// deleteExpired deletes expired items from the cache.
// It updates the count atomically.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
// It iterates over the items map and deletes items that have expired.
func (c *safeMapCache[K, V]) deleteExpired(ctx context.Context) error {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		return err
	}

	// Get the current time
	now := c.clock.Now()

	// Acquire a read lock to check the current items
	c.mu.RLock()

	// Find expired items by key
	expiredKeys := make([]K, 0)
	for key, storedItem := range c.items {
		if !storedItem.expiresAt.IsZero() && now.After(storedItem.expiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Release the read lock
	c.mu.RUnlock()

	// Use a write lock to ensure thread safety
	c.mu.Lock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		c.mu.Unlock() // Release Write Lock
		return err
	}

	// Iterate over the items map and delete expired items
	var delta int64
	for _, key := range expiredKeys {
		delete(c.items, key)
		delta--
	}

	// Update the count atomically
	c.count.Add(delta)

	c.mu.Unlock() // Release Write Lock
	return nil
}

// Stop stops the garbage collector and clears the cache.
// This method should be called when the cache is no longer needed.
func (c *safeMapCache[K, V]) Stop(ctx context.Context) error {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		return err
	}

	// Use a write lock to ensure thread safety
	c.mu.Lock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		c.mu.Unlock() // Release Write Lock
		return err
	}

	if err := c.garbageCollector.Stop(ctx); err != nil {
		c.mu.Unlock() // Release Write Lock
		return err
	}

	c.mu.Unlock()
	c.Clear(ctx) // Clear the cache

	return nil
}
