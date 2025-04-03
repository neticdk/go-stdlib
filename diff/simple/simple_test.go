package simple_test

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/diff/simple"
)

func TestSimpleDifferInterface(t *testing.T) {
	// Test data for both string and string slice inputs
	stringTests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{
			name: "identical strings",
			a:    "hello\nworld",
			b:    "hello\nworld",
			expected: "   1    1   hello\n" +
				"   2    2   world\n",
		},
		{
			name: "insert line",
			a:    "hello\nworld",
			b:    "hello\nmiddle\nworld",
			expected: "   1    1   hello\n" +
				"        2 + middle\n" +
				"   2    3   world\n",
		},
		{
			name:     "empty strings",
			a:        "",
			b:        "",
			expected: "",
		},
	}

	sliceTests := []struct {
		name     string
		a        []string
		b        []string
		expected string
	}{
		{
			name: "identical slices",
			a:    []string{"hello", "world"},
			b:    []string{"hello", "world"},
			expected: "   1    1   hello\n" +
				"   2    2   world\n",
		},
		{
			name: "insert element",
			a:    []string{"hello", "world"},
			b:    []string{"hello", "middle", "world"},
			expected: "   1    1   hello\n" +
				"        2 + middle\n" +
				"   2    3   world\n",
		},
		{
			name:     "empty slices",
			a:        []string{},
			b:        []string{},
			expected: "",
		},
	}

	// Test default differ
	differ := simple.NewDiffer()

	// Test string input
	for _, tt := range stringTests {
		t.Run("Default/String/"+tt.name, func(t *testing.T) {
			result, err := differ.Diff(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Error running Diff with default differ: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}

	// Test string slice input
	for _, tt := range sliceTests {
		t.Run("Default/Slice/"+tt.name, func(t *testing.T) {
			result, err := differ.DiffStrings(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Error running DiffStrings with default differ: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}

	// Test custom differ with options
	customDiffer := simple.NewCustomDiffer(
		simple.WithShowLineNumbers(false),
	)

	// Custom expected results for options
	customExpected := "  hello\n+ middle\n  world\n"

	t.Run("Custom/String", func(t *testing.T) {
		result, err := customDiffer.Diff("hello\nworld", "hello\nmiddle\nworld")
		if err != nil {
			t.Fatalf("Error running Diff with custom differ: %v", err)
		}
		if result != customExpected {
			t.Errorf("Expected:\n%s\nGot:\n%s", customExpected, result)
		}
	})

	t.Run("Custom/Slice", func(t *testing.T) {
		result, err := customDiffer.DiffStrings(
			[]string{"hello", "world"},
			[]string{"hello", "middle", "world"},
		)
		if err != nil {
			t.Fatalf("Error running DiffStrings with custom differ: %v", err)
		}
		if result != customExpected {
			t.Errorf("Expected:\n%s\nGot:\n%s", customExpected, result)
		}
	})
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{
			name: "identical strings",
			a:    "hello\nworld",
			b:    "hello\nworld",
			expected: "   1    1   hello\n" +
				"   2    2   world\n",
		},
		{
			name: "insert line",
			a:    "hello\nworld",
			b:    "hello\nmiddle\nworld",
			expected: "   1    1   hello\n" +
				"        2 + middle\n" +
				"   2    3   world\n",
		},
		{
			name: "delete line",
			a:    "hello\nmiddle\nworld",
			b:    "hello\nworld",
			expected: "   1    1   hello\n" +
				"   2      - middle\n" +
				"   3    2   world\n",
		},
		{
			name: "modify line",
			a:    "hello\nold line\nworld",
			b:    "hello\nnew line\nworld",
			expected: "   1    1   hello\n" +
				"   2      - old line\n" +
				"        2 + new line\n" +
				"   3    3   world\n",
		},
		{
			name:     "empty strings",
			a:        "",
			b:        "",
			expected: "",
		},
		{
			name: "first string empty",
			a:    "",
			b:    "hello\nworld",
			expected: "        1 + hello\n" +
				"        2 + world\n",
		},
		{
			name: "second string empty",
			a:    "hello\nworld",
			b:    "",
			expected: "   1      - hello\n" +
				"   2      - world\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := simple.Diff(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Error running Diff: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestDiffStrings(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected string
	}{
		{
			name: "identical slices",
			a:    []string{"hello", "world"},
			b:    []string{"hello", "world"},
			expected: "   1    1   hello\n" +
				"   2    2   world\n",
		},
		{
			name: "insert element",
			a:    []string{"hello", "world"},
			b:    []string{"hello", "middle", "world"},
			expected: "   1    1   hello\n" +
				"        2 + middle\n" +
				"   2    3   world\n",
		},
		{
			name: "delete element",
			a:    []string{"hello", "middle", "world"},
			b:    []string{"hello", "world"},
			expected: "   1    1   hello\n" +
				"   2      - middle\n" +
				"   3    2   world\n",
		},
		{
			name:     "empty slices",
			a:        []string{},
			b:        []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := simple.DiffStrings(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Error running DiffStrings: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestWithShowLineNumbers(t *testing.T) {
	a := "hello\nworld"
	b := "hello\neveryone"

	// With line numbers (default)
	withNumbers, err := simple.Diff(a, b)
	if err != nil {
		t.Fatalf("Error running Diff with line numbers: %v", err)
	}
	if !strings.Contains(withNumbers, "   1    1   hello") {
		t.Errorf("Expected line numbers, got: %s", withNumbers)
	}

	// Without line numbers
	withoutNumbers, err := simple.Diff(a, b, simple.WithShowLineNumbers(false))
	if err != nil {
		t.Fatalf("Error running Diff without line numbers: %v", err)
	}
	if strings.Contains(withoutNumbers, "   1    1") {
		t.Errorf("Did not expect line numbers, got: %s", withoutNumbers)
	}
	if !strings.Contains(withoutNumbers, "  hello") {
		t.Errorf("Expected content without line numbers, got: %s", withoutNumbers)
	}
}

func TestLongTextDiff(t *testing.T) {
	// Test with a large number of lines to ensure the algorithm handles it efficiently
	aLines := make([]string, 100)
	bLines := make([]string, 100)

	for i := range 100 {
		aLines[i] = "Line A " + string(rune(i%26+'a'))
		bLines[i] = "Line B " + string(rune(i%26+'a'))
	}

	// Make a few lines identical
	bLines[10] = aLines[10]
	bLines[50] = aLines[50]
	bLines[90] = aLines[90]

	// Run the diff
	_, err := simple.DiffStrings(aLines, bLines)
	if err != nil {
		t.Fatalf("Error running DiffStrings on long text: %v", err)
	}

	// Success if it completes without error
}

func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name string
		a    []string
		b    []string
	}{
		{
			name: "identical content",
			a:    []string{"same", "same", "same"},
			b:    []string{"same", "same", "same"},
		},
		{
			name: "completely different content",
			a:    []string{"only", "in", "a"},
			b:    []string{"only", "in", "b"},
		},
		{
			name: "one empty, one with content",
			a:    []string{},
			b:    []string{"content", "here"},
		},
		{
			name: "content with empty lines",
			a:    []string{"line1", "", "line3"},
			b:    []string{"line1", "line2", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := simple.DiffStrings(tc.a, tc.b)
			if err != nil {
				t.Fatalf("Error running DiffStrings on edge case: %v", err)
			}
		})
	}
}

func BenchmarkSimpleDiff(b *testing.B) {
	// Create two different texts to diff
	aLines := make([]string, 50)
	bLines := make([]string, 60)

	for i := range 50 {
		aLines[i] = "Line A " + string(rune(i%26+'a'))
	}

	for i := range 60 {
		bLines[i] = "Line B " + string(rune(i%26+'a'))
	}

	// Make some lines the same to create a realistic diff scenario
	for i := range 20 {
		pos := i * 2
		if pos < 50 && pos < 60 {
			bLines[pos] = aLines[pos]
		}
	}

	// Run the benchmark
	b.ResetTimer()
	for b.Loop() {
		_, err := simple.DiffStrings(aLines, bLines)
		if err != nil {
			b.Fatalf("Error in benchmark: %v", err)
		}
	}
}
