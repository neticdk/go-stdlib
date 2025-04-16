package cache

import "time"

// garbageCollector is responsible for periodically deleting expired items
// from the cache. It runs in a separate goroutine and uses a ticker to
// trigger the deletion process at regular intervals.
type garbageCollector struct {
	// active indicates whether the garbage collector is currently running.
	active bool

	// interval is the duration between each garbage collection cycle.
	interval time.Duration

	// stop is a channel used to signal the garbage collector to stop.
	stop chan bool
}

// start begins the garbage collection process. It runs in a separate
// goroutine and periodically checks for expired items in the cache.
// When a tick happens, it calls the DeleteExpired method on the cache to
// remove expired items.
func (gc *garbageCollector) start(c Cache) {
	ticker := c.GetClock().NewTicker(gc.interval)
	defer ticker.Stop()
	gc.active = true

	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
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
func stopGarbageCollector(c Cache) {
	// Prevent panic if GC was never started
	gc := c.GetGarbageCollector()
	if gc == nil {
		return
	}
	gc.stop <- true
}

// startGarbageCollector initializes and starts the garbage collector
// for the cache. It creates a new garbageCollector instance with the
// specified interval and assigns it to the cache's garbageCollector field.
// If the interval is zero or negative, the garbage collector will not be started.
func startGarbageCollector(c Cache, interval time.Duration) {
	if interval <= 0 {
		return
	}
	gc := &garbageCollector{
		interval: interval,
		stop:     make(chan bool),
		active:   true,
	}
	c.SetGarbageCollector(gc)
	go gc.start(c)
}
