package simple

import (
	"github.com/neticdk/go-stdlib/diff"
)

type differ struct{}

// NewDiffer creates a new SimpleDiffer instance
func NewDiffer() diff.Differ {
	return &differ{}
}

func (d *differ) Diff(a, b string) string {
	return Diff(a, b)
}

func (s *differ) DiffStrings(a, b []string) string {
	return DiffStrings(a, b)
}

type customDiffer struct {
	opts []Option
}

// NewCustomDiffer creates a new CustomSimpleDiffer instance
func NewCustomDiffer(opts ...Option) diff.Differ {
	return &customDiffer{opts: opts}
}

func (d *customDiffer) Diff(a, b string) string {
	return Diff(a, b, d.opts...)
}

func (d *customDiffer) DiffStrings(a, b []string) string {
	return DiffStrings(a, b, d.opts...)
}
