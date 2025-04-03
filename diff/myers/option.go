package myers

const (
	defaultContextLines    = 3
	defaultLinearSpace     = false
	defaultShowLineNumbers = true
	defaultMaxEditDistance = -1 // no constraint

	// DefaultLinearRecursionMaxDepth specifies the maximum depth of recursion
	// for the linear space algorithm.
	DefaultLinearRecursionMaxDepth = 30

	// DefaultSimpleDiffFallbackSize specifies the maximum size of input for which the
	// simple diff algorithm is used.
	DefaultSimpleDiffFallbackSize = 10000
)

type Option func(*options)

// options configures the myers diff algorithm
type options struct {
	maxEditDistance         int  // MaxEditDistance specifies the maximum edit distance to consider.
	linearSpace             bool // LinearSpace specifies whether to use linear space algorithm.
	linearRecursionMaxDepth int  // LinearRecursionMaxDepth specifies the maximum depth of recursion for the linear space algorithm.
	simpleDiffFallbackSize  int  // SimpleDiffFallbackSize specifies the maximum size of input for which the simple diff algorithm is used.

	// Formatting options
	contextLines    int
	showLineNumbers bool
}

// WithContextLines sets the number of context lines
func WithContextLines(n int) Option {
	return func(o *options) {
		o.contextLines = n
	}
}

// WithShowLineNumbers sets whether to show line numbers
func WithShowLineNumbers(show bool) Option {
	return func(o *options) {
		o.showLineNumbers = show
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

// WithSimpleDiffFallbackSize sets the maximum size of input for which the
// simple diff algorithm is used
func WithSimpleDiffFallbackSize(size int) Option {
	return func(o *options) {
		o.simpleDiffFallbackSize = size
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
		linearSpace:             defaultLinearSpace,
		contextLines:            defaultContextLines,
		showLineNumbers:         defaultShowLineNumbers,
		maxEditDistance:         defaultMaxEditDistance,
		simpleDiffFallbackSize:  DefaultSimpleDiffFallbackSize,
		linearRecursionMaxDepth: DefaultLinearRecursionMaxDepth,
	}

	for _, opt := range opts {
		opt(&defaultOpts)
	}

	return defaultOpts
}
