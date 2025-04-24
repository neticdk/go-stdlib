package cache

// ErrCacheMiss is an error that indicates that a requested item
// was not found in the cache.
type ErrCacheMiss struct{}

// Error implements the error interface for ErrCacheMiss.
func (e *ErrCacheMiss) Error() string {
	return "cache: miss"
}

// Unwrap implements the Unwrap method for ErrCacheMiss.
func (e *ErrCacheMiss) Unwrap() error {
	return nil
}

// Code returns the error code for ErrCacheMiss.
func (e *ErrCacheMiss) Code() int {
	return 0
}

// NewErrCacheMiss creates a new instance of ErrCacheMiss.
func NewErrCacheMiss() error {
	return &ErrCacheMiss{}
}

// ErrExpired is an error that indicates that a requested item
// has expired in the cache.
type ErrExpired struct{}

// Error implements the error interface for ErrExpired.
func (e *ErrExpired) Error() string {
	return "cache: item expired"
}

// Unwrap implements the Unwrap method for ErrExpired.
func (e *ErrExpired) Unwrap() error {
	return nil
}

// Code returns the error code for ErrExpired.
func (e *ErrExpired) Code() int {
	return 0
}

// NewErrExpired creates a new instance of ErrExpired.
func NewErrExpired() error {
	return &ErrExpired{}
}

// ErrCacheFull is an error that indicates that the cache is full
// and cannot accept new items.
type ErrCacheFull struct{}

// Error implements the error interface for ErrCacheFull.
func (e *ErrCacheFull) Error() string {
	return "cache: full"
}

// Unwrap implements the Unwrap method for ErrCacheFull.
func (e *ErrCacheFull) Unwrap() error {
	return nil
}

// Code returns the error code for ErrCacheFull.
func (e *ErrCacheFull) Code() int {
	return 0
}

// NewErrCacheFull creates a new instance of ErrCacheFull.
func NewErrCacheFull() error {
	return &ErrCacheFull{}
}

// ErrCacheNotStopped is an error that indicates that the cache
// has not been stopped properly.
type ErrCacheNotStopped struct{}

// Error implements the error interface for ErrCacheNotStopped.
func (e *ErrCacheNotStopped) Error() string {
	return "cache: not stopped"
}

// Unwrap implements the Unwrap method for ErrCacheNotStopped.
func (e *ErrCacheNotStopped) Unwrap() error {
	return nil
}

// Code returns the error code for ErrCacheNotStopped.
func (e *ErrCacheNotStopped) Code() int {
	return 0
}

// NewErrCacheNotStopped creates a new instance of ErrCacheNotStopped.
func NewErrCacheNotStopped() error {
	return &ErrCacheNotStopped{}
}
