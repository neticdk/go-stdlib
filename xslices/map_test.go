package xslices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) string
		expected []string
	}{
		{
			name:  "convert ints to strings",
			input: []int{1, 2, 3, 4, 5},
			fn: func(i int) string {
				return string(rune('A' - 1 + i))
			},
			expected: []string{"A", "B", "C", "D", "E"},
		},
		{
			name:  "empty slice",
			input: []int{},
			fn: func(i int) string {
				return string(rune('A' - 1 + i))
			},
			expected: []string{},
		},
		{
			name:  "single element",
			input: []int{42},
			fn: func(i int) string {
				if i == 42 {
					return "answer"
				}
				return "wrong"
			},
			expected: []string{"answer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.fn)
			assert.Equal(t, tt.expected, result)
		})
	}
}
