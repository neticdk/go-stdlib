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

// Must implement this to allow garbage collection
type garbageCollection[K comparable, V any] interface {
	deleteExpired(ctx context.Context) error
}

type item[V any] struct {
	value     V
	expiresAt time.Time
}

type safeMapCache[K comparable, V any] struct {
	items map[K]item[V]

	garbageCollector GarbageCollector[K, V]
	clock            Clock
	mu               sync.RWMutex
}

func NewSafeMapCache[K comparable, V any]() *safeMapCache[K, V] {
	c := &safeMapCache[K, V]{
		items: make(map[K]item[V]),
		clock: &realClock{},
	}
	gc := NewGarbageCollector[K, V](5 * time.Minute)
	if gc != nil {
		c.garbageCollector = gc
		go c.garbageCollector.Start(c)
	}
	runtime.AddCleanup(c, stopGarbageCollector[K, V], c.garbageCollector)
	return c
}

func (c *safeMapCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	if err := ctx.Err(); err != nil {
		var zeroV V
		return zeroV, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		var zeroV V
		return zeroV, err
	}

	storedItem, ok := c.items[key]
	if !ok {
		var zeroV V
		return zeroV, cache.NewErrNotFound()
	}
	if !storedItem.expiresAt.IsZero() && c.clock.Now().After(storedItem.expiresAt) {
		var zeroV V
		return zeroV, cache.NewErrExpired()
	}
	return storedItem.value, nil
}

func (c *safeMapCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return c.SetWithTTL(ctx, key, value, 0)
}

func (c *safeMapCache[K, V]) SetWithTTL(ctx context.Context, key K, value V, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	expiry := time.Time{} // Default: No expiry
	// you'd need a way to pass TTL here, e.g., via options or if V is a struct containing TTL
	if ttl > 0 {
		expiry = c.clock.Now().Add(ttl)
	}
	c.items[key] = item[V]{value: value, expiresAt: expiry}

	return nil
}

func (c *safeMapCache[K, V]) Delete(ctx context.Context, key K) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	delete(c.items, key)
	return nil
}

func (c *safeMapCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	nonExpired := xslices.Filter(slices.Collect(maps.Values(c.items)), func(v item[V]) bool {
		return v.expiresAt.IsZero() || !c.clock.Now().After(v.expiresAt)
	})

	return len(nonExpired)
}

func (c *safeMapCache[K, V]) Clear(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	c.items = make(map[K]item[V])

	return nil
}

func (c *safeMapCache[K, V]) deleteExpired(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	for key, storedItem := range c.items {
		if !storedItem.expiresAt.IsZero() && c.clock.Now().After(storedItem.expiresAt) {
			delete(c.items, key)
		}
	}

	return nil
}
