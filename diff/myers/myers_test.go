package myers_test

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/diff/myers"
)

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
			result, err := myers.Diff(tt.a, tt.b)
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
			result, err := myers.DiffStrings(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Error running DiffStrings: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestWithContextLines(t *testing.T) {
	a := "line1\nline2\nline3\nold\nline5\nline6\nline7\nline8"
	b := "line1\nline2\nline3\nnew\nline5\nline6\nline7\nline8"

	tests := []struct {
		name         string
		contextLines int
		expectLines  int
	}{
		{"no context", 0, 2},      // Just the changed line
		{"default context", 3, 8}, // 3 lines before, changed line, 3 lines after
		{"full context", 10, 9},   // All lines (since file is only 9 lines)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := myers.Diff(a, b, myers.WithContextLines(tt.contextLines))
			if err != nil {
				t.Fatalf("Error running Diff: %v", err)
			}

			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			if len(lines) != tt.expectLines {
				t.Errorf("Expected %d lines, got %d lines: %s", tt.expectLines, len(lines), result)
			}
		})
	}
}

func TestWithShowLineNumbers(t *testing.T) {
	a := "hello\nworld"
	b := "hello\neveryone"

	// With line numbers (default)
	withNumbers, err := myers.Diff(a, b)
	if err != nil {
		t.Fatalf("Error running Diff with line numbers: %v", err)
	}
	if !strings.Contains(withNumbers, "   1    1   hello") {
		t.Errorf("Expected line numbers, got: %s", withNumbers)
	}

	// Without line numbers
	withoutNumbers, err := myers.Diff(a, b, myers.WithShowLineNumbers(false))
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

func TestWithMaxEditDistance(t *testing.T) {
	// Create strings with significant differences
	a := strings.Repeat("a\n", 100)
	b := strings.Repeat("b\n", 100)

	// With unlimited edit distance
	_, err := myers.Diff(a, b)
	if err != nil {
		t.Fatalf("Error running Diff with default max edit distance: %v", err)
	}

	// With very limited edit distance
	_, err = myers.Diff(a, b, myers.WithMaxEditDistance(5))
	if err != nil {
		t.Fatalf("Error running Diff with limited max edit distance: %v", err)
	}

	// The test is successful if both calls complete without errors
	// The limited edit distance version should fall back to the simple diff method
}

func TestLongTextDiff(t *testing.T) {
	// Test with a large number of lines to ensure the algorithm handles it efficiently
	aLines := make([]string, 500)
	bLines := make([]string, 500)

	for i := 0; i < 500; i++ {
		aLines[i] = "Line A " + string(rune(i%26+'a'))
		bLines[i] = "Line B " + string(rune(i%26+'a'))
	}

	// Introduce a few changes
	bLines[100] = aLines[100]
	bLines[200] = aLines[200]
	bLines[300] = aLines[300]

	// With context
	_, err := myers.DiffStrings(aLines, bLines, myers.WithContextLines(3))
	if err != nil {
		t.Fatalf("Error running DiffStrings on long text: %v", err)
	}

	// Without context
	_, err = myers.DiffStrings(aLines, bLines, myers.WithContextLines(0))
	if err != nil {
		t.Fatalf("Error running DiffStrings on long text without context: %v", err)
	}
}

func TestCombinedOptions(t *testing.T) {
	a := "line1\nline2\nline3\nold\nline5"
	b := "line1\nline2\nline3\nnew\nline5"

	// Combine multiple options
	result, err := myers.Diff(a, b,
		myers.WithContextLines(1),
		myers.WithShowLineNumbers(false),
		myers.WithMaxEditDistance(50),
		myers.WithLinearSpace(true),
	)
	if err != nil {
		t.Fatalf("Error running Diff with combined options: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	// With context 1, we should have 4 lines:
	// line3 (context before)
	// - old (deletion)
	// + new (addition)
	// line5 (context after)
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d lines: %s", len(lines), result)
	}

	// Without line numbers, lines should start with spaces, +, or -
	for _, line := range lines {
		if !strings.HasPrefix(line, "  ") &&
			!strings.HasPrefix(line, "+ ") &&
			!strings.HasPrefix(line, "- ") {
			t.Errorf("Expected line to start with '  ', '+ ', or '- ', got: %s", line)
		}
	}
}
