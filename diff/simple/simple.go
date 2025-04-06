package simple

import (
	"github.com/neticdk/go-stdlib/diff/internal/diffcore"
)

// Diff computes differences between two strings using a simple diff algorithm.
// Returns an error if option validation fails.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := diffcore.SplitLines(a)
	bLines := diffcore.SplitLines(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
// Returns an error if option validation fails.
func DiffStrings(a, b []string, opts ...Option) (string, error) {
	appliedOptions := applyOptions(opts...)
	if err := appliedOptions.validate(); err != nil {
		return "", err
	}
	return simpleDiffStrings(a, b, appliedOptions), nil
}

func simpleDiffStrings(a, b []string, opts options) string {
	// Compute edit script using the shared simple LCS-based diff algorithm
	script := diffcore.ComputeEditsLCS(a, b)
	return opts.formatter.Format(script, opts.FormatOptions)
}
