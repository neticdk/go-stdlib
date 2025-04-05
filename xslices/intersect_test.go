package xslices

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestIntersectionOfTwoSlices(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected []int
	}{
		{
			name:     "common elements",
			a:        []int{1, 2, 3},
			b:        []int{2, 3, 4},
			expected: []int{2, 3},
		},
		{
			name:     "no common elements",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: []int{},
		},
		{
			name:     "empty first slice",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "empty second slice",
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "both slices empty",
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "duplicate elements in first slice",
			a:        []int{1, 2, 2, 3},
			b:        []int{2, 3, 4},
			expected: []int{2, 3},
		},
		{
			name:     "duplicate elements in second slice",
			a:        []int{1, 2, 3},
			b:        []int{2, 2, 3, 4},
			expected: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersect(tt.a, tt.b)
			assert.Equal(t, result, tt.expected, "Intersect/%q", tt.name)
		})
	}
}
