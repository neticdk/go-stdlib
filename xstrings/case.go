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
	runes := []rune(s)

	for i, r := range runes {
		if isDelimiter(r) {
			// Replace any delimiter with the requested one, avoiding duplicates
			if !prevIsDelimiter && i > 0 {
				builder.WriteString(delimiter)
			}
			prevIsDelimiter = true
			prev = r
			continue
		}

		// We're processing a non-delimiter character
		if i > 0 && !prevIsDelimiter {
			currIsUpper := unicode.IsUpper(r)
			currIsDigit := unicode.IsDigit(r)

			prevIsUpper := unicode.IsUpper(prev)
			prevIsLetter := unicode.IsLetter(prev)
			prevIsDigit := unicode.IsDigit(prev)

			nextIsLower := false
			nextIsS := false
			if i+1 < len(runes) {
				nextIsLower = unicode.IsLower(runes[i+1])
				//Special handling for 's' character
				//"APIs" -> "apis", "HTTPs" -> "https", "URLs" -> "urls"
				nextIsS = (runes[i+1]) == 's'
			}

			// Conditions for adding a delimiter:
			// 1. Lowercase letter -> Uppercase letter (e.g., dC in MixedCase)
			cond1 := prevIsLetter && !prevIsUpper && currIsUpper
			// 2. Uppercase letter -> Uppercase letter and not followed by Lowercase "s" (e.g., TS in HTTPSRequest)
			cond2 := prevIsUpper && currIsUpper && nextIsLower && !nextIsS
			// 3. Uppercase Letter -> Digit (e.g., P2 in HTTP2)
			cond3 := prevIsUpper && currIsDigit
			// 4. Digit -> Uppercase Letter (e.g., 2D in HTTP2Data)
			cond4 := prevIsDigit && currIsUpper

			needsDelimiter := cond1 || cond2 || cond3 || cond4

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
