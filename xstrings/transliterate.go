package xstrings

import "github.com/gosimple/unidecode"

// Transliterate converts a string to its transliterated form.
func Transliterate(s string) string {
	return unidecode.Unidecode(s)
}
