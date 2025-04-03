package simple

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/diffcore" // Import the new internal package
)

// Diff computes differences between two strings using a simple diff algorithm.
// Currently never returns an error.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := diffcore.SplitLines(a)
	bLines := diffcore.SplitLines(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
// Currently never returns an error.
func DiffStrings(a, b []string, opts ...Option) (string, error) {
	options := applyOptions(opts...)

	// Compute edit script using the shared simple LCS-based diff algorithm
	edits := diffcore.ComputeEditsLCS(a, b)

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
