package xslices

import (
	"testing"
)

func TestFindFunc(t *testing.T) {
	type testCase[T any] struct {
		name     string
		haystack []T
		f        func(T) bool
		expected T
		found    bool
	}

	tests := []testCase[int]{
		{
			name:     "empty slice",
			haystack: []int{},
			f:        func(i int) bool { return i > 0 },
			expected: 0,
			found:    false,
		},
		{
			name:     "match found",
			haystack: []int{1, 2, 3, 4, 5},
			f:        func(i int) bool { return i == 3 },
			expected: 3,
			found:    true,
		},
		{
			name:     "no match",
			haystack: []int{1, 2, 3, 4, 5},
			f:        func(i int) bool { return i > 5 },
			expected: 0,
			found:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, found := FindFunc(tc.haystack, tc.f)
			if actual != tc.expected {
				t.Errorf("FindFunc() = %v, want %v", actual, tc.expected)
			}
			if found != tc.found {
				t.Errorf("FindFunc() found = %v, want %v", found, tc.found)
			}

			actualPtr, foundPtr := FindFunc(tc.haystack, tc.f)
			if actualPtr != tc.expected {
				t.Errorf("FindFunc() = %v, want %v", actualPtr, tc.expected)
			}
			if foundPtr != tc.found {
				t.Errorf("FindFunc() found = %v, want %v", foundPtr, tc.found)
			}
		})
	}
}

func TestFindIFunc(t *testing.T) {
	tests := []struct {
		name          string
		data          []int
		predicate     func(int) bool
		expectedIndex int
		expectedFound bool
	}{
		{
			name:          "element found",
			data:          []int{1, 2, 3, 4, 5},
			predicate:     func(n int) bool { return n == 3 },
			expectedIndex: 2,
			expectedFound: true,
		},
		{
			name:          "element not found",
			data:          []int{1, 2, 3, 4, 5},
			predicate:     func(n int) bool { return n == 6 },
			expectedIndex: -1,
			expectedFound: false,
		},
		{
			name:          "empty slice",
			data:          []int{},
			predicate:     func(n int) bool { return n == 1 },
			expectedIndex: -1,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, found := FindIFunc(tt.data, tt.predicate)
			if index != tt.expectedIndex || found != tt.expectedFound {
				t.Errorf("FindIFunc() = (%d, %v), want (%d, %v)", index, found, tt.expectedIndex, tt.expectedFound)
			}
		})
	}
}
