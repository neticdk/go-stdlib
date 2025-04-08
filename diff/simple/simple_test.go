package simple_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/diff"
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
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected, "Diff/%q", tt.name)
		})
	}

	// Test string slice input
	for _, tt := range sliceTests {
		t.Run("Default/Slice/"+tt.name, func(t *testing.T) {
			result, err := differ.DiffStrings(tt.a, tt.b)
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected, "DiffStrings/%q", tt.name)
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
		assert.NoError(t, err)
		assert.Equal(t, result, customExpected)
	})

	t.Run("Custom/Slice", func(t *testing.T) {
		result, err := customDiffer.DiffStrings(
			[]string{"hello", "world"},
			[]string{"hello", "middle", "world"},
		)
		assert.NoError(t, err)
		assert.Equal(t, result, customExpected)
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
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected, "Diff/%q", tt.name)
		})
	}
}

// TestOptionsValidationViaDiffStrings tests the internal validation logic indirectly
// by calling DiffStrings and checking the returned error.
func TestOptionsValidationViaDiffStrings(t *testing.T) {
	// Dummy inputs, content doesn't matter for option validation
	dummyA := []string{"a"}
	dummyB := []string{"b"}

	tests := []struct {
		name         string
		opts         []simple.Option
		expectErr    bool
		errSubstring string // Substring expected in the error message
	}{
		{
			name:      "default options are valid",
			opts:      []simple.Option{},
			expectErr: false,
		},
		{
			name: "valid custom options",
			opts: []simple.Option{
				simple.WithContextLines(5),
				simple.WithShowLineNumbers(false),
				simple.WithUnifiedFormatter(), // Use With... option to set formatter
			},
			expectErr: false,
		},
		// Test validation inherited from FormatOptions
		{
			name: "invalid context lines (inherited)",
			opts: []simple.Option{
				simple.WithContextLines(-1),
			},
			expectErr:    true,
			errSubstring: "ContextLines must be non-negative",
		},
		// Note: Testing invalid OutputFormat enum value is difficult via functional options,
		// but the FormatOptions.Validate test covers the internal check.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := simple.DiffStrings(dummyA, dummyB, tt.opts...)

			if tt.expectErr {
				assert.Error(t, err, "Expected an error but got none")
				if err != nil && tt.errSubstring != "" {
					assert.Contains(t, err.Error(), tt.errSubstring, "Error message mismatch")
				}
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
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
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected, "DiffStrings/%q", tt.name)
		})
	}
}

func TestWithShowLineNumbers(t *testing.T) {
	a := "hello\nworld"
	b := "hello\neveryone"

	// With line numbers (default)
	withNumbers, err := simple.Diff(a, b)
	assert.NoError(t, err)
	assert.Contains(t, withNumbers, "   1    1", "Expected line numbers")

	// Without line numbers
	withoutNumbers, err := simple.Diff(a, b, simple.WithShowLineNumbers(false))
	assert.NoError(t, err)
	assert.NotContains(t, withoutNumbers, "   1    1", "Did not expect line numbers")
	assert.NotContains(t, withoutNumbers, "   2    2", "Did not expect line numbers")
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
	result, err := simple.DiffStrings(aLines, bLines)
	assert.NoError(t, err)
	assert.Contains(t, result, "Line A a", "Expected content with line numbers")
}

func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		a        []string
		b        []string
		expected string
	}{
		{
			name: "identical content",
			a:    []string{"same", "same", "same"},
			b:    []string{"same", "same", "same"},
			expected: "   1    1   same\n" +
				"   2    2   same\n" +
				"   3    3   same\n",
		},
		{
			name: "completely different content",
			a:    []string{"only", "in", "a"},
			b:    []string{"only", "in", "b"},
			expected: "   1    1   only\n" +
				"   2    2   in\n" +
				"   3      - a\n" +
				"        3 + b\n",
		},
		{
			name: "one empty, one with content",
			a:    []string{},
			b:    []string{"content", "here"},
			expected: "        1 + content\n" +
				"        2 + here\n",
		},
		{
			name: "content with empty lines",
			a:    []string{"line1", "", "line3"},
			b:    []string{"line1", "line2", ""},
			expected: "   1    1   line1\n" +
				"        2 + line2\n" +
				"   2    3   \n" +
				"   3      - line3\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := simple.DiffStrings(tc.a, tc.b)
			assert.NoError(t, err)
			assert.Equal(t, result, tc.expected, "DiffStrings/%q", tc.name)
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
		_, _ = simple.DiffStrings(aLines, bLines)
	}
}

// mockFormatter for testing WithFormatter option (can be shared or redefined)
type mockFormatter struct {
	returnValue string
}

func (m mockFormatter) Format(edits []diff.Line, options diff.FormatOptions) string {
	return m.returnValue
}

// capturingMockFormatter captures the options passed to it
type capturingMockFormatter struct {
	capturedOptions diff.FormatOptions
}

func (m *capturingMockFormatter) Format(edits []diff.Line, options diff.FormatOptions) string {
	m.capturedOptions = options
	// Return a string representing the captured options for verification
	return fmt.Sprintf("Formatted with options: {OutputFormat:%s ContextLines:%d ShowLineNumbers:%v}",
		options.OutputFormat, options.ContextLines, options.ShowLineNumbers)
}

func TestWithFormatterOptions(t *testing.T) {
	a := "line1\nold\nline3"
	b := "line1\nnew\nline3"

	tests := []struct {
		name     string
		options  []simple.Option
		expected string
	}{
		{
			name: "Default (ContextFormatter)",
			options: []simple.Option{
				simple.WithShowLineNumbers(false), // Simplify expected output
			},
			expected: "  line1\n- old\n+ new\n  line3\n",
		},
		{
			name: "WithContextFormatter",
			options: []simple.Option{
				simple.WithContextFormatter(),
				simple.WithShowLineNumbers(false),
			},
			expected: "  line1\n- old\n+ new\n  line3\n",
		},
		{
			name: "WithUnifiedFormatter",
			options: []simple.Option{
				simple.WithUnifiedFormatter(),
				simple.WithShowLineNumbers(false),
			},
			expected: "--- a\n+++ b\n@@ -1,3 +1,3 @@\n line1\n-old\n+new\n line3\n",
		},
		{
			name: "WithCustomFormatter",
			options: []simple.Option{
				simple.WithFormatter(mockFormatter{returnValue: "Mock simple output!"}),
			},
			expected: "Mock simple output!",
		},
		{
			name: "WithCustomFormatter nil (should use default)",
			options: []simple.Option{
				simple.WithFormatter(nil), // Setting nil should be ignored, default ContextFormatter used
				simple.WithShowLineNumbers(false),
			},
			expected: "  line1\n- old\n+ new\n  line3\n",
		},
		{
			name: "WithOutputFormat (used by formatter)",
			options: []simple.Option{
				simple.WithFormatter(&capturingMockFormatter{}),
				simple.WithOutputFormat(diff.FormatUnified), // Set the format
				simple.WithContextLines(5),                  // Set other format options
				simple.WithShowLineNumbers(false),
			},
			// Expected value comes from the mock formatter which should echo the option
			expected: "Formatted with options: {OutputFormat:unified ContextLines:5 ShowLineNumbers:false}",
		},
		{
			name: "Formatter with specific context options",
			options: []simple.Option{
				simple.WithContextFormatter(), // Explicitly context
				simple.WithContextLines(0),    // No context lines
				simple.WithShowLineNumbers(false),
			},
			expected: "- old\n+ new\n", // Only changed lines
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using DiffStrings directly with options
			result, err := simple.DiffStrings(strings.Split(a, "\n"), strings.Split(b, "\n"), tt.options...)
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected, "DiffStrings/%q", tt.name)

			// Also test via CustomDiffer factory
			customDiffer := simple.NewCustomDiffer(tt.options...)
			resultFactory, err := customDiffer.Diff(a, b)
			assert.NoError(t, err)
			assert.Equal(t, resultFactory, tt.expected, "CustomDiffer.Diff/%q", tt.name)
		})
	}
}
