package xtime

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestClean(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "1 hour 30 minutes",
			expected: "1 hour 30 minutes",
		},
		{
			input:    "1h30m",
			expected: "1h30m",
		},
		{
			input:    "1.5 seconds",
			expected: "1.5 seconds",
		},
		{
			input:    "1 year, 2 hour, and 5s",
			expected: "1 year  2 hour  5s",
		},
		{
			input:    ",1 years",
			expected: " 1 years",
		},
		{
			input:    "and 1 minute",
			expected: "1 minute",
		},
		{
			input:    ", 4 minutes",
			expected: "  4 minutes",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := clean(tc.input)
			assert.Equal(t, actual, tc.expected)
		})
	}
}
