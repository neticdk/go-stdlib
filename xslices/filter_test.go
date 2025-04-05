package xslices

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) bool
		expected []int
	}{
		{
			name:  "filter even numbers",
			input: []int{1, 2, 3, 4, 5, 6},
			fn: func(i int) bool {
				return i%2 == 0
			},
			expected: []int{2, 4, 6},
		},
		{
			name:  "empty slice",
			input: []int{},
			fn: func(i int) bool {
				return true
			},
			expected: []int{},
		},
		{
			name:  "filter all elements",
			input: []int{1, 2, 3},
			fn: func(i int) bool {
				return false
			},
			expected: []int{},
		},
		{
			name:  "keep all elements",
			input: []int{1, 2, 3},
			fn: func(i int) bool {
				return true
			},
			expected: []int{1, 2, 3},
		},
		{
			name:  "filter specific value",
			input: []int{1, 42, 2, 42, 3},
			fn: func(i int) bool {
				return i == 42
			},
			expected: []int{42, 42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.fn)
			assert.Equal(t, result, tt.expected, "Filter/%q", tt.name)
		})
	}
}
