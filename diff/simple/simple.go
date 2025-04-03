package simple

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
)

// Diff computes differences between two strings using a simple diff algorithm.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := splitLines(a)
	bLines := splitLines(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
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
	lcs := longestCommonSubsequence(a, b)
	aIndex, bIndex := 0, 0

	for _, item := range lcs {
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

// MatchItem represents a matching item in two sequences
type MatchItem struct {
	AIndex, BIndex int
}

// longestCommonSubsequence finds the longest common subsequence between two string slices
func longestCommonSubsequence(a, b []string) []MatchItem {
	// Create a map of values to positions in B for faster lookup
	bValues := make(map[string][]int)
	for i, val := range b {
		bValues[val] = append(bValues[val], i)
	}

	var lcs []MatchItem

	// Find matches
	for i, aVal := range a {
		if positions, ok := bValues[aVal]; ok {
			// Try to find the best match position
			bestMatch := -1

			for _, pos := range positions {
				// Check if this match extends the current LCS
				valid := true
				for j := len(lcs) - 1; j >= 0; j-- {
					if lcs[j].BIndex >= pos {
						valid = false
						break
					}
				}

				if valid && (bestMatch == -1 || pos < bestMatch) {
					bestMatch = pos
				}
			}

			if bestMatch != -1 {
				lcs = append(lcs, MatchItem{i, bestMatch})
			}
		}
	}

	return lcs
}

// splitLines splits a string into lines, handling empty strings and trailing newlines
func splitLines(s string) []string {
	// Special case for completely empty string
	if s == "" {
		return []string{}
	}

	// Split the string by newlines
	lines := strings.Split(s, "\n")

	// If the string ends with a newline, the split will produce an empty string
	// at the end - we should remove it to avoid confusing diff output
	if s[len(s)-1] == '\n' {
		lines = lines[:len(lines)-1]
	}

	return lines
}
