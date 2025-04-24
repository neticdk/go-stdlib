package cache

import (
	"context"
	"time"
)

// Cache is an interface that defines the basic operations for a cache.
// It allows for storing, retrieving, and deleting items in the cache.
// The cache is generic and can work with any key-value pair.
// The key type must be comparable, and the value type can be any type.
type Cache[K comparable, V any] interface {
	// Get retrieves a value from the cache by its key.
	// The context can be used to control the operation, such as
	// setting a timeout or cancellation.
	Get(ctx context.Context, key K) (value V, err error)

	// Set stores a value in the cache with the specified key.
	// The context can be used to control the operation, such as
	// setting a timeout or cancellation.
	Set(ctx context.Context, key K, value V) error

	// Delete removes a value from the cache by its key.
	// The context can be used to control the operation, such as
	// setting a timeout or cancellation.
	Delete(ctx context.Context, key K) error

	// Len returns the number of items in the cache.
	Len() int

	// Clear removes all items from the cache.
	// The context can be used to control the operation, such as
	// setting a timeout or cancellation.
	Clear(ctx context.Context) error

	// Stop stops the cache and releases any resources it holds.
	// The context can be used to control the operation, such as
	// setting a timeout or cancellation.
	Stop(ctx context.Context) error
}

// TTLCache is an interface that extends the Cache interface
// to support time-to-live (TTL) functionality.
// It provides a method to set a value with a specific TTL.
type TTLCache[K comparable, V any] interface {
	// Cache defines the basic operations for a cache.
	// It is embedded to provide the basic cache functionality.
	Cache[K, V]

	// SetWithTTL stores a value in the cache with the specified key
	// and a time-to-live (TTL) duration.
	// The context can be used to control the operation, such as
	// setting a timeout or cancellation.
	SetWithTTL(ctx context.Context, key K, value V, ttl time.Duration) error
}
