package myers

import (
	"fmt"

	"github.com/neticdk/go-stdlib/diff"
)

const (
	// Algorithm options
	defaultLinearSpace     = false
	defaultMaxEditDistance = -1 // no constraint

	// defaultLinearRecursionMaxDepth specifies the maximum depth of recursion
	// for the linear space algorithm.
	defaultLinearRecursionMaxDepth = 30

	// defaultSmallInputThreshold specifies the size of input (in number of
	// lines/string slice elements) below which the Myers algorithm is used.
	defaultSmallInputThreshold = 100

	// defaultLargeInputThreshold specifies the size of input (in number of
	// lines/string slice elements) above which the simple diff algorithm is
	// used as a fallback.
	defaultLargeInputThreshold = 10000
)

// options configures the myers diff algorithm
type options struct {
	// Formatting options
	formatter diff.Formatter
	diff.FormatOptions

	// Algorithm options

	// maxEditDistance specifies the maximum edit distance to consider.
	maxEditDistance int
	// linearSpace specifies whether to use linear space algorithm.
	linearSpace bool
	// linearRecursionMaxDepth specifies the maximum depth of recursion for the
	// linear space algorithm.
	linearRecursionMaxDepth int
	// smallInputThreshold specifies the size of input (in number of
	// lines/string slice elements) below which the Myers algorithm is used.
	smallInputThreshold int
	// largeInputThreshold specifies the maximum size of input for which the
	// simple diff algorithm is used.
	largeInputThreshold int
}

// Validate validates the options
func (o *options) validate() error {
	if err := o.Validate(); err != nil {
		return err
	}
	if o.maxEditDistance < -1 {
		return fmt.Errorf("myers: MaxEditDistance cannot be less than -1 (got %d)", o.maxEditDistance)
	}
	if o.smallInputThreshold < 0 || o.largeInputThreshold < 0 || o.linearRecursionMaxDepth < 0 {
		return fmt.Errorf("myers: Thresholds and max depth must be non-negative")
	}
	if o.smallInputThreshold >= o.largeInputThreshold {
		return fmt.Errorf("myers: smallInputThreshold (%d) must be less than largeInputThreshold (%d)", o.smallInputThreshold, o.largeInputThreshold)
	}
	return nil
}

type Option func(*options)

// Option to set the formatter directly
func WithFormatter(f diff.Formatter) Option {
	return func(o *options) {
		if f != nil {
			o.formatter = f
		}
	}
}

// WithContextFormatter uses the context formatter for diff output
func WithContextFormatter() Option {
	return func(o *options) {
		o.formatter = diff.ContextFormatter{}
	}
}

// WithUnifiedFormatter uses the unified formatter for diff output
func WithUnifiedFormatter() Option {
	return func(o *options) {
		o.formatter = diff.UnifiedFormatter{}
	}
}

// WithOutputFormat sets the output format
func WithOutputFormat(format diff.OutputFormat) Option {
	return func(o *options) {
		o.OutputFormat = format
	}
}

// WithContextLines sets the number of context lines
func WithContextLines(n int) Option {
	return func(o *options) {
		o.ContextLines = n
	}
}

// WithShowLineNumbers sets whether to show line numbers
func WithShowLineNumbers(show bool) Option {
	return func(o *options) {
		o.ShowLineNumbers = show
	}
}

// WithMaxEditDistance sets the maximum edit distance
func WithMaxEditDistance(n int) Option {
	return func(o *options) {
		o.maxEditDistance = n
	}
}

// WithLinearSpace sets whether to use linear space algorithm
func WithLinearSpace(linear bool) Option {
	return func(o *options) {
		o.linearSpace = linear
	}
}

// WithSmallInputThreshold sets the size of input (in number of lines/string
// slice elements) below which the Myers algorithm is used.
func WithSmallInputThreshold(size int) Option {
	return func(o *options) {
		o.smallInputThreshold = size
	}
}

// WithLargeInputThreshold sets the maximum size of input for which the simple
// diff algorithm is used.
func WithLargeInputThreshold(size int) Option {
	return func(o *options) {
		o.largeInputThreshold = size
	}
}

// WithLinearRecursionMaxDepth sets the maximum depth of recursion for the
// linear space algorithm
func WithLinearRecursionMaxDepth(depth int) Option {
	return func(o *options) {
		o.linearRecursionMaxDepth = depth
	}
}

func applyOptions(opts ...Option) options {
	defaultOpts := options{
		formatter: diff.ContextFormatter{},
		FormatOptions: diff.FormatOptions{
			OutputFormat:    diff.DefaultOutputFormat,
			ContextLines:    diff.DefaultContextLines,
			ShowLineNumbers: diff.DefaultShowLineNumbers,
		},
		linearSpace:             defaultLinearSpace,
		maxEditDistance:         defaultMaxEditDistance,
		largeInputThreshold:     defaultLargeInputThreshold,
		smallInputThreshold:     defaultSmallInputThreshold,
		linearRecursionMaxDepth: defaultLinearRecursionMaxDepth,
	}

	for _, opt := range opts {
		opt(&defaultOpts)
	}

	return defaultOpts
}
