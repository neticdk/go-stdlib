// Package inmem provides a thread-safe, configurable in-memory implementation
// of the cache.Cache and cache.TTLCache interfaces.
//
// # Overview
//
// This package offers `SafeMapCache`, an in-memory cache using a Go map protected
// by a single `sync.RWMutex`. It supports item expiration (TTL), configurable capacity limits,
// and automatic removal of expired items via a background garbage collector.
//
// Note: Due to the single mutex design, this implementation is best suited for
// low-to-moderate concurrency scenarios. High contention on the lock might become a bottleneck.
//
// # Features
//
//   - Thread-Safe: Uses `sync.RWMutex` for safe concurrent operations.
//   - Configurable: Accepts options via `NewSafeMap` for:
//   - Default TTL (`WithDefaultTTL`, default: 0/no expiry).
//   - Maximum number of items (`WithMaxSize`, default: 10000).
//   - Garbage collection interval (`WithGCInterval`, default: 5 minutes, min: 5 seconds).
//   - TTL Support: Items expire after their TTL (set via `SetWithTTL` or the configured default).
//   - Size Limiting: Returns `cache.ErrCacheFull` if `Set` attempts to add a new item when the cache is at `MaxSize`.
//   - Garbage Collection: A background goroutine periodically removes expired items based on the configured interval.
//   - Efficient `Len()`: Provides an O(1) approximate count of active items using an atomic counter.
//   - Explicit Cleanup: Requires calling the `Stop()` method to gracefully shut down the garbage collector and clear the cache.
//   - Context Aware: Operations respect context cancellation/deadlines.
//   - Testable Time: Uses a `Clock` interface for reliable time-based testing.
//
// # Usage
//
// Create a new cache using `NewSafeMap`. It is **essential** to call the `Stop()`
// method when the cache is no longer needed to ensure the background garbage
// collector goroutine is terminated cleanly. Using `defer` is a common pattern.
//
//	// Cache with a default 1-hour TTL, max 5000 items, and 10-minute GC interval
//	cacheOpts := []inmem.SafeMapCacheOption[string, []byte]{
//	    inmem.WithDefaultTTL[string, []byte](1*time.Hour),
//	    inmem.WithMaxSize[string, []byte](5000),
//	    inmem.WithGCInterval[string, []byte](10*time.Minute),
//	}
//	c := inmem.NewSafeMap(cacheOpts...)
//	defer c.Stop(context.Background()) // Ensure Stop is called for cleanup
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	_ = c.Set(ctx, "config:data", someBytes) // Uses default TTL (1 hour)
//	_ = c.SetWithTTL(ctx, "session:abc", tokenBytes, 15*time.Minute) // Explicit TTL
//
//	val, err := c.Get(ctx, "config:data")
//	// ... handle error or use val ...
//
//	count := c.Len() // Get approximate count
//	fmt.Printf("Approximate items in cache: %d\n", count)
//
//	_ = c.Delete(ctx, "session:abc") // Mark session for deletion
//
// # Important Considerations
//
//   - `Delete` Behavior: Calling `Delete(ctx, key)` marks the item as expired immediately
//     (subsequent `Get(ctx, key)` will return `cache.ErrExpired`) and decrements the
//     counter used by `Len()`. However, the key and value **remain in memory**
//     until the next garbage collection cycle physically removes them from the underlying map.
//     Memory is not reclaimed instantly upon calling `Delete`.
//   - `Len()` Approximation: `Len()` returns the count of items considered active (added minus deleted/GC'd).
//     It may include items that have passed their TTL but have not yet been removed by the garbage collector.
//   - Memory Usage: Size is limited by available RAM and the `MaxSize` setting.
//   - Resource Cleanup: **Failure to call the `Stop()` method will result in leaking the garbage collector goroutine.**
package inmem
