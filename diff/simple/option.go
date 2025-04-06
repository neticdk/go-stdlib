package simple

import "github.com/neticdk/go-stdlib/diff"

// options configures the simple diff algorithm
type options struct {
	// Formatting options
	formatter diff.Formatter
	diff.FormatOptions
}

// validate validates the options
func (o *options) validate() error {
	// Call validate on the embedded FormatOptions
	return o.Validate()
}

// Option is a function that configures the diff options
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

// WithShowLineNumbers configures whether line numbers are shown in the diff output
func WithShowLineNumbers(show bool) Option {
	return func(o *options) {
		o.ShowLineNumbers = show
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
	}

	for _, opt := range opts {
		opt(&defaultOpts)
	}

	return defaultOpts
}
