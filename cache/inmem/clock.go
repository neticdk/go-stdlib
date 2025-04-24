package inmem

import "time"

// Clock is an interface that provides methods to get the current time,
// wait for a duration, and create a new ticker.
// It allows for easy mocking and testing of time-dependent code.
// The default implementation uses the real time package.
type Clock interface {
	// Now returns the current time.
	Now() time.Time

	// NewTicker creates a new Ticker that will send ticks on the
	// returned channel. The ticks will be sent at intervals of d.
	NewTicker(d time.Duration) *time.Ticker
}

// realClock is the default implementation of the Clock interface.
type realClock struct{}

// Now returns the current time using the real time package.
// This is the default implementation and should be used in production code.
func (c *realClock) Now() time.Time {
	return time.Now()
}

// NewTicker creates a new Ticker that will send ticks on the
// returned channel. The ticks will be sent at intervals of d.
// This is the default implementation and should be used in production code.
func (c *realClock) NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
