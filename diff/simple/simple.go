package simple

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/lcs"
	"github.com/neticdk/go-stdlib/diff/internal/lines"
)

// Diff computes differences between two strings using a simple diff algorithm.
// Currently never returns an error.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := lines.Split(a)
	bLines := lines.Split(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
// Currently never returns an error.
func DiffStrings(a, b []string, opts ...Option) (string, error) {
	options := applyOptions(opts...)

	// Compute edit script using simple diff algorithm
	edits := computeEditScript(a, b)

	// Format the diff output
	var sb strings.Builder

	// Track line numbers if enabled
	aLineNum := 1
	bLineNum := 1

	for _, edit := range edits {
		switch edit.Kind {
		case diff.Equal:
			if options.showLineNumbers {
				sb.WriteString(fmt.Sprintf("%4d %4d   ", aLineNum, bLineNum))
			} else {
				sb.WriteString("  ")
			}
			sb.WriteString(edit.Text)
			sb.WriteString("\n")
			aLineNum++
			bLineNum++
		case diff.Delete:
			if options.showLineNumbers {
				sb.WriteString(fmt.Sprintf("%4d      - ", aLineNum))
			} else {
				sb.WriteString("- ")
			}
			sb.WriteString(edit.Text)
			sb.WriteString("\n")
			aLineNum++
		case diff.Insert:
			if options.showLineNumbers {
				sb.WriteString(fmt.Sprintf("     %4d + ", bLineNum))
			} else {
				sb.WriteString("+ ")
			}
			sb.WriteString(edit.Text)
			sb.WriteString("\n")
			bLineNum++
		}
	}

	return sb.String(), nil
}

// computeEditScript implements a simple diff algorithm based on longest common subsequence
func computeEditScript(a, b []string) []diff.Line {
	edits := []diff.Line{}

	// Use a simple longest common subsequence approach
	longest := lcs.LongestCommonSubsequence(a, b)
	aIndex, bIndex := 0, 0

	for _, item := range longest {
		// Add deletions for unmatched items in a
		for aIndex < item.AIndex {
			edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
			aIndex++
		}

		// Add insertions for unmatched items in b
		for bIndex < item.BIndex {
			edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
			bIndex++
		}

		// Add the matching item
		edits = append(edits, diff.Line{Kind: diff.Equal, Text: a[aIndex]})
		aIndex++
		bIndex++
	}

	// Handle remaining items in a
	for aIndex < len(a) {
		edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
		aIndex++
	}

	// Handle remaining items in b
	for bIndex < len(b) {
		edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
		bIndex++
	}

	return edits
}
