package cache

// ErrNotFound is an error that indicates that a requested item
// was not found in the cache.
type ErrNotFound struct{}

// Error implements the error interface for ErrNotFound.
func (e *ErrNotFound) Error() string {
	return "cache: not found"
}

// Unwrap implements the Unwrap method for ErrNotFound.
func (e *ErrNotFound) Unwrap() error {
	return nil
}

// Code returns the error code for ErrNotFound.
func (e *ErrNotFound) Code() int {
	return 0
}

// NewErrNotFound creates a new instance of ErrNotFound.
func NewErrNotFound() error {
	return &ErrNotFound{}
}

// ErrExpired is an error that indicates that a requested item
// has expired in the cache.
type ErrExpired struct{}

// Error implements the error interface for ErrExpired.
func (e *ErrExpired) Error() string {
	return "cache: element expired"
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
