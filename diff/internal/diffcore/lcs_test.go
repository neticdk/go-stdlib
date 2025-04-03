package diffcore

import (
	"reflect"
	"testing"

	"github.com/neticdk/go-stdlib/diff"
)

func TestComputeEditsLCS(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected []diff.Line
	}{
		{
			name:     "identical slices",
			a:        []string{"hello", "world"},
			b:        []string{"hello", "world"},
			expected: []diff.Line{{Kind: diff.Equal, Text: "hello"}, {Kind: diff.Equal, Text: "world"}},
		},
		{
			name:     "insert element",
			a:        []string{"hello", "world"},
			b:        []string{"hello", "middle", "world"},
			expected: []diff.Line{{Kind: diff.Equal, Text: "hello"}, {Kind: diff.Insert, Text: "middle"}, {Kind: diff.Equal, Text: "world"}},
		},
		{
			name:     "delete element",
			a:        []string{"hello", "middle", "world"},
			b:        []string{"hello", "world"},
			expected: []diff.Line{{Kind: diff.Equal, Text: "hello"}, {Kind: diff.Delete, Text: "middle"}, {Kind: diff.Equal, Text: "world"}},
		},
		{
			name:     "modify element",
			a:        []string{"hello", "old", "world"},
			b:        []string{"hello", "new", "world"},
			expected: []diff.Line{{Kind: diff.Equal, Text: "hello"}, {Kind: diff.Delete, Text: "old"}, {Kind: diff.Insert, Text: "new"}, {Kind: diff.Equal, Text: "world"}},
		},
		{
			name:     "empty slices",
			a:        []string{},
			b:        []string{},
			expected: []diff.Line{},
		},
		{
			name:     "first slice empty",
			a:        []string{},
			b:        []string{"hello", "world"},
			expected: []diff.Line{{Kind: diff.Insert, Text: "hello"}, {Kind: diff.Insert, Text: "world"}},
		},
		{
			name:     "second slice empty",
			a:        []string{"hello", "world"},
			b:        []string{},
			expected: []diff.Line{{Kind: diff.Delete, Text: "hello"}, {Kind: diff.Delete, Text: "world"}},
		},
		{
			name:     "complex case",
			a:        []string{"a", "b", "c", "d", "e"},
			b:        []string{"a", "c", "f", "e"},
			expected: []diff.Line{{Kind: diff.Equal, Text: "a"}, {Kind: diff.Delete, Text: "b"}, {Kind: diff.Equal, Text: "c"},{Kind: diff.Delete, Text: "d"}, {Kind: diff.Insert, Text: "f"}, {Kind: diff.Equal, Text: "e"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeEditsLCS(tt.a, tt.b)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ComputeEditsLCS() = %v, want %v", result, tt.expected)
			}
		})
	}
}
