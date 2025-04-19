package xstrings

import "github.com/neticdk/go-stdlib/xslices"

// Coalesce returns the first non-empty string from the given slice.
func Coalesce(strs ...string) string {
	s, found := xslices.FindFunc(strs, func(s string) bool {
		return s != ""
	})
	if found {
		return s
	}
	return ""
}
