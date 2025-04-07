package simple

import (
	"github.com/neticdk/go-stdlib/diff"
)

type differ struct{}

// NewDiffer creates a new SimpleDiffer instance
func NewDiffer() diff.Differ {
	return &differ{}
}

// Diff computes differences between two strings using a simple diff algorithm.
// Returns an error if option validation fails.
func (d *differ) Diff(a, b string) (string, error) {
	return Diff(a, b)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
// Returns an error if option validation fails.
func (s *differ) DiffStrings(a, b []string) (string, error) {
	return DiffStrings(a, b)
}

type customDiffer struct {
	opts []Option
}

// NewCustomDiffer creates a new CustomSimpleDiffer instance
func NewCustomDiffer(opts ...Option) diff.Differ {
	return &customDiffer{opts: opts}
}

// Diff computes differences between two strings using a simple diff algorithm.
// Returns an error if option validation fails.
func (d *customDiffer) Diff(a, b string) (string, error) {
	return Diff(a, b, d.opts...)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
// Returns an error if option validation fails.
func (d *customDiffer) DiffStrings(a, b []string) (string, error) {
	return DiffStrings(a, b, d.opts...)
}
