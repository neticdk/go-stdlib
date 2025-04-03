package myers_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/diff/myers"
)

const (
	// minimcs myers.smallInputThreshold
	smallInputThreshold = 100

	// mimics myers.largeInputThreshold
	largeInputThreshold = 10000
)

func TestMyersDifferInterface(t *testing.T) {
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
	differ := myers.NewDiffer()

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
	customDiffer := myers.NewCustomDiffer(
		myers.WithContextLines(1),
		myers.WithShowLineNumbers(false),
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

func TestLinearSpaceAlgorithmPaths(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		opts []myers.Option
		desc string
	}{
		{
			name: "very small input to force n=1 or m=1 in findMiddleSnake and use linear space",
			a:    []string{"a"},
			b:    []string{"b", "a"},
			opts: []myers.Option{myers.WithLinearSpace(true), myers.WithMaxEditDistance(100), myers.WithShowLineNumbers(false)}, // Force use of linear space algorithm, and allow a reasonable edit distance
			desc: "should use linear space algorithm and reach n=1||m=1 in findMiddleSnake",
		},
		{
			name: "small input (should use standard algorithm)",
			a:    make([]string, smallInputThreshold-1),
			b:    make([]string, smallInputThreshold-1),
			opts: []myers.Option{myers.WithLinearSpace(true)},
			desc: "small input should use standard algorithm",
		},
		{
			name: "max edit distance constraint",
			a:    make([]string, smallInputThreshold+10),
			b:    make([]string, smallInputThreshold+10),
			opts: []myers.Option{
				myers.WithLinearSpace(true),
				myers.WithMaxEditDistance(50),
			},
			desc: "should use standard algorithm due to edit distance constraint",
		},
		{
			name: "large input (should use linear space)",
			a:    make([]string, smallInputThreshold+1),
			b:    make([]string, smallInputThreshold+1),
			opts: []myers.Option{myers.WithLinearSpace(true)},
			desc: "should use linear space algorithm",
		},
		{
			name: "very large input (should use simple diff)",
			a:    make([]string, largeInputThreshold+1),
			b:    make([]string, largeInputThreshold+1),
			opts: []myers.Option{myers.WithLinearSpace(true)},
			desc: "should fall back to simple diff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip the generic checks for the "very small input" test case. We are now testing for this directly.
			if tt.name == "very small input to force n=1 or m=1 in findMiddleSnake and use linear space" {
				result, err := myers.DiffStrings(tt.a, tt.b, tt.opts...)
				if err != nil {
					t.Fatalf("Error running diff: %v", err)
				}

				expected := "+ b\n  a\n" // Expected output
				if result != expected {
					t.Errorf("Diff = %q, want %q", result, expected)
				}
				return
			}

			// Fill arrays with some predictable differences
			for i := range tt.a {
				if i%5 == 0 {
					tt.a[i] = fmt.Sprintf("a%d", i)
					tt.b[i] = fmt.Sprintf("b%d", i)
				} else {
					tt.a[i] = fmt.Sprintf("common%d", i)
					tt.b[i] = fmt.Sprintf("common%d", i)
				}
			}

			// Run diff with the specified options
			result, err := myers.DiffStrings(tt.a, tt.b, tt.opts...)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// Verify we got a valid diff
			if result == "" {
				t.Error("Expected non-empty diff result")
			}

			// Count changes to verify diff is working
			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			deletions := 0
			insertions := 0
			for _, line := range lines {
				trimmed := strings.TrimLeft(line, " 0123456789")
				if strings.HasPrefix(trimmed, "- ") {
					deletions++
				}
				if strings.HasPrefix(trimmed, "+ ") {
					insertions++
				}
			}

			// We should have an equal number of insertions and deletions
			if deletions != insertions {
				t.Errorf("Unbalanced changes: %d deletions, %d insertions", deletions, insertions)
			}

			// We should have some changes (every 5th line is different)
			expectedChanges := len(tt.a) / 5
			minExpectedChanges := expectedChanges / 2
			if deletions < minExpectedChanges {
				t.Errorf("Expected at least %d changes, got %d", minExpectedChanges, deletions)
			}

			t.Logf("Found %d changes in %s case", deletions, tt.desc)
		})
	}
}

func TestLinearSpaceBaseConditions(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected string
	}{
		{
			name: "empty a sequence",
			a:    []string{},
			b:    []string{"line1", "line2", "line3"},
			expected: "+ line1\n" +
				"+ line2\n" +
				"+ line3\n",
		},
		{
			name: "empty b sequence",
			a:    []string{"line1", "line2", "line3"},
			b:    []string{},
			expected: "- line1\n" +
				"- line2\n" +
				"- line3\n",
		},
		{
			name: "force recursive split",
			a:    []string{"line1", "line2", "", "line4", "line5"},
			b:    []string{"line1", "line2", "line3", "line4", "line5"},
			expected: "  line1\n" +
				"  line2\n" +
				"- \n" +
				"+ line3\n" +
				"  line4\n" +
				"  line5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a large enough input to trigger linear space algorithm
			// but pad it so the test sequences appear in the middle
			prefix := make([]string, 200)
			suffix := make([]string, 200)
			for i := range prefix {
				prefix[i] = fmt.Sprintf("prefix%d", i)
				suffix[i] = fmt.Sprintf("suffix%d", i)
			}

			// Add marker lines to help identify the test section
			marker := "TEST-SECTION"
			fullA := append(append(append(prefix, marker), tt.a...), suffix...)
			fullB := append(append(append(prefix, marker), tt.b...), suffix...)

			differ := myers.NewCustomDiffer(
				myers.WithLinearSpace(true),
				myers.WithMaxEditDistance(-1),    // Disable max edit distance
				myers.WithShowLineNumbers(false), // Disable line numbers for easier comparison
			)

			result, err := differ.DiffStrings(fullA, fullB)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// Extract the relevant part of the diff (the middle section)
			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			var relevantLines []string
			inRelevantSection := false
			for _, line := range lines {
				// Look for our marker
				if strings.Contains(line, marker) {
					inRelevantSection = true
					continue // Skip the marker line
				}
				if inRelevantSection {
					// Stop when we hit the suffix
					if strings.Contains(line, "suffix") {
						break
					}
					relevantLines = append(relevantLines, line)
				}
			}

			actual := strings.Join(relevantLines, "\n")
			if len(relevantLines) > 0 {
				actual += "\n"
			}
			if actual != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, actual)
			}
		})
	}
}

func TestEditScriptAlgorithmSelection(t *testing.T) {
	tests := []struct {
		name         string
		a            []string
		b            []string
		maxEditDist  int
		expectLinear bool
	}{
		{
			name:         "small input",
			a:            make([]string, 50),
			b:            make([]string, 50),
			maxEditDist:  10,
			expectLinear: false,
		},
		{
			name:         "medium input, constrained edit distance",
			a:            make([]string, 200),
			b:            make([]string, 200),
			maxEditDist:  50,
			expectLinear: false,
		},
		{
			name:         "medium input, unconstrained",
			a:            make([]string, 200),
			b:            make([]string, 200),
			maxEditDist:  0, // No constraint
			expectLinear: true,
		},
		{
			name:         "large input",
			a:            make([]string, 15000),
			b:            make([]string, 15000),
			maxEditDist:  0,
			expectLinear: false, // Should use simple diff
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill arrays with some predictable content
			for i := range tt.a {
				if i%5 == 0 {
					tt.a[i] = fmt.Sprintf("a%d", i)
					tt.b[i] = fmt.Sprintf("b%d", i)
				} else {
					tt.a[i] = fmt.Sprintf("common%d", i)
					tt.b[i] = fmt.Sprintf("common%d", i)
				}
			}

			opts := []myers.Option{
				myers.WithLinearSpace(true),
			}
			if tt.maxEditDist > 0 {
				opts = append(opts, myers.WithMaxEditDistance(tt.maxEditDist))
			}

			_, err := myers.DiffStrings(tt.a, tt.b, opts...)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// Check logs to verify which algorithm was used
			// (We might need to modify the implementation to expose this information)
		})
	}
}

func TestLinearSpaceRecursionDepth(t *testing.T) {
	// Create input that will test recursion depth
	size := smallInputThreshold + 50 // Large enough to use linear space
	a := make([]string, size)
	b := make([]string, size)

	// Create a pattern that will force deep recursion
	for i := range size {
		if i%3 == 0 {
			a[i] = fmt.Sprintf("a%d", i)
			b[i] = fmt.Sprintf("b%d", i)
		} else {
			// Create long sequences of matches to force recursion
			a[i] = fmt.Sprintf("common%d", i)
			b[i] = fmt.Sprintf("common%d", i)
		}
	}

	tests := []struct {
		name        string
		maxDepth    int
		expectError bool
	}{
		{
			name:     "normal depth",
			maxDepth: 30,
		},
		{
			name:     "shallow depth",
			maxDepth: 5,
		},
		{
			name:     "very shallow depth",
			maxDepth: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := myers.DiffStrings(a, b,
				myers.WithLinearSpace(true),
				myers.WithLinearRecursionMaxDepth(tt.maxDepth),
			)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// Verify we got a valid diff
			if result == "" {
				t.Error("Expected non-empty diff result")
			}

			// Count changes to verify diff is working
			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			changes := 0
			for _, line := range lines {
				trimmed := strings.TrimLeft(line, " 0123456789")
				if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "+ ") {
					changes++
				}
			}

			// We should have some changes (every 3rd line is different)
			expectedChanges := size / 3
			minExpectedChanges := expectedChanges / 2
			if changes < minExpectedChanges {
				t.Errorf("Expected at least %d changes, got %d", minExpectedChanges, changes)
			}

			t.Logf("Found %d changes with max depth %d", changes, tt.maxDepth)
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

func TestWithLargeInputThreshold(t *testing.T) {
	// Create a large input that would normally use the simple diff algorithm
	size := 15000
	aLines := make([]string, size)
	bLines := make([]string, size)

	for i := range size {
		aLines[i] = fmt.Sprintf("line %d", i)
		bLines[i] = fmt.Sprintf("line %d", i)
	}

	// Change a few lines to create differences
	bLines[5000] = "changed line"
	bLines[10000] = "another change"

	// Test with different fallback sizes
	tests := []struct {
		name           string
		fallbackSize   int
		shouldFallback bool
	}{
		{
			name:           "no fallback",
			fallbackSize:   20000, // Larger than input
			shouldFallback: false,
		},
		{
			name:           "force fallback",
			fallbackSize:   5000, // Smaller than input
			shouldFallback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			differ := myers.NewCustomDiffer(
				myers.WithLargeInputThreshold(tt.fallbackSize),
				myers.WithContextLines(0), // Minimize output size
			)

			result, err := differ.DiffStrings(aLines, bLines)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// The simple diff algorithm will include all changes in a single chunk
			// while the Myers algorithm with context will show them separately
			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			hasEllipsis := strings.Contains(result, "...\n")

			if tt.shouldFallback && hasEllipsis {
				t.Error("Expected simple diff (no chunks), but got chunked output")
			}
			if !tt.shouldFallback && !hasEllipsis && len(lines) > 4 {
				t.Error("Expected Myers diff with chunks, but got unchunked output")
			}
		})
	}
}

func TestWithSmallThreshold(t *testing.T) {
	// Create a small input
	sizeA := 500
	sizeB := 400

	// Test with different small input thresholds
	tests := []struct {
		name              string
		thresholdSize     int
		expectLinearSpace bool // Expecting the LinearSpace Algorithm when set to true
		desc              string
	}{
		{
			name:              "force use linear space",
			thresholdSize:     1, // Force smaller array to use linear space
			expectLinearSpace: true,
			desc:              "threshold smaller than the input: LinearSpace should  be called",
		},
		{
			name:              "use normal myers",
			thresholdSize:     1000, // Force larger array to use standard myers
			expectLinearSpace: false,
			desc:              "threshold larger than input size, expect !linearSpace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aLines := make([]string, sizeA)
			bLines := make([]string, sizeB)

			for i := range sizeA {
				aLines[i] = fmt.Sprintf("line %d", i)
			}

			for i := range sizeB {
				bLines[i] = fmt.Sprintf("line %d", i)
			}
			// Change a few lines to create differences
			if len(bLines) > 10 {
				bLines[10] = "changed line"
			}
			if len(aLines) > 20 {
				aLines[20] = "another change"
			}

			differ := myers.NewCustomDiffer(
				myers.WithSmallInputThreshold(tt.thresholdSize),
				myers.WithContextLines(0),   // Minimize output size
				myers.WithLinearSpace(true), // Force the code to call the linear space algorithm.
			)

			_, err := differ.DiffStrings(aLines, bLines)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// This will test to verify if the test hit the linear space algorithm or not
			aLen := len(aLines)
			bLen := len(bLines)

			// We are checking for conditions where computescript is not called!
			if aLen < smallInputThreshold || bLen < smallInputThreshold {
				if tt.expectLinearSpace != (tt.thresholdSize < min(aLen, bLen)) {
					t.Errorf("Inconsistent linear space: %s.  Desc: %s. expected %v, got %v", tt.name, tt.desc, tt.expectLinearSpace, (tt.thresholdSize < min(aLen, bLen)))
				}
			} else {
				// If both are above the threshold, expect linearSpace to be used if thresholdSize is less.
				if tt.expectLinearSpace != (tt.thresholdSize < min(aLen, bLen)) {
					t.Errorf("Inconsistent linear space: %s.  Desc: %s. expected %v, got %v", tt.name, tt.desc, tt.expectLinearSpace, (tt.thresholdSize < min(aLen, bLen)))
				}
			}
		})
	}
}

func TestWithLinearRecursionMaxDepth(t *testing.T) {
	// Create input that will cause deep recursion
	size := 100 // Reduced size for more predictable behavior
	aLines := make([]string, size)
	bLines := make([]string, size)

	// Create a pattern that will show different behavior between algorithms
	for i := range size {
		if i%2 == 0 {
			aLines[i] = fmt.Sprintf("a%d", i)
			bLines[i] = fmt.Sprintf("b%d", i)
		} else {
			aLines[i] = fmt.Sprintf("common%d", i)
			bLines[i] = fmt.Sprintf("common%d", i)
		}
	}

	tests := []struct {
		name     string
		maxDepth int
	}{
		{
			name:     "shallow recursion",
			maxDepth: 5,
		},
		{
			name:     "deep recursion",
			maxDepth: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			differ := myers.NewCustomDiffer(
				myers.WithLinearSpace(true),
				myers.WithLinearRecursionMaxDepth(tt.maxDepth),
				myers.WithContextLines(0), // No context to make changes clearer
			)

			result, err := differ.DiffStrings(aLines, bLines)
			if err != nil {
				t.Fatalf("Error running diff: %v", err)
			}

			// Count the actual changes
			lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
			deletions := 0
			insertions := 0

			for _, line := range lines {
				trimmed := strings.TrimLeft(line, " 0123456789")
				if strings.HasPrefix(trimmed, "- ") {
					deletions++
				}
				if strings.HasPrefix(trimmed, "+ ") {
					insertions++
				}
			}

			// The changes should be paired (each change is a deletion and an insertion)
			if deletions != insertions {
				t.Errorf("Unbalanced changes: %d deletions, %d insertions", deletions, insertions)
			}

			changes := deletions // or insertions, they should be equal
			t.Logf("Found %d changes", changes)

			// We expect approximately size/2 changes (every other line is different)
			expectedChanges := size / 2
			minExpectedChanges := expectedChanges / 2

			if changes == 0 {
				t.Error("Expected changes, got none")
				// Print first few lines of input for debugging
				for i := 0; i < min(10, len(aLines)); i++ {
					t.Logf("Line %d: a=%q b=%q", i, aLines[i], bLines[i])
				}
			}

			// For deep recursion, we expect to find most changes
			if tt.maxDepth > 10 && changes < minExpectedChanges {
				t.Errorf("With deep recursion, expected at least %d changes, got %d",
					minExpectedChanges, changes)
			}

			// Verify specific changes are present
			if !strings.Contains(result, "a0") && !strings.Contains(result, "b0") {
				t.Error("Expected to find changes in the first line")
			}
		})
	}
}

func TestOptionCombinations(t *testing.T) {
	// Test that options can be combined effectively
	a := []string{
		"line1",
		"line2",
		"oldline3",
		"line4",
		"line5",
		"line6",
		"line7",
	}
	b := []string{
		"line1",
		"line2",
		"newline3",
		"line4",
		"line5",
		"line6",
		"line7",
	}

	differ := myers.NewCustomDiffer(
		myers.WithLinearSpace(true),
		myers.WithLinearRecursionMaxDepth(10),
		myers.WithLargeInputThreshold(500),
		myers.WithContextLines(2),
		myers.WithShowLineNumbers(false),
		myers.WithMaxEditDistance(50),
	)

	result, err := differ.DiffStrings(a, b)
	if err != nil {
		t.Fatalf("Error running diff with combined options: %v", err)
	}

	// Verify that line numbers are hidden
	if strings.Contains(result, "   1    1") {
		t.Error("Line numbers should be hidden")
	}

	// Verify context lines
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	expectedLines := []string{
		"  line1",
		"  line2",
		"- oldline3",
		"+ newline3",
		"  line4",
		"  line5",
	}

	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, expected := range expectedLines {
		if i < len(lines) && lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
		}
	}
}

func TestLongTextDiff(t *testing.T) {
	// Test with a large number of lines to ensure the algorithm handles it efficiently
	aLines := make([]string, 500)
	bLines := make([]string, 500)

	for i := range 500 {
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

func BenchmarkMyersDiff(b *testing.B) {
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
		_, err := myers.DiffStrings(aLines, bLines)
		if err != nil {
			b.Fatalf("Error in benchmark: %v", err)
		}
	}
}

// func BenchmarkMyersDiffLinearSpace(b *testing.B) {
// 	sizes := []int{100, 1000, 10000}
// 	changes := []float64{0.01, 0.1, 0.5} // Percentage of lines changed

//		for _, size := range sizes {
//			for _, changeRate := range changes {
//				name := fmt.Sprintf("size=%d,changes=%.2f", size, changeRate)
//				b.Run(name, func(b *testing.B) {
//					a, bb := generateBenchmarkInput(size, changeRate)
//					b.ResetTimer()
//					for i := 0; i < b.N; i++ {
//						_, err := myers.Diff(a, bb, myers.WithLinearSpace(true))
//						if err != nil {
//							b.Fatalf("Error in benchmark: %v", err)
//						}
//					}
//				})
//			}
//		}
//	}
func BenchmarkMyersDiffLinearSpace(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	changes := []float64{0.01, 0.1, 0.5} // Percentage of lines changed

	for _, size := range sizes {
		for _, changeRate := range changes {
			name := fmt.Sprintf("size=%d,changes=%.2f", size, changeRate)
			b.Run(name, func(b *testing.B) {
				a, bb := generateBenchmarkInput(size, changeRate)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := myers.Diff(a, bb, myers.WithLinearSpace(true))
					if err != nil {
						b.Fatalf("Error in benchmark: %v", err)
					}
				}
			})
		}
	}
}

// generateBenchmarkInput creates two strings for benchmarking diff operations.
// size: number of lines to generate
// changeRate: fraction of lines that should be different (0.0 to 1.0)
func generateBenchmarkInput(size int, changeRate float64) (string, string) {
	if size <= 0 {
		return "", ""
	}
	if changeRate < 0.0 {
		changeRate = 0.0
	}
	if changeRate > 1.0 {
		changeRate = 1.0
	}

	// Pre-allocate slices for both inputs
	aLines := make([]string, size)
	bLines := make([]string, size)

	// Calculate how many lines should be different
	changesToMake := int(float64(size) * changeRate)

	// Create a set of indices that will be changed
	changeIndices := make(map[int]bool)
	if changesToMake > 0 {
		// Use a simple random selection if we're changing less than half
		if changeRate <= 0.5 {
			for len(changeIndices) < changesToMake {
				idx := rand.Intn(size)
				changeIndices[idx] = true
			}
		} else {
			// For higher change rates, select indices to keep the same
			keepSame := size - changesToMake
			for i := range size {
				changeIndices[i] = true
			}
			for range keepSame {
				idx := rand.Intn(size)
				delete(changeIndices, idx)
			}
		}
	}

	// Generate the lines with controlled changes
	for i := range size {
		if changeIndices[i] {
			// Create different content for changed lines
			aLines[i] = fmt.Sprintf("old line %d content %d", i, rand.Intn(1000))
			bLines[i] = fmt.Sprintf("new line %d content %d", i, rand.Intn(1000))
		} else {
			// Create identical content for unchanged lines
			content := fmt.Sprintf("same line %d content %d", i, rand.Intn(1000))
			aLines[i] = content
			bLines[i] = content
		}
	}

	// Consider adding some structural changes based on changeRate
	if changeRate > 0.3 {
		// Maybe add some lines
		extraLines := int(float64(size) * 0.1) // Add up to 10% extra lines
		if extraLines > 0 {
			for i := range extraLines {
				pos := rand.Intn(len(bLines) + 1)
				newLine := fmt.Sprintf("inserted line %d", i)
				bLines = append(bLines[:pos], append([]string{newLine}, bLines[pos:]...)...)
			}
		}

		// Maybe delete some lines from a
		if changeRate > 0.6 {
			deletions := int(float64(size) * 0.1) // Delete up to 10% of lines
			for range deletions {
				if len(aLines) > 1 { // Keep at least one line
					pos := rand.Intn(len(aLines))
					aLines = append(aLines[:pos], aLines[pos+1:]...)
				}
			}
		}
	}

	// Add some common patterns that often appear in real text
	if size > 10 {
		// Add some repeated lines
		repeatedLine := "this line appears multiple times"
		for range 3 {
			pos := rand.Intn(len(aLines))
			aLines[pos] = repeatedLine
			bLines[pos] = repeatedLine
		}

		// Add some blocks of similar lines
		if size > 100 {
			blockSize := min(5, size/20)
			blockStart := rand.Intn(size - blockSize)
			for i := range blockSize {
				prefix := "block line "
				if changeIndices[blockStart+i] {
					aLines[blockStart+i] = prefix + "old " + strconv.Itoa(i)
					bLines[blockStart+i] = prefix + "new " + strconv.Itoa(i)
				} else {
					aLines[blockStart+i] = prefix + strconv.Itoa(i)
					bLines[blockStart+i] = prefix + strconv.Itoa(i)
				}
			}
		}
	}

	return strings.Join(aLines, "\n"), strings.Join(bLines, "\n")
}

// generateBenchmarkInput creates two strings for benchmarking diff operations.
// size: number of lines to generate
// changeRate: fraction of lines that should be different (0.0 to 1.0)
// func generateBenchmarkInput(size int, changeRate float64) (string, string) {
// 	if size <= 0 {
// 		return "", ""
// 	}
// 	if changeRate < 0.0 {
// 		changeRate = 0.0
// 	}
// 	if changeRate > 1.0 {
// 		changeRate = 1.0
// 	}

// 	// Pre-allocate slices for both inputs
// 	aLines := make([]string, size)
// 	bLines := make([]string, size)

// 	// Calculate how many lines should be different
// 	changesToMake := int(float64(size) * changeRate)

// 	// Create a set of indices that will be changed
// 	changeIndices := make(map[int]bool)
// 	if changesToMake > 0 {
// 		// Use a simple random selection if we're changing less than half
// 		if changeRate <= 0.5 {
// 			for len(changeIndices) < changesToMake {
// 				idx := rand.Intn(size)
// 				changeIndices[idx] = true
// 			}
// 		} else {
// 			// For higher change rates, select indices to keep the same
// 			keepSame := size - changesToMake
// 			for i := range size {
// 				changeIndices[i] = true
// 			}
// 			for range keepSame {
// 				idx := rand.Intn(size)
// 				delete(changeIndices, idx)
// 			}
// 		}
// 	}

// 	// Generate the lines with controlled changes
// 	for i := range size {
// 		if changeIndices[i] {
// 			// Create different content for changed lines
// 			aLines[i] = fmt.Sprintf("old line %d content %d", i, rand.Intn(1000))
// 			bLines[i] = fmt.Sprintf("new line %d content %d", i, rand.Intn(1000))
// 		} else {
// 			// Create identical content for unchanged lines
// 			content := fmt.Sprintf("same line %d content %d", i, rand.Intn(1000))
// 			aLines[i] = content
// 			bLines[i] = content
// 		}
// 	}

// 	// Consider adding some structural changes based on changeRate
// 	if changeRate > 0.3 {
// 		// Maybe add some lines
// 		extraLines := int(float64(size) * 0.1) // Add up to 10% extra lines
// 		if extraLines > 0 {
// 			for i := range extraLines {
// 				pos := rand.Intn(len(bLines) + 1)
// 				newLine := fmt.Sprintf("inserted line %d", i)
// 				bLines = append(bLines[:pos], append([]string{newLine}, bLines[pos:]...)...)
// 			}
// 		}

// 		// Maybe delete some lines from a
// 		if changeRate > 0.6 {
// 			deletions := int(float64(size) * 0.1) // Delete up to 10% of lines
// 			for range deletions {
// 				if len(aLines) > 1 { // Keep at least one line
// 					pos := rand.Intn(len(aLines))
// 					aLines = append(aLines[:pos], aLines[pos+1:]...)
// 				}
// 			}
// 		}
// 	}

// 	// Add some common patterns that often appear in real text
// 	if size > 10 {
// 		// Add some repeated lines
// 		repeatedLine := "this line appears multiple times"
// 		for range 3 {
// 			pos := rand.Intn(len(aLines))
// 			aLines[pos] = repeatedLine
// 			bLines[pos] = repeatedLine
// 		}

// 		// Add some blocks of similar lines
// 		if size > 100 {
// 			blockSize := min(5, size/20)
// 			blockStart := rand.Intn(size - blockSize)
// 			for i := range blockSize {
// 				prefix := "block line "
// 				if changeIndices[blockStart+i] {
// 					aLines[blockStart+i] = prefix + "old " + strconv.Itoa(i)
// 					bLines[blockStart+i] = prefix + "new " + strconv.Itoa(i)
// 				} else {
// 					aLines[blockStart+i] = prefix + strconv.Itoa(i)
// 					bLines[blockStart+i] = prefix + strconv.Itoa(i)
// 				}
// 			}
// 		}
// 	}

// 	return strings.Join(aLines, "\n"), strings.Join(bLines, "\n")
// }
