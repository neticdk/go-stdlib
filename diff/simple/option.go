package simple

// DefaultShowLineNumbers determines if line numbers are shown by default
var DefaultShowLineNumbers = true

// options configures the simple diff algorithm
type options struct {
	showLineNumbers bool
}

// Option is a function that configures the diff options
type Option func(*options)

// WithShowLineNumbers configures whether line numbers are shown in the diff output
func WithShowLineNumbers(show bool) Option {
	return func(o *options) {
		o.showLineNumbers = show
	}
}

func applyOptions(opts ...Option) options {
	defaultOpts := options{
		showLineNumbers: DefaultShowLineNumbers,
	}

	for _, opt := range opts {
		opt(&defaultOpts)
	}

	return defaultOpts
}
