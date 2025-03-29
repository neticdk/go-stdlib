package xstrings

const DefaultDelimiter = "-"

type TransformOption func(*TransformOptions)

// TransformOptions represents the options for transforming a string.
type TransformOptions struct {
	// lowercase is used to convert the string to lowercase.
	lowercase bool
	// decamelize is used to remove camelcase.
	decamelize bool
	// transliterate is used to transliterate the string.
	transliterate bool
	// delimiter is used to specify a delimiter for separating words.
	delimiter string
}

// WithTransliterate sets the transliterate option.
func WithTransliterate(transliterate bool) TransformOption {
	return func(o *TransformOptions) {
		o.transliterate = transliterate
	}
}

// WithLowercase sets the lowercase option.
func WithLowercase(lowercase bool) TransformOption {
	return func(o *TransformOptions) {
		o.lowercase = lowercase
	}
}

// WithDecamelize sets the decamelize option.
func WithDecamelize(decamelize bool) TransformOption {
	return func(o *TransformOptions) {
		o.decamelize = decamelize
	}
}

// WithDelimiter sets the delimiter option.
func WithDelimiter(delimiter string) TransformOption {
	return func(o *TransformOptions) {
		o.delimiter = delimiter
	}
}
