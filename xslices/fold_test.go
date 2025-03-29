package xslices

import (
	"reflect"
	"testing"
)

func TestFold(t *testing.T) {
	tests := []struct {
		name     string
		data     []any
		acc      any
		f        func(any, any) any
		expected any
	}{
		{
			name:     "sum",
			data:     []any{1, 2, 3, 4, 5},
			acc:      0,
			f:        func(acc, e any) any { return acc.(int) + e.(int) },
			expected: 15,
		},
		{
			name:     "concat",
			data:     []any{"a", "b", "c", "d", "e"},
			acc:      "",
			f:        func(acc, e any) any { return acc.(string) + e.(string) },
			expected: "abcde",
		},
		{
			name: "product",
			data: []any{2, 3, 4, 5, 6},
			acc:  1,
			f: func(acc, e any) any {
				return acc.(int) * e.(int)
			},
			expected: 720,
		},
		{
			name:     "concat with acc",
			data:     []any{"a", "b", "c", "d", "e"},
			acc:      "result is: ",
			f:        func(acc, e any) any { return acc.(string) + e.(string) },
			expected: "result is: abcde",
		},
		{
			name: "reverse",
			data: []any{1, 2, 3, 4, 5},
			acc:  []any{},
			f:    func(acc, e any) any { return append([]any{e}, acc.([]any)...) },
			expected: []any{
				5, 4, 3, 2, 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Fold(tt.acc, tt.data, tt.f)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

func TestFoldR(t *testing.T) {
	tests := []struct {
		name     string
		data     []any
		acc      any
		f        func(any, any) any
		expected any
	}{
		{
			name:     "sum",
			data:     []any{1, 2, 3, 4, 5},
			acc:      0,
			f:        func(acc, e any) any { return acc.(int) + e.(int) },
			expected: 15,
		},
		{
			name:     "concat",
			data:     []any{"a", "b", "c", "d", "e"},
			acc:      "",
			f:        func(acc, e any) any { return acc.(string) + e.(string) },
			expected: "edcba",
		},
		{
			name: "product",
			data: []any{2, 3, 4, 5, 6},
			acc:  1,
			f: func(acc, e any) any {
				return acc.(int) * e.(int)
			},
			expected: 720,
		},
		{
			name:     "concat with acc",
			data:     []any{"a", "b", "c", "d", "e"},
			acc:      "result is: ",
			f:        func(acc, e any) any { return acc.(string) + e.(string) },
			expected: "result is: edcba",
		},
		{
			name: "reverse",
			data: []any{1, 2, 3, 4, 5},
			acc:  []any{},
			f:    func(acc, e any) any { return append(acc.([]any), e) },
			expected: []any{
				5, 4, 3, 2, 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FoldR(tt.acc, tt.data, tt.f)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}
