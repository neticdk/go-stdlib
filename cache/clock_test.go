package cache

import (
	"sync"
	"time"
)

// newMockClock creates a new test clock with the specified time.
// This is useful for testing time-dependent code without relying on
// the actual system time.
func newMockClock(t int64) *mockClock {
	return &mockClock{
		mu:         sync.Mutex{},
		now:        time.Unix(t, 0),
		mockTicker: newMockTicker(),
	}
}

// mockTicker is a mock implementation of time.Ticker for testing.
type mockTicker struct {
	C    chan time.Time // Channel similar to time.Ticker.C
	stop chan bool      // Internal channel to signal stopping the mock ticker goroutine
}

// NewMockTicker creates a new mockTicker.
func newMockTicker() *mockTicker {
	return &mockTicker{
		// Use a buffered channel to prevent Tick() from blocking indefinitely
		// if the receiver isn't ready immediately.
		C:    make(chan time.Time, 1), // Buffered channel with capacity 1
		stop: make(chan bool),
	}
}

// Tick simulates a ticker event by sending the current time on the channel.
// Note: If the channel buffer is full, this could still block.
func (mt *mockTicker) Tick(t time.Time) {
	mt.C <- t
}

// Stop stops the mock ticker. In a real scenario, this might clean up internal goroutines if any.
func (mt *mockTicker) Stop() {
	close(mt.stop) // Signal any internal goroutines to stop
}

// mockClock is a mock implementation of the Clock interface.
type mockClock struct {
	mu         sync.Mutex
	now        time.Time
	mockTicker *mockTicker // The ticker to return from NewTicker
}

// Now returns the current mock time.
func (mc *mockClock) Now() time.Time {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return mc.now
}

// Set sets the current mock time.
func (mc *mockClock) Set(t time.Time) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.now = t
}

// Add advances the current mock time by the specified duration.
func (mc *mockClock) Add(d time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.now = mc.now.Add(d)
}

// NewTicker returns the pre-configured mockTicker.
// NOTE: This mock implementation returns the *same* ticker instance every time
// for simplicity in this test example.
func (mc *mockClock) NewTicker(d time.Duration) *time.Ticker {
	// Return a real ticker but the test will interact with mc.mockTicker directly.
	if mc.mockTicker == nil {
		panic("mockClock.mockTicker is not set") // Ensure mock ticker is set before use
	}

	realTicker := &time.Ticker{C: mc.mockTicker.C} // Directly assign the mock channel
	return realTicker
}

// After is not needed for GC test, but required by interface.
func (mc *mockClock) After(d time.Duration) <-chan time.Time {
	// Return a channel that closes after mock time advances past duration d
	return time.After(1 * time.Millisecond) // Placeholder, not used in GC test
}
