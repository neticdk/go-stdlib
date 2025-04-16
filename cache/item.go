package cache

import "time"

// Item represents a cache item with a value and an expiration time.
// The expiration time is set to the current time plus the TTL (time-to-live).
// The item is considered expired if the current time is greater than
// the expiration time.
type Item struct {
	Value      any
	Expiration time.Time
}

// isExpired checks if the item is expired based on the current time
// from the provided clock.
func (i *Item) isExpired(clock Clock) bool {
	return clock.Now().After(i.Expiration)
}
