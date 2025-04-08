package diffcore_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/diff/internal/diffcore"
)

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single line",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "multiple lines",
			input:    "hello\nworld",
			expected: []string{"hello", "world"},
		},
		{
			name:     "string with trailing newline",
			input:    "hello\nworld\n",
			expected: []string{"hello", "world"},
		},
		{
			name:     "string with only newline",
			input:    "\n",
			expected: []string{""},
		},
		{
			name:     "multiple newlines",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "leading and trailing newlines",
			input:    "\nline1\nline2\nline3\n",
			expected: []string{"", "line1", "line2", "line3"},
		},
		{
			name:     "multiple consecutive newlines",
			input:    "line1\n\nline2",
			expected: []string{"line1", "", "line2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diffcore.SplitLines(tt.input)
			assert.Equal(t, result, tt.expected, "SplitLines/%q", tt.name)
		})
	}
}
