package cache

import "context"

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
}
