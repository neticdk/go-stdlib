package myers

type Option func(*options)

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

func applyOptions(opts ...Option) options {
	defaultOpts := options{
		linearSpace:     DefaultLinearSpace,
		contextLines:    DefaultContextLines,
		showLineNumbers: DefaultShowLineNumbers,
		maxEditDistance: DefaultMaxEditDistance,
	}

	for _, opt := range opts {
		opt(&defaultOpts)
	}

	return defaultOpts
}
