package myers

import (
	"github.com/neticdk/go-stdlib/diff"
)

type differ struct{}

// NewMyersDiffer creates a new MyersDiffer instance
func NewDiffer() diff.Differ {
	return &differ{}
}

func (d *differ) Diff(a, b string) (string, error) {
	return Diff(a, b)
}

func (d *differ) DiffStrings(a, b []string) (string, error) {
	return DiffStrings(a, b)
}

type customDiffer struct {
	opts []Option
}

// NewCustomDiffer creates a new CustomMyersDiffer instance
func NewCustomDiffer(opts ...Option) diff.Differ {
	return &customDiffer{opts: opts}
}

func (d *customDiffer) Diff(a, b string) (string, error) {
	return Diff(a, b, d.opts...)
}

func (d *customDiffer) DiffStrings(a, b []string) (string, error) {
	return DiffStrings(a, b, d.opts...)
}
