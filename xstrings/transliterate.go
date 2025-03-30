package xstrings

import "github.com/neticdk/stdlib/xstrings/transliterate"

// Transliterate converts a string to its transliterated form.
//
// Deprecated: Use the String() method from github.com/neticdk/stdlib/xstrings/transliterate package directly.
func Transliterate(s string) string {
	return transliterate.String(s)
}
