package cache

import (
	"testing"
	"time"
)

type mockCache struct {
	*cache
	called int
}

func newMockCache() *mockCache {
	return &mockCache{
		cache: newDefaultCache(),
	}
}

// Mock implementation of the Cache interface for testing
func (m *mockCache) Add(key string, value any) {
	// Mock implementation
}

func (m *mockCache) AddInterval(key string, value any, interval time.Duration) {
	// Mock implementation
}
func (m *mockCache) GetItems() map[string]Item {
	// Mock implementation
	return nil
}
func (m *mockCache) GetAllItems() map[string]Item {
	// Mock implementation
	return nil
}
func (m *mockCache) Get(key string) (any, bool) {
	// Mock implementation
	return nil, false
}
func (m *mockCache) GetTTL(key string) (time.Duration, bool) {
	// Mock implementation
	return 0, false
}
func (m *mockCache) Renew(key string) {
	// Mock implementation
}
func (m *mockCache) RenewInterval(key string, interval time.Duration) {
	// Mock implementation
}
func (m *mockCache) DeleteExpired() {
	// Mock implementation
	m.called++
}
func (m *mockCache) Stop() {
	// Mock implementation
}
func (m *mockCache) Set(key string, value any) {
	// Mock implementation
}
func (m *mockCache) SetInterval(key string, value any, interval time.Duration) {
	// Mock implementation
}
func (m *mockCache) GetInterval(key string) (any, bool) {
	// Mock implementation
	return nil, false
}
func (m *mockCache) Delete(key string) {
	// Mock implementation
}
func (m *mockCache) Clear() {
	// Mock implementation
}
func (m *mockCache) GetGarbageCollector() *garbageCollector {
	return m.garbageCollector
}
func (m *mockCache) SetGarbageCollector(gc *garbageCollector) {
	m.garbageCollector = gc
}
func (m *mockCache) GetClock() Clock {
	return m.clock
}
func (m *mockCache) SetClock(clock Clock) {
	m.clock = clock
}

func TestGarbageCollector(t *testing.T) {
	mc := newMockClock(time.Unix(0, 0).UnixNano()) // Initialize mock clock
	mt := newMockTicker()
	mc.mockTicker = mt // Associate the mock ticker with the mock clock

	// Use newDefaultCache and override the clock
	c := newMockCache()
	c.clock = mc
	// Set a default TTL for items added via Add()
	c.cacheInterval = 5 * time.Second

	mc.Set(time.Unix(0, 0)) // Start time at 0

	// 5. Start Garbage Collector
	gcInterval := 1 * time.Second // Set GC to run every second (mock time)
	// Ensure the nil check is in stopGarbageCollector before running this test
	startGarbageCollector(c, gcInterval) // Starts the gc goroutine
	// Check if GC was actually started (optional)
	if c.garbageCollector == nil {
		t.Fatal("Garbage collector was not started")
	}

	// Advance time past the expiration of "expiredKey" but before "validKey"
	mc.Set(time.Unix(3, 0)) // t=3s

	// Simulate the ticker firing
	mt.Tick(mc.Now()) // Send the tick (should not block due to buffered channel)

	time.Sleep(100 * time.Millisecond) // Allow some time for the tick to be processed

	if c.called != 1 {
		t.Errorf("Expected DeleteExpired to be called once, got %d", c.called)
	}

	// 8. Test Stopping the Garbage Collector
	t.Log("Stopping garbage collector...")
	stopGarbageCollector(c) // Send stop signal

	// Allow time for the stop signal to be potentially processed
	time.Sleep(50 * time.Millisecond)

	// Check if the garbage collector is stopped
	if c.garbageCollector.active {
		t.Error("Garbage collector is still active after stop signal")
	}

	// Cleanup mock ticker if necessary (e.g., close channels)
	mt.Stop()
}
