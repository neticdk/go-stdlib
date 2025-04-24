package inmem

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Must implement this to allow garbage collection
type garbageCollection[K comparable, V any] interface {
	deleteExpired(ctx context.Context) error
}

// GarbageCollector is an interface that defines the methods for a garbage collector.
// It is responsible for periodically deleting expired items from the cache.
type GarbageCollector[K comparable, V any] interface {
	// Start begins the garbage collection process.
	// This should be called in a separate goroutine.
	// E.g. go gc.Start(c)
	Start(ctx context.Context, c garbageCollection[K, V]) error

	// Stop stops the garbage collection process.
	Stop(ctx context.Context) error

	// IsActive returns true if the garbage collector is currently running.
	IsActive() bool
}

// garbageCollector is responsible for periodically deleting expired items
// from the cache. It should be run in a separate goroutine.
// It uses a ticker to trigger the deletion process at regular intervals.
type garbageCollector[K comparable, V any] struct {
	// active indicates whether the garbage collector is currently running.
	active bool

	// interval is the duration between each garbage collection cycle.
	interval time.Duration

	// stop is a channel used to signal the garbage collector to stop.
	stop chan bool

	// stopped is a channel used to signal that the garbage collector has stopped.
	stopped chan bool

	// clock is an interface that provides the current time and
	// a ticker for scheduling garbage collection.
	clock Clock

	// mu is a mutex used to synchronize access to the active field.
	mu sync.RWMutex

	// logger is used for logging messages.
	// Only used to log garbage collection errors.
	logger *slog.Logger
}

// NewGarbageCollector creates a new garbage collector with the specified
// interval. The interval must be greater than 0. If not, nil is returned.
func NewGarbageCollector[K comparable, V any](interval time.Duration) *garbageCollector[K, V] {
	if interval <= 0 {
		return nil
	}
	return &garbageCollector[K, V]{
		interval: interval,
		stop:     make(chan bool, 1),
		stopped:  make(chan bool, 1),
		active:   false,
		clock:    &realClock{},
		mu:       sync.RWMutex{},
		logger:   slog.Default(),
	}
}

// Start begins the garbage collection process.
// It should be run in a separate goroutine.
// It periodically checks for expired items in the cache.
// When a tick happens, it calls the DeleteExpired method on the cache to
// remove expired items.
func (gc *garbageCollector[K, V]) Start(ctx context.Context, c garbageCollection[K, V]) error {
	ticker := gc.clock.NewTicker(gc.interval)

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	defer ticker.Stop()
	gc.mu.Lock()
	gc.active = true
	gc.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	for {
		select {
		case <-ticker.C:
			// Perform garbage collection
			ctx, cancel := context.WithTimeout(context.Background(), gc.interval*9/10)
			err := c.deleteExpired(ctx)
			if err != nil {
				// Log the error if it occurs
				gc.logger.Error("Garbage collection", "error", err)
			}
			cancel()
		case <-gc.stop:
			// Stop the garbage collector
			gc.mu.Lock()
			gc.active = false
			close(gc.stopped) // Signal that the garbage collector has stopped
			gc.mu.Unlock()
			return nil
		}
	}
}

// stopGarbageCollector stops the garbage collector by calling the Stop method.
// It is important to call this function when the cache is no longer
// needed to avoid resource leaks and ensure that the goroutine is
// properly cleaned up.
func stopGarbageCollector[K comparable, V any](gc GarbageCollector[K, V]) {
	// Prevent panic if GC was never started
	if gc == nil {
		return
	}
	// Prepare timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	gc.Stop(ctx)
}

// Stop stops the garbage collector by sending a signal
// to the stop channel. This will cause the garbage collector to exit
// its loop and stop running.
// It is important to call this function when the cache is no longer
// needed to avoid resource leaks and ensure that the goroutine is
// properly cleaned up.
func (gc *garbageCollector[K, V]) Stop(ctx context.Context) error {
	// Prevent panic if GC was never started
	if gc == nil {
		return fmt.Errorf("garbage collector is nil")
	}
	close(gc.stop) // Close the stop channel to signal the goroutine to stop
	select {
	case <-ctx.Done(): // Wait for the context to be done
		return fmt.Errorf("context done: %w", ctx.Err())
	case <-gc.stopped: // Wait for the stopped channel to be closed
		gc.logger.Info("Garbage collector stopped")
	}
	return nil
}

// IsActive returns true if the garbage collector is currently running.
func (gc *garbageCollector[K, V]) IsActive() bool {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.active
}
