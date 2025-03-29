package xslices

import (
	"reflect"
	"testing"
)

func TestUnfold(t *testing.T) {
	tests := []struct {
		name     string
		acc      any
		f        func(any) any
		p        func(any) bool
		opts     []UnfoldOption
		expected []any
	}{
		{
			name:     "double",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			p:        func(acc any) bool { return acc.(int) < 100 },
			expected: []any{1, 2, 4, 8, 16, 32, 64},
		},
		{
			name:     "always false p",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			p:        func(acc any) bool { return false },
			expected: nil,
		},
		{
			name:     "always true p",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			p:        func(acc any) bool { return true },
			opts:     []UnfoldOption{WithMax(5)},
			expected: []any{1, 2, 4, 8, 16, 32},
		},
		{
			name:     "runes",
			acc:      'a',
			f:        func(acc any) any { return acc.(rune) + 1 },
			p:        func(acc any) bool { return acc.(rune) < 'f' },
			expected: []any{'a', 'b', 'c', 'd', 'e'},
		},
		{
			name:     "strings",
			acc:      "a",
			f:        func(acc any) any { return acc.(string) + "a" },
			p:        func(acc any) bool { return len(acc.(string)) < 5 },
			expected: []any{"a", "aa", "aaa", "aaaa"},
		},
		{
			name:     "with step",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			p:        func(acc any) bool { return acc.(int) < 100 },
			opts:     []UnfoldOption{WithStep(2)},
			expected: []any{1, 4, 16, 64},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Unfold(tt.acc, tt.f, tt.p, tt.opts...)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v but got %v", tt.expected, actual)
			}
		})
	}
}

func TestUnfoldI(t *testing.T) {
	tests := []struct {
		name     string
		acc      any
		f        func(any) any
		i        int
		opts     []UnfoldOption
		expected []any
	}{
		{
			name:     "integers",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			i:        7,
			expected: []any{1, 2, 4, 8, 16, 32, 64, 128},
		},
		{
			name:     "negative i",
			acc:      1,
			f:        func(acc any) any { return acc.(int) * 2 },
			i:        -7,
			expected: nil,
		},
		{
			name:     "strings",
			acc:      "a",
			f:        func(acc any) any { return acc.(string) + "a" },
			i:        5,
			expected: []any{"a", "aa", "aaa", "aaaa", "aaaaa", "aaaaaa"},
		},
		{
			name: "with max",
			acc:  1,
			f:    func(acc any) any { return acc.(int) * 2 },
			i:    7,
			opts: []UnfoldOption{WithMax(5)},
			expected: []any{
				1, 2, 4, 8, 16, 32,
			},
		},
		{
			name: "with step",
			acc:  1,
			f:    func(acc any) any { return acc.(int) * 2 },
			i:    7,
			opts: []UnfoldOption{WithStep(2)},
			expected: []any{
				1, 4, 16, 64,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := UnfoldI(tt.acc, tt.f, tt.i, tt.opts...)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v but got %v", tt.expected, actual)
			}
		})
	}
}
