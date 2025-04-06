package simple

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/diffcore"
)

// Diff computes differences between two strings using a simple diff algorithm.
// Currently never returns an error.
func Diff(a, b string, opts ...Option) string {
	if a == "" && b == "" {
		return ""
	}

	aLines := diffcore.SplitLines(a)
	bLines := diffcore.SplitLines(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using a simple diff algorithm.
// Currently never returns an error.
func DiffStrings(a, b []string, opts ...Option) string {
	return simpleDiffStrings(a, b, applyOptions(opts...))
}

func simpleDiffStrings(a, b []string, opts options) string {
	// Compute edit script using the shared simple LCS-based diff algorithm
	script := diffcore.ComputeEditsLCS(a, b)

	// Format the diff output
	var sb strings.Builder

	// Track line numbers if enabled
	aLineNum := 1
	bLineNum := 1

	// Group edits by type for context-aware output
	chunks := diffcore.GroupEditsByContext(script, opts.contextLines)

	for i, chunk := range chunks {
		// Add separator between chunks
		if i > 0 && opts.contextLines > 0 {
			sb.WriteString("...\n")
		}

		for _, edit := range chunk {
			switch edit.Kind {
			case diff.Equal:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d %4d   ", aLineNum, bLineNum))
				} else {
					sb.WriteString("  ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				aLineNum++
				bLineNum++
			case diff.Delete:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d      - ", aLineNum))
				} else {
					sb.WriteString("- ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				aLineNum++
			case diff.Insert:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("     %4d + ", bLineNum))
				} else {
					sb.WriteString("+ ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				bLineNum++
			}
		}
	}

	return sb.String()
}
