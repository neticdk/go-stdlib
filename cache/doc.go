// Package cache defines generic interfaces for caching operations.
// It provides standard contracts for key-value caches, allowing different
// implementations (e.g., in-memory, distributed) to be used interchangeably.
//
// # Overview
//
// The core of this package is the `Cache[K comparable, V any]` interface,
// supporting basic Get, Set, Delete, Len, and Clear operations. Keys must
// be comparable, and values can be of any type. All operations accept a
// `context.Context` for cancellation and deadline propagation.
//
// For caches supporting time-based expiration, the `TTLCache[K, V]` interface
// extends `Cache[K, V]` with a `SetWithTTL` method.
//
// # Implementations
//
// This package defines interfaces; concrete implementations are provided
// in sub-packages. See the `inmem` sub-package for a thread-safe,
// configurable in-memory cache implementation supporting TTL, size limits, and
// background garbage collection.
//
// # Errors
//
// The package defines standard errors that implementations can return:
//   - ErrCacheMiss: The requested key was not found.
//   - ErrExpired: The requested key was found but has passed its expiration time.
//   - ErrCacheFull: The cache cannot store a new item because it reached its size limit.
//   - ErrCacheNotStopped: Returned by operations if the cache was not properly stopped (usage may vary by implementation).
package cache
