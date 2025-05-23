package inmem

import (
	"context"
	"testing"
	"time"
)

// mockClock is a mock implementation of the safeMapCache struct.
// It allows for testing the calls to methods of safeMapCache
// without using the actual implementations.
type mockCache[K comparable, V any] struct {
	// embeds the safeMapCache struct to inherit its methods
	*safeMapCache[K, V]

	// mutex for synchronizing access to deleteExpiredCalled
	called chan bool
}

// newMockCache creates a new instance of mockCache.
func newMockCache[K comparable, V any]() *mockCache[K, V] {
	return &mockCache[K, V]{
		safeMapCache: NewSafeMap[K, V](),
		called:       make(chan bool, 1),
	}
}

// deleteExpired is a mock implementation of the deleteExpired method.
// It increments the deleteExpiredCalled counter each time it is called.
func (m *mockCache[K, V]) deleteExpired(ctx context.Context) error {
	// Mock implementation
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called <- true
	return nil
}

func TestGarbageCollector(t *testing.T) {
	mc := newMockClock(time.Unix(0, 0).UnixNano()) // Initialize mock clock
	mt := newMockTicker()
	mc.mockTicker = mt // Associate the mock ticker with the mock clock

	gcInterval := 1 * time.Second                 // Set GC to run every second (mock time)
	gc := newTestGarbageCollector(gcInterval, mc) // Starts the gc goroutine

	// Use newDefaultCache and override the clock
	c := newMockCache[string, any]()
	c.garbageCollector = gc
	go c.garbageCollector.Start(c) // Start the garbage collector
	c.clock = mc
	// Set a default TTL for items added via Add()

	mc.Set(time.Unix(0, 0)) // Start time at 0

	// 5. Start Garbage Collector
	// Ensure the nil check is in stopGarbageCollector before running this test
	// Check if GC was actually started (optional)
	if c.garbageCollector == nil {
		t.Fatal("Garbage collector was not started")
	}

	// Advance time past the expiration of "expiredKey" but before "validKey"
	mc.Set(time.Unix(3, 0)) // t=3s

	// Simulate the ticker firing
	mt.Tick(mc.Now()) // Send the tick (should not block due to buffered channel)

	<-c.called

	// 8. Test Stopping the Garbage Collector
	//stopChan := make(chan bool, 1) // Channel to signal stop
	t.Log("Stopping garbage collector...")
	c.garbageCollector.Stop(context.Background()) // Send stop signal

	// <-stopChan // Wait for the stop signal to be processed

	// Check if the garbage collector is stopped
	if c.garbageCollector.IsActive() {
		t.Error("Garbage collector is still active after stop signal")
	}

	// Cleanup mock ticker if necessary (e.g., close channels)
	mt.Stop()
}
