package myers

type Option func(*options)

// WithContextLines sets the number of context lines
func WithContextLines(n int) Option {
	return func(o *options) {
		o.contextLines = n
	}
}

func WithShowLineNumbers(show bool) Option {
	return func(o *options) {
		o.showLineNumbers = show
	}
}

func WithMaxEditDistance(n int) Option {
	return func(o *options) {
		o.maxEditDistance = n
	}
}

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
