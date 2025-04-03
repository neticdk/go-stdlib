package myers

const (
	defaultContextLines    = 3
	defaultLinearSpace     = false
	defaultShowLineNumbers = true
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

type Option func(*options)

// options configures the myers diff algorithm
type options struct {
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
		linearSpace:             defaultLinearSpace,
		contextLines:            defaultContextLines,
		showLineNumbers:         defaultShowLineNumbers,
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
