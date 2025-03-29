package xstrings

import (
	"strings"
	"unicode"
)

// ToKebabCase converts a string to kebab case.
func ToKebabCase(s string) string {
	return ToDelimited(s, "-")
}

// ToSnakeCase converts a string to snake case.
func ToSnakeCase(s string) string {
	return ToDelimited(s, "_")
}

// ToDotCase converts a string to dot case.
func ToDotCase(s string) string {
	return ToDelimited(s, ".")
}

// ToCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	// Estimate the capacity to avoid reallocations
	var builder strings.Builder
	builder.Grow(len(s))

	words := splitIntoWords(s)

	for i, word := range words {
		if word == "" {
			continue
		}

		runes := []rune(word)
		if i == 0 {
			// First word starts with lowercase for camelCase
			builder.WriteRune(unicode.ToLower(runes[0]))
		} else {
			builder.WriteRune(unicode.ToUpper(runes[0]))
		}

		// Add the rest of the word (already lowercased from splitIntoWords)
		if len(runes) > 1 {
			builder.WriteString(string(runes[1:]))
		}
	}

	return builder.String()
}

// ToPascalCase converts a string to pascal case.
func ToPascalCase(s string) string {
	s = ToCamelCase(s)
	if len(s) > 0 {
		s = string(unicode.ToUpper(rune(s[0]))) + s[1:]
	}
	return s
}

// ToDelimited converts a string to a delimited format using the specified
// delimiter
func ToDelimited(s string, delimiter string) string {
	if s == "" {
		return s
	}

	// Allocate a builder with reasonable initial capacity
	var builder strings.Builder
	builder.Grow(len(s) * 2)

	// Single-pass algorithm to handle both case transitions and existing
	// delimiters
	var prev rune
	var prevIsDelimiter bool

	for i, r := range s {
		if isDelimiter(r) {
			// Replace any delimiter with the requested one, avoiding duplicates
			if !prevIsDelimiter && i > 0 {
				builder.WriteString(delimiter)
			}
			prevIsDelimiter = true
			continue
		}

		// We're processing a non-delimiter character
		if i > 0 && !prevIsDelimiter {
			// Add delimiter before uppercase letters
			currIsUpper := unicode.IsUpper(r)
			prevIsUpper := unicode.IsUpper(prev)

			nextIsLower := false
			nextIsS := false
			if i+1 < len(s) {
				nextIsLower = unicode.IsLower(rune(s[i+1]))
				// Special handling for 's' character
				// "APIs" -> "apis", "HTTPs" -> "https", "URLs" -> "urls"
				nextIsS = rune(s[i+1]) == 's'
			}

			isPrevDigit := unicode.IsDigit(prev)
			isCurrentDigit := unicode.IsDigit(r)

			// 1. "camelCase" -> "camel_case"
			// 2. "HTTPRequest" -> "http_request" (but preserves "HTTPs" ->
			//    "https" without delimiter)
			// 3. "user123" -> "user_123" (digit transitions)
			needsDelimiter := (currIsUpper && !prevIsUpper) ||
				(prevIsUpper && currIsUpper && nextIsLower && !nextIsS) ||
				(isCurrentDigit != isPrevDigit)

			if needsDelimiter {
				builder.WriteString(delimiter)
			}
		}

		// Add lowercase version of the current character
		builder.WriteRune(unicode.ToLower(r))
		prev = r
		prevIsDelimiter = false
	}

	result := builder.String()

	// Trim any leading or trailing delimiters
	result = strings.Trim(result, delimiter)

	// Clean up any consecutive delimiters
	if delimiter != "" {
		for strings.Contains(result, delimiter+delimiter) {
			result = strings.ReplaceAll(result, delimiter+delimiter, delimiter)
		}
	}

	return result
}

// splitIntoWords splits the input string into words based on delimiters and
// case changes
func splitIntoWords(s string) []string {
	runes := []rune(s)
	if len(runes) == 0 {
		return nil
	}

	// Pre-allocate a reasonable capacity for the slice
	words := make([]string, 0, len(runes)/4+1)
	wordStart := 0

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Handle delimiters
		if isDelimiter(r) {
			if i > wordStart {
				words = append(words, strings.ToLower(string(runes[wordStart:i])))
			}
			wordStart = i + 1
			continue
		}

		// Check for case transitions if not at the start
		if i > wordStart {
			// Lowercase -> Uppercase transition (new word starts)
			if unicode.IsUpper(r) && !unicode.IsUpper(runes[i-1]) {
				words = append(words, strings.ToLower(string(runes[wordStart:i])))
				wordStart = i
				continue
			}

			// Handle acronyms: Uppercase -> Lowercase transition after multiple
			// uppercase
			// Like in "APIRequest" -> ["api", "request"]
			if i > wordStart+1 && unicode.IsLower(r) && unicode.IsUpper(runes[i-1]) && unicode.IsUpper(runes[i-2]) {
				words = append(words, strings.ToLower(string(runes[wordStart:i-1])))
				wordStart = i - 1
			}
		}
	}

	// Add the last word if there is one
	if wordStart < len(runes) {
		words = append(words, strings.ToLower(string(runes[wordStart:])))
	}

	return words
}

// isDelimiter checks if the rune is a delimiter (space, hyphen, underscore,
// period)
func isDelimiter(r rune) bool {
	return r == ' ' || r == '-' || r == '_' || r == '.'
}
