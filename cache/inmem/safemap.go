package inmem

import (
	"context"
	"maps"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/neticdk/go-stdlib/cache"
	"github.com/neticdk/go-stdlib/xslices"
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

// safeMapCache is a thread-safe in-memory cache implementation.
// It uses a map to store items and a mutex to ensure thread safety.
type safeMapCache[K comparable, V any] struct {
	// items is the map that stores the cached items.
	items map[K]item[V]

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
func NewSafeMap[K comparable, V any]() *safeMapCache[K, V] {
	// Create the cache instance
	c := &safeMapCache[K, V]{
		items: make(map[K]item[V]),
		clock: &realClock{},
	}

	// Initialize the garbage collector
	gc := NewGarbageCollector[K, V](5 * time.Minute)
	if gc != nil {
		c.garbageCollector = gc
		go c.garbageCollector.Start(c)
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
	defer c.mu.RUnlock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the read lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		var zeroV V
		return zeroV, err
	}

	// Retrieve the item from the cache
	storedItem, ok := c.items[key]
	if !ok {
		var zeroV V
		return zeroV, cache.NewErrNotFound()
	}

	// Check if the item has expired
	if !storedItem.expiresAt.IsZero() && c.clock.Now().After(storedItem.expiresAt) {
		var zeroV V
		return zeroV, cache.NewErrExpired()
	}

	// Return the value of the item
	return storedItem.value, nil
}

// Set stores an item in the cache with no expiration.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
func (c *safeMapCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return c.SetWithTTL(ctx, key, value, 0)
}

// SetWithTTL stores an item in the cache with a specified expiration time.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
// It also checks if the TTL is valid (greater than 0).
func (c *safeMapCache[K, V]) SetWithTTL(ctx context.Context, key K, value V, ttl time.Duration) error {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		return err
	}

	// Use a write lock to ensure thread safety
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		return err
	}

	expiry := time.Time{} // Default: No expiry

	// If the TTL is greater than 0, set the expiry time
	if ttl > 0 {
		expiry = c.clock.Now().Add(ttl)
	}

	// Set the item in the cache
	c.items[key] = item[V]{value: value, expiresAt: expiry}
	return nil
}

// Delete removes an item from the cache.
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
	defer c.mu.Unlock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Delete the item from the cache
	delete(c.items, key)
	return nil
}

// Len returns the number of unexpired items in the cache.
// It uses a read lock to ensure thread safety.
func (c *safeMapCache[K, V]) Len() int {
	// Use a read lock to ensure thread safety
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Find all non-expired items
	nonExpired := xslices.Filter(slices.Collect(maps.Values(c.items)), func(v item[V]) bool {
		return v.expiresAt.IsZero() || !c.clock.Now().After(v.expiresAt)
	})

	// Return the count of non-expired items
	return len(nonExpired)
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
	defer c.mu.Unlock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Clear the items map
	c.items = make(map[K]item[V])

	return nil
}

// deleteExpired deletes expired items from the cache.
// It uses a write lock to ensure thread safety.
// It checks if the context is canceled or timed out.
// If the context is canceled or times out, it returns an error.
// It iterates over the items map and deletes items that have expired.
func (c *safeMapCache[K, V]) deleteExpired(ctx context.Context) error {
	// Check if the context is canceled or timed out
	if err := ctx.Err(); err != nil {
		return err
	}

	// Use a write lock to ensure thread safety
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the context is canceled or timed out again
	// This is important because the context can be canceled
	// while the write lock is held, and we need to ensure
	// that we respect the cancellation.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Iterate over the items map and delete expired items
	for key, storedItem := range c.items {
		if !storedItem.expiresAt.IsZero() && c.clock.Now().After(storedItem.expiresAt) {
			delete(c.items, key)
		}
	}

	return nil
}
