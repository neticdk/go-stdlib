package diff

import (
	"fmt"
	"strings"
)

// OutputFormat represents the format of the diff output
type OutputFormat int

const (
	DefaultContextLines    = 3
	DefaultOutputFormat    = FormatContext
	DefaultShowLineNumbers = true
)

const (
	FormatContext OutputFormat = iota
	FormatUnified
)

func (o OutputFormat) String() string {
	switch o {
	case FormatContext:
		return "context"
	case FormatUnified:
		return "unified"
	default:
		return "Unknown"
	}
}

// The Formatter interface
type Formatter interface {
	// Format formats a diff
	Format(edits []Line, options FormatOptions) string
}

// FormatOptions represents the options for formatting the diff output
type FormatOptions struct {
	// OutputFormat specifies the format of the diff output
	// Default: FormatContext
	OutputFormat OutputFormat
	// ContextLines specifies the number of lines of context to show around each change
	// If ContextLines is 0, no context lines will be shown
	// Default: 3
	ContextLines int
	// ShowLineNumbers indicates whether line numbers should be shown in the diff output
	// Default: false
	ShowLineNumbers bool
}

// Validate validates the FormatOptions instance
func (fo *FormatOptions) Validate() error {
	if fo.ContextLines < 0 {
		return fmt.Errorf("ContextLines must be non-negative")
	}
	if fo.OutputFormat != FormatContext && fo.OutputFormat != FormatUnified {
		return fmt.Errorf("OutputFormat must be FormatContext or FormatUnified")
	}
	return nil
}

// ContextFormatter is a Formatter implementation using context diffs
type ContextFormatter struct{}

// Format formats the diff output using the ContextFormatter
func (c ContextFormatter) Format(edits []Line, options FormatOptions) string {
	var sb strings.Builder
	aLineNum, bLineNum := 1, 1 // Global line counters

	// Determine which edits belong to printable chunks
	chunkRanges := calculateChunkRanges(edits, options.ContextLines)

	firstChunkPrinted := false

	for idx, edit := range edits {
		// Determine if the current edit falls within any printable chunk range
		shouldPrint := false
		isStartOfNewGroup := false // Flag if this edit starts a new printable group
		for _, r := range chunkRanges {
			if idx >= r.start && idx < r.end {
				shouldPrint = true
				// Check if this is the very first line of a group marked as 'new'
				if idx == r.start && r.isNewGroup {
					isStartOfNewGroup = true
				}
				break // No need to check other ranges for this edit
			}
		}

		// Print separator *before* printing the first line of a new group
		// Only if context > 0 and it's not the absolute first chunk being printed.
		if isStartOfNewGroup && firstChunkPrinted && options.ContextLines > 0 {
			sb.WriteString("...\n")
		}

		// Print the edit line if it's within a chunk range
		if shouldPrint {
			firstChunkPrinted = true // Mark that we've started printing content
			switch edit.Kind {
			case Equal:
				if options.ShowLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d %4d   ", aLineNum, bLineNum))
				} else {
					sb.WriteString("  ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
			case Delete:
				if options.ShowLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d      - ", aLineNum))
				} else {
					sb.WriteString("- ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
			case Insert:
				if options.ShowLineNumbers {
					sb.WriteString(fmt.Sprintf("     %4d + ", bLineNum))
				} else {
					sb.WriteString("+ ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
			}
		}

		// ALWAYS update global line numbers based on the original edit kind
		switch edit.Kind {
		case Equal:
			aLineNum++
			bLineNum++
		case Delete:
			aLineNum++
		case Insert:
			bLineNum++
		}
	}
	return sb.String()
}

// UnifiedFormatter is a Formatter implementation using unified diffs
type UnifiedFormatter struct{}

// Format formats the diff output using the UnifiedFormatter
func (u UnifiedFormatter) Format(edits []Line, options FormatOptions) string {
	var sb strings.Builder

	// Standard headers for unified format (generic filenames)
	sb.WriteString("--- a\n")
	sb.WriteString("+++ b\n")

	chunkRanges := calculateChunkRanges(edits, options.ContextLines)

	for _, r := range chunkRanges {
		// Skip empty ranges (shouldn't happen with current logic, but safe)
		if r.start >= r.end {
			continue
		}

		// Calculate Hunk Header
		hunkHeader := calculateHunkHeader(edits, r)

		// Format the @@ hunk header
		sb.WriteString(fmt.Sprintf("@@ -%s +%s @@\n", hunkHeader.aRange, hunkHeader.bRange))

		// Write Hunk Body
		// Iterate through the chunk range to write the body
		for i := r.start; i < r.end; i++ {
			// Check index bounds just in case
			if i >= len(edits) {
				break
			}
			edit := edits[i]
			switch edit.Kind {
			case Equal:
				sb.WriteString(" ")
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
			case Delete:
				sb.WriteString("-")
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
			case Insert:
				sb.WriteString("+")
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

type hunkHeader struct {
	aRange string
	bRange string
}

// calculateHunkHeader calculates the header string for a unified diff hunk.
func calculateHunkHeader(edits []Line, r chunkRange) hunkHeader {
	// The starting line numbers for the hunk are the *global* line numbers
	// corresponding to the *first line* included in this chunk range (r.Start).
	hunkAStart, hunkBStart := calculateStartLinesForIndex(edits, r.start)

	// Calculate lengths by iterating *only* through the edits within this range
	hunkALen := 0
	hunkBLen := 0
	for i := r.start; i < r.end; i++ {
		// Check index bounds just in case
		if i >= len(edits) {
			break
		}
		switch edits[i].Kind {
		case Equal, Delete:
			hunkALen++
		}
		switch edits[i].Kind {
		case Equal, Insert:
			hunkBLen++
		}
	}

	// Adjust display start for 0-length hunks
	hunkAStartDisplay := hunkAStart
	if hunkALen == 0 && hunkAStart > 0 { // Avoid 0,0 for empty files, keep 0,0 if start was 0
		hunkAStartDisplay = hunkAStart - 1
	} else if hunkALen == 0 && hunkAStart == 0 {
		hunkAStartDisplay = 0 // Handle edge case for empty file change
	}

	hunkBStartDisplay := hunkBStart
	if hunkBLen == 0 && hunkBStart > 0 {
		hunkBStartDisplay = hunkBStart - 1
	} else if hunkBLen == 0 && hunkBStart == 0 {
		hunkBStartDisplay = 0
	}

	aRange := fmt.Sprintf("%d,%d", hunkAStartDisplay, hunkALen)
	bRange := fmt.Sprintf("%d,%d", hunkBStartDisplay, hunkBLen)

	return hunkHeader{aRange: aRange, bRange: bRange}
}

// calculateStartLinesForIndex determines the original a and b line numbers
// corresponding to the start of the edit at the given targetIndex.
func calculateStartLinesForIndex(edits []Line, targetIndex int) (aLineNum, bLineNum int) {
	aLine, bLine := 1, 1
	for i := range targetIndex {
		if i >= len(edits) {
			break
		}
		switch edits[i].Kind {
		case Equal:
			aLine++
			bLine++
		case Delete:
			aLine++
		case Insert:
			bLine++
		}
	}
	return aLine, bLine
}
