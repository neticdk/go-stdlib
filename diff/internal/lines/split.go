package lines

import "strings"

// Split splits a string into lines, handling empty strings and trailing newlines
func Split(s string) []string {
	if s == "" {
		return []string{}
	}

	lines := strings.Split(s, "\n")

	// If the string ends with a newline, the split will produce an empty string
	// at the end - we should remove it to avoid confusing diff output
	if s[len(s)-1] == '\n' {
		lines = lines[:len(lines)-1]
	}

	return lines
}
