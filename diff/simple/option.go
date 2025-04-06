package simple

// defaultShowLineNumbers determines if line numbers are shown by default
const (
	// Formatting options
	defaultContextLines    = 3
	defaultShowLineNumbers = true
)

// options configures the simple diff algorithm
type options struct {
	// Formatting options
	contextLines    int
	showLineNumbers bool
}

// Option is a function that configures the diff options
type Option func(*options)

// WithContextLines sets the number of context lines
func WithContextLines(n int) Option {
	return func(o *options) {
		o.contextLines = n
	}
}

// WithShowLineNumbers configures whether line numbers are shown in the diff output
func WithShowLineNumbers(show bool) Option {
	return func(o *options) {
		o.showLineNumbers = show
	}
}

func applyOptions(opts ...Option) options {
	defaultOpts := options{
		contextLines:    defaultContextLines,
		showLineNumbers: defaultShowLineNumbers,
	}

	for _, opt := range opts {
		opt(&defaultOpts)
	}

	return defaultOpts
}
