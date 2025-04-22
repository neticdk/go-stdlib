package inmem

import (
	"context"
	"time"
)

type GarbageCollector[K comparable, V any] interface {
	// Start begins the garbage collection process.
	Start(c garbageCollection[K, V])
	// Stop stops the garbage collection process.
	Stop()
	// IsActive returns true if the garbage collector is currently running.
	IsActive() bool
}

// garbageCollector is responsible for periodically deleting expired items
// from the cache. It runs in a separate goroutine and uses a ticker to
// trigger the deletion process at regular intervals.
type garbageCollector[K comparable, V any] struct {
	// active indicates whether the garbage collector is currently running.
	active bool

	// interval is the duration between each garbage collection cycle.
	interval time.Duration

	// stop is a channel used to signal the garbage collector to stop.
	stop chan bool

	// clock is an interface that provides the current time and
	// a ticker for scheduling garbage collection.
	clock Clock
}

func NewGarbageCollector[K comparable, V any](interval time.Duration) *garbageCollector[K, V] {
	if interval <= 0 {
		return nil
	}
	return &garbageCollector[K, V]{
		interval: interval,
		stop:     make(chan bool),
		active:   false,
		clock:    &realClock{},
	}
}

// start begins the garbage collection process. It runs in a separate
// goroutine and periodically checks for expired items in the cache.
// When a tick happens, it calls the DeleteExpired method on the cache to
// remove expired items.
func (gc *garbageCollector[K, V]) Start(c garbageCollection[K, V]) {
	ticker := gc.clock.NewTicker(gc.interval)

	defer ticker.Stop()
	gc.active = true

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), gc.interval)
			_ = c.deleteExpired(ctx)
			cancel()
		case <-gc.stop:
			// Stop the garbage collector
			gc.active = false
			return
		}
	}
}

// stopGarbageCollector stops the garbage collector by sending a signal
// to the stop channel. This will cause the garbage collector to exit
// its loop and stop running.
// It is important to call this function when the cache is no longer
// needed to avoid resource leaks and ensure that the goroutine is
// properly cleaned up.
func stopGarbageCollector[K comparable, V any](gc GarbageCollector[K, V]) {
	// Prevent panic if GC was never started
	gc.Stop()
}

// stopGarbageCollector stops the garbage collector by sending a signal
// to the stop channel. This will cause the garbage collector to exit
// its loop and stop running.
// It is important to call this function when the cache is no longer
// needed to avoid resource leaks and ensure that the goroutine is
// properly cleaned up.
func (gc *garbageCollector[K, V]) Stop() {
	// Prevent panic if GC was never started
	if gc == nil {
		return
	}
	gc.stop <- true
}

// IsActive returns true if the garbage collector is currently running.
func (gc *garbageCollector[K, V]) IsActive() bool {
	return gc.active
}
