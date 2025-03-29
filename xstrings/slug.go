package xstrings

import (
	"fmt"
	"regexp"
	"strings"
)

// Slugify converts a string to a slug.
// By default, it converts the string to lowercase, transliterates it, and
// removes camel case.
func Slugify(s string, options ...TransformOption) string {
	opts := TransformOptions{
		lowercase:     true,
		transliterate: true,
		decamelize:    true,
		delimiter:     DefaultDelimiter,
	}

	for _, option := range options {
		option(&opts)
	}

	if opts.transliterate {
		s = Transliterate(s)
	}

	if opts.decamelize {
		s = ToDelimited(s, opts.delimiter)
	}

	var (
		reLowerCaseSeparator   = regexp.MustCompile(fmt.Sprintf(`[^a-z\d%s]+`, opts.delimiter))
		reMixedCaseSeparator   = regexp.MustCompile(fmt.Sprintf(`[^a-zA-Z\d%s]+`, opts.delimiter))
		reConsecutiveSeparator = regexp.MustCompile(fmt.Sprintf(`%s+`, opts.delimiter))
	)

	if opts.lowercase {
		s = strings.ToLower(s)
		s = reLowerCaseSeparator.ReplaceAllString(s, opts.delimiter)
	} else {
		s = reMixedCaseSeparator.ReplaceAllString(s, opts.delimiter)
	}

	// Replace multiple consecutive separators with a single separator.
	s = reConsecutiveSeparator.ReplaceAllString(s, opts.delimiter)

	// Remove leading and trailing hyphens.
	s = strings.Trim(s, "-")

	return s
}
