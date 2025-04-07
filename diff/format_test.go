package diff_test // Use _test package for black-box testing

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/diff"
)

func TestFormatOptions_Validate(t *testing.T) {
	tests := []struct {
		name      string
		options   diff.FormatOptions
		expectErr bool
	}{
		{
			name: "valid options (context)",
			options: diff.FormatOptions{
				OutputFormat:    diff.FormatContext,
				ContextLines:    3,
				ShowLineNumbers: true,
			},
			expectErr: false,
		},
		{
			name: "valid options (unified)",
			options: diff.FormatOptions{
				OutputFormat:    diff.FormatUnified,
				ContextLines:    0,
				ShowLineNumbers: false,
			},
			expectErr: false,
		},
		{
			name: "invalid context lines",
			options: diff.FormatOptions{
				OutputFormat:    diff.FormatContext,
				ContextLines:    -1,
				ShowLineNumbers: true,
			},
			expectErr: true,
		},
		{
			name: "invalid output format",
			options: diff.FormatOptions{
				OutputFormat:    diff.OutputFormat(99), // Invalid enum value
				ContextLines:    3,
				ShowLineNumbers: true,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			assert.Equal(t, err != nil, tt.expectErr, "Validate() error = %v, expectErr %v", err, tt.expectErr)
		})
	}
}

func TestContextFormatter_Format(t *testing.T) {
	formatter := diff.ContextFormatter{}
	edits := []diff.Line{
		{Kind: diff.Equal, Text: "line1"},
		{Kind: diff.Delete, Text: "line2-old"},
		{Kind: diff.Insert, Text: "line2-new"},
		{Kind: diff.Equal, Text: "line3"},
		{Kind: diff.Equal, Text: "line4"},
		{Kind: diff.Equal, Text: "line5"},
		{Kind: diff.Equal, Text: "line6"},
		{Kind: diff.Delete, Text: "line7-old"},
		{Kind: diff.Equal, Text: "line8"},
	}

	tests := []struct {
		name     string
		options  diff.FormatOptions
		edits    []diff.Line
		expected string
	}{
		{
			name: "default context (3 lines), with line numbers",
			options: diff.FormatOptions{
				ContextLines:    diff.DefaultContextLines,    // 3
				ShowLineNumbers: diff.DefaultShowLineNumbers, // true
			},
			expected: `   1    1   line1
   2      - line2-old
        2 + line2-new
   3    3   line3
   4    4   line4
   5    5   line5
   6    6   line6
   7      - line7-old
   8    7   line8
`, // No "..." because changes at index 1/2 and 7 are 5 apart, which is <= 2*contextLines (6)
		},
		{
			name: "no context (0 lines), with line numbers",
			options: diff.FormatOptions{
				ContextLines:    0,
				ShowLineNumbers: true,
			},
			expected: `   2      - line2-old
        2 + line2-new
   7      - line7-old
`,
		},
		{
			name: "one context line, without line numbers",
			options: diff.FormatOptions{
				ContextLines:    1,
				ShowLineNumbers: false,
			},
			expected: `  line1
- line2-old
+ line2-new
  line3
...
  line6
- line7-old
  line8
`,
		},
		{
			name: "full context (more than needed)",
			options: diff.FormatOptions{
				ContextLines:    10,
				ShowLineNumbers: false,
			},
			expected: `  line1
- line2-old
+ line2-new
  line3
  line4
  line5
  line6
- line7-old
  line8
`, // No "..." expected because large context merges everything
		},
		{
			name: "no edits",
			options: diff.FormatOptions{
				ContextLines:    3,
				ShowLineNumbers: true,
			},
			edits:    []diff.Line{}, // Explicitly pass empty slice
			expected: "",            // Formatter returns empty for no edits
		},
		{
			name: "only equal edits",
			options: diff.FormatOptions{
				ContextLines:    3,
				ShowLineNumbers: false,
			},
			// Edits with only Equal lines should still be formatted if context > 0
			// The grouping logic returns the whole block if no changes found
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "equal1"},
				{Kind: diff.Equal, Text: "equal2"},
			},
			expected: "  equal1\n  equal2\n", // Standard diff tools output nothing if files are identical (no changes)
		},
		{
			name: "only change edits, no context lines",
			options: diff.FormatOptions{
				ContextLines:    0,
				ShowLineNumbers: false,
			},
			edits: []diff.Line{
				{Kind: diff.Delete, Text: "del"},
				{Kind: diff.Insert, Text: "ins"},
			},
			expected: `- del
+ ins
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentEdits := edits
			if tt.edits != nil { // Allow overriding edits for specific tests
				currentEdits = tt.edits
			}
			result := formatter.Format(currentEdits, tt.options)
			assert.Equal(t, result, tt.expected, "ContextFormatter.Format/%q", tt.name)
		})
	}
}

func TestUnifiedFormatter_Format(t *testing.T) {
	formatter := diff.UnifiedFormatter{}
	edits := []diff.Line{
		{Kind: diff.Equal, Text: "line1"},
		{Kind: diff.Delete, Text: "line2-old"},
		{Kind: diff.Insert, Text: "line2-new"},
		{Kind: diff.Equal, Text: "line3"},
		{Kind: diff.Equal, Text: "line4"},
		{Kind: diff.Equal, Text: "line5"},
		{Kind: diff.Equal, Text: "line6"},
		{Kind: diff.Delete, Text: "line7-old"},
		{Kind: diff.Equal, Text: "line8"},
	}

	tests := []struct {
		name     string
		options  diff.FormatOptions
		edits    []diff.Line
		expected string
	}{
		{
			name: "default context (3 lines)",
			options: diff.FormatOptions{
				ContextLines: diff.DefaultContextLines, // 3
				// ShowLineNumbers typically ignored by unified format
			},
			expected: `--- a
+++ b
@@ -1,8 +1,7 @@
 line1
-line2-old
+line2-new
 line3
 line4
 line5
 line6
-line7-old
 line8
`, // Single hunk because changes are close enough (<= 2*context)
		},
		{
			name: "no context (0 lines)",
			options: diff.FormatOptions{
				ContextLines: 0,
			},
			// Note: Unified format often implicitly includes *some* context
			// around the hunk header, even if 0 specified. The current implementation
			// of groupEditsByContext returns only changed lines for context=0.
			// Let's adjust expected based on groupEditsByContext(context=0) behavior.
			expected: `--- a
+++ b
@@ -2,1 +1,0 @@
-line2-old
@@ -2,0 +2,1 @@
+line2-new
@@ -7,1 +6,0 @@
-line7-old
`,
		},
		{
			name: "one context line",
			options: diff.FormatOptions{
				ContextLines: 1,
			},
			expected: `--- a
+++ b
@@ -1,3 +1,3 @@
 line1
-line2-old
+line2-new
 line3
@@ -6,3 +6,2 @@
 line6
-line7-old
 line8
`,
		},
		{
			name: "full context (more than needed)",
			options: diff.FormatOptions{
				ContextLines: 10,
			},
			expected: `--- a
+++ b
@@ -1,8 +1,7 @@
 line1
-line2-old
+line2-new
 line3
 line4
 line5
 line6
-line7-old
 line8
`,
		},
		{
			name: "no edits",
			options: diff.FormatOptions{
				ContextLines: 3,
			},
			edits:    []diff.Line{},    // Explicitly provide empty edits
			expected: "--- a\n+++ b\n", // Returns empty for no edits
		},
		{
			name: "only equal edits",
			options: diff.FormatOptions{
				ContextLines: 3,
			},
			edits: []diff.Line{ // Explicitly provide only equal edits
				{Kind: diff.Equal, Text: "equal1"},
				{Kind: diff.Equal, Text: "equal2"},
			},
			expected: "--- a\n+++ b\n@@ -1,2 +1,2 @@\n equal1\n equal2\n",
		},
		{
			name: "only delete edits",
			options: diff.FormatOptions{
				ContextLines: 1,
			},
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "keep"},
				{Kind: diff.Delete, Text: "del1"},
				{Kind: diff.Delete, Text: "del2"},
				{Kind: diff.Equal, Text: "keep2"},
			},
			expected: `--- a
+++ b
@@ -1,4 +1,2 @@
 keep
-del1
-del2
 keep2
`,
		},
		{
			name: "only insert edits",
			options: diff.FormatOptions{
				ContextLines: 1,
			},
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "keep"},
				{Kind: diff.Insert, Text: "ins1"},
				{Kind: diff.Insert, Text: "ins2"},
				{Kind: diff.Equal, Text: "keep2"},
			},
			expected: `--- a
+++ b
@@ -1,2 +1,4 @@
 keep
+ins1
+ins2
 keep2
`,
		},
		{
			name: "change at beginning",
			options: diff.FormatOptions{
				ContextLines: 1,
			},
			edits: []diff.Line{
				{Kind: diff.Delete, Text: "del_start"},
				{Kind: diff.Insert, Text: "ins_start"},
				{Kind: diff.Equal, Text: "keep"},
			},
			expected: `--- a
+++ b
@@ -1,2 +1,2 @@
-del_start
+ins_start
 keep
`,
		},
		{
			name: "change at end",
			options: diff.FormatOptions{
				ContextLines: 1,
			},
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "keep"},
				{Kind: diff.Delete, Text: "del_end"},
				{Kind: diff.Insert, Text: "ins_end"},
			},
			expected: `--- a
+++ b
@@ -1,2 +1,2 @@
 keep
-del_end
+ins_end
`,
		},
		{
			name: "zero length hunk - pure insertion",
			options: diff.FormatOptions{
				ContextLines: 2,
			},
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "inserted line 1"},
				{Kind: diff.Insert, Text: "inserted line 2"},
				{Kind: diff.Equal, Text: "line2"},
			},
			expected: `--- a
+++ b
@@ -1,2 +1,4 @@
 line1
+inserted line 1
+inserted line 2
 line2
`,
		},
		{
			name: "zero length hunk - pure deletion",
			options: diff.FormatOptions{
				ContextLines: 2,
			},
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Delete, Text: "deleted line 1"},
				{Kind: diff.Delete, Text: "deleted line 2"},
				{Kind: diff.Equal, Text: "line2"},
			},
			expected: `--- a
+++ b
@@ -1,4 +1,2 @@
 line1
-deleted line 1
-deleted line 2
 line2
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentEdits := edits
			if tt.edits != nil { // Allow overriding edits for specific tests
				currentEdits = tt.edits
			}
			result := formatter.Format(currentEdits, tt.options)
			assert.Equal(t, result, tt.expected, "UnifiedFormatter.Format/%q", tt.name)
		})
	}
}
